package main

import (
	"context"
	"log"
	"net/http"
	"os"

	//"fmt"
	"strconv"

	"github.com/MohamedHossam2004/Event-Planner/notification-service/internal/mailer"
	"github.com/go-chi/chi/v5"

	// "github.com/go-chi/chi/v5/middleware"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	webPort = ":8080"
	webEnv  = "development"
)

var (
	mongoClient *mongo.Client
	mailingList *mongo.Collection
)

type config struct {
	port string
	env  string
	smtp struct {
		host     string
		port     int
		username string
		password string
		sender   string
	}
}
type payload struct {
	Topic string         `json:"topic"`
	Data  map[string]any `json:"data"`
}

type application struct {
	Config config
	Logger *log.Logger
	Mailer mailer.Mailer
}

func connectToDb() {
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27018")
	var err error
	mongoClient, err = mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	err = mongoClient.Ping(context.Background(), nil)
	if err != nil {
		log.Fatal(err)
	}

	mongoClient.Database("Notification-Service").CreateCollection(context.Background(), "mailingList")
	mailingList = mongoClient.Database("Notification-Service").Collection("mailingList")
	log.Println("Connected to MongoDB!")
}

func (app *application) subscribe(w http.ResponseWriter, r *http.Request) {
	userEmail := r.Header.Get("Email")
	category := chi.URLParam(r, "category")

	checker := mailingList.FindOne(context.Background(), bson.M{"category": category})
	if checker.Err() != nil {
		log.Println("There is no category with this name, making one...: ", checker.Err())
		mailingList.InsertOne(context.Background(), bson.M{"category": category, "emails": []string{userEmail}})
	}
	_, err := mailingList.UpdateOne(context.Background(), bson.M{"category": category}, bson.M{"$addToSet": bson.M{"emails": userEmail}})
	if err != nil {
		log.Println("Error adding email: ", err)
		return
	}

	log.Printf("Successfully subscribed %s to category %s", userEmail, category)
	w.WriteHeader(http.StatusOK)
}

func (app *application) unsubscribe(w http.ResponseWriter, r *http.Request) {
	userEmail := r.Header.Get("Email")
	category := chi.URLParam(r, "category")

	checker := mailingList.FindOne(context.Background(), bson.M{"category": category})
	if checker.Err() != nil {
		log.Println("There is no category with this name, making one...: ", checker.Err())
		mailingList.InsertOne(context.Background(), bson.M{"category": category, "emails": []string{userEmail}})
	}
	_, err := mailingList.UpdateOne(context.Background(), bson.M{"category": category}, bson.M{"$pull": bson.M{"emails": userEmail}})
	if err != nil {
		log.Println("Error removing email: ", err)
		return
	}

	log.Printf("Successfully unsubscribed %s from category %s", userEmail, category)
	w.WriteHeader(http.StatusOK)
}

// func (app * application) unsubscribeGeneral(w http.ResponseWriter, r *http.Request) {
// 	userEmail := r.Header.Get("Email")

// 	if userEmail == "" {
// 		http.Error(w, "Email header is required", http.StatusBadRequest)
// 		return
// 	}

// 	res, err := generalMailingList.DeleteOne(context.Background(), bson.M{"email": userEmail})

// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	if res.DeletedCount == 0 {
// 		http.Error(w, "Email not found in the mailing list or already unsubscribed", http.StatusNotFound)
// 		return
// 	}

// 	w.WriteHeader(http.StatusOK)
// 	w.Write([]byte("Successfully unsubscribed"))
// 	log.Printf("Email '%s' has been unsubscribed", userEmail)
// 	app.background(func() {
// 		err := app.mailer.Send(userEmail,"UnsubscribeTemplate.tmpl",nil)
// 		if err !=nil{
// 			app.logger.Println(err)

// 		}
// 	})
// }

// func (app * application) notify(w http.ResponseWriter, r *http.Request){
// 	var Payload payload
// 	err:=app.readJSON(w,r,&Payload)

// 	if err !=nil{
// 		app.serverErrorResponse(w,r,err)
// 		return
// 	}

// 	topic:=Payload.Topic

// 	// if err == nil{
// 	// 	return
// 	// }
// 	switch topic {
//     case "event_add":
//         app.eventAdd(Payload)
// 	case "event_remove":
//         app.eventRemove(Payload)
// 	case "event_update":
//         app.eventUpdate(Payload)
//     default:
//        return
//     }

// }

// func (app * application) eventAdd(Payload payload){

// }
// func (app * application) eventRemove(Payload payload){
// emails, ok := Payload.Data["emails"].([]string)
// 	if !ok {
// 		// Handle the case where the assertion fails
// 		return
// 	}
// eventName, ok := Payload.Data["event_name"].(string)
// 	if !ok {
// 		return
// 	}
// eventDate, ok := Payload.Data["date"].(string)
// 	if !ok {
// 		return
// 	}

// app.background(func() {
// 		err := app.Mailer.Send(emails,"SubscribeTemplate.tmpl",nil)
// 		if err !=nil{
// 			app.Logger.Println(err)

// 		}
// })

func main() {
	var cfg config
	cfg.port = webPort
	cfg.env = webEnv
	cfg.smtp.host = "localhost"
	portStr := "1025"
	cfg.smtp.username = ""
	cfg.smtp.password = ""
	cfg.smtp.sender = "giu-event-hub@giu-uni.de"

	if cfg.smtp.host == "" || portStr == "" {
		log.Fatal("Environment variables for Mailhog are not set")
	}
	var err error
	cfg.smtp.port, err = strconv.Atoi(portStr)
	if err != nil {
		log.Fatal("Error: Invalid MAILHOG_PORT value: %s\n", portStr)
		return
	}

	logger := log.New(os.Stdout, "", log.Ldate|log.LUTC)
	app := &application{
		Config: cfg,
		Logger: logger,
		Mailer: mailer.New(cfg.smtp.host, cfg.smtp.port, cfg.smtp.username, cfg.smtp.password, cfg.smtp.sender),
	}

	log.Printf("starting user service on %s\n", cfg.port)

	r := chi.NewRouter()
	connectToDb()

	r.Post("/subscribe/{category}", app.subscribe)
	r.Delete("/unsubscribe/{category}", app.unsubscribe)
	// r.Post("/notify",app.notify)

	err = http.ListenAndServe(app.Config.port, r)

	log.Fatal(err)
}
