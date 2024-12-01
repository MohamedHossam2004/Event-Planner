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

	mailingList = mongoClient.Database("Notification-Service").Collection("mailingList")
	log.Println("Connected to MongoDB!")
}

func (app *application) subscribe(w http.ResponseWriter, r *http.Request) {
	userEmail := r.Header.Get("Email")
	category := chi.URLParam(r, "category")

	switch category {
	case "general":
		{
			cursor, err := mailingList.Find(context.Background(), bson.M{})
			if err != nil {
				log.Fatal("Error fetching categories: ", err)
				http.Error(w, "Failed to fetch categories", http.StatusInternalServerError)
				return
			}
			defer cursor.Close(context.Background())

			for cursor.Next(context.Background()) {
				var doc bson.M
				if err := cursor.Decode(&doc); err != nil {
					log.Println("Error decoding document: ", err)
					continue
				}

				_, err = mailingList.UpdateOne(context.Background(), bson.M{"category": doc["category"]}, bson.M{"$addToSet": bson.M{"emails": userEmail}})
				if err != nil {
					log.Println("Error adding email: ", err)
				}
			}
			log.Println("Email added to all categories!")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("Email successfully subscribed to all categories!"))
		}
	default:
		{
			_, err := mailingList.UpdateOne(context.Background(), bson.M{"category": category}, bson.M{"$addToSet": bson.M{"emails": userEmail}}, options.Update().SetUpsert(true))
			if err != nil {
				log.Fatal("Error updating category: ", err)
				http.Error(w, "Failed to subscribe", http.StatusInternalServerError)
				return
			}
		}
		log.Printf("Successfully subscribed %s to category %s", userEmail, category)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Successfully subscribed to category!"))
	}

}

func (app *application) unsubscribe(w http.ResponseWriter, r *http.Request) {
	userEmail := r.Header.Get("Email")
	log.Printf("Unsubscribing user: %s", userEmail)
	category := chi.URLParam(r, "category")
	log.Printf("from %s", category)

	switch category {
	case "general":
		{
			cursor, err := mailingList.Find(context.Background(), bson.M{})
			if err != nil {
				log.Fatal("Error fetching categories: ", err)
				http.Error(w, "Failed to fetch categories", http.StatusInternalServerError)
				return
			}
			defer cursor.Close(context.Background())

			for cursor.Next(context.Background()) {
				doc := bson.M{}
				if err := cursor.Decode(&doc); err != nil {
					log.Println("Error decoding document: ", err)
					continue
				}

				err := mailingList.FindOneAndDelete(context.Background(), bson.M{"category": doc["category"], "emails": userEmail})
				if err.Err() != nil {
					log.Println("Error removing email: ", err.Err())
				}
				log.Println("Email removed from all categories!")
				w.WriteHeader(http.StatusOK)
				w.Write([]byte("Successfully unsubscribed from all categories!"))
			}
		}
	default:
		{
			err := mailingList.FindOneAndDelete(context.Background(), bson.M{"category": category, "emails": userEmail})
			if err.Err() != nil {
				log.Fatal("Error removing email: ", err)
				http.Error(w, "Failed to unsubscribe", http.StatusInternalServerError)
				return
			}
			log.Printf("Successfully unsubscribed %s from categorty %s", userEmail, category)
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("Successfully unsubscribed from category!"))
		}
	}
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
	log.Printf("Setting up route: /unsubscribe/{category}")
	r.Delete("/unsubscribe/{category}", app.unsubscribe)
	// r.Post("/notify",app.notify)

	err = http.ListenAndServe(app.Config.port, r)

	log.Fatal(err)
}
