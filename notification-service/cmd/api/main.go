package main

import (
	"context"
	"log"
	"net/http"
	"os"
	//"fmt"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/MohamedHossam2004/Event-Planner/notification-service/internal/mailer"
	// "github.com/go-chi/chi/v5/middleware"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	webPort = ":8080"
	webEnv  = "development"
)
var (
	mongoClient        *mongo.Client
	generalMailingList *mongo.Collection
	eventMailingList   *mongo.Collection
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

type application struct {
	config config
	logger *log.Logger
	mailer mailer.Mailer
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

	generalMailingList = mongoClient.Database("Notification-Service").Collection("generalMailingList")
	eventMailingList = mongoClient.Database("Notification-Service").Collection("eventMailingList")
	log.Println("Connected to MongoDB!")

}

func (app * application) subscribeGeneral(w http.ResponseWriter, r *http.Request) {
	userEmail := r.Header.Get("Email")

	var result bson.M
	err := generalMailingList.FindOne(context.Background(), bson.M{"email": userEmail}).Decode(&result)

	if err == mongo.ErrNoDocuments {
		_, insertErr := generalMailingList.InsertOne(context.Background(), bson.M{"email": userEmail})
			if insertErr != nil {
				log.Println("Error inserting document: ", insertErr)
				http.Error(w, "Failed to subscribe", http.StatusInternalServerError)
				return
			}
	
			log.Println("Email inserted successfully!")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("Email successfully subscribed!"))
			app.background(func() {
				err := app.mailer.Send(userEmail,"SubscribeTemplate.tmpl",nil)
				if err !=nil{
					app.logger.Println(err)

				}
			})
			



	} else {
		log.Println("Email already subscribed")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Email already subscribed"))
	}

	
}

func subscribeEvent(w http.ResponseWriter, r *http.Request) {
	userEmail := r.Header.Get("Email")
	eventID := chi.URLParam(r, "id")

	_, err := eventMailingList.UpdateOne(context.Background(), bson.M{"event_id": eventID}, bson.M{"$addToSet": bson.M{"emails": userEmail}},options.Update().SetUpsert(true))

	if err != nil {
		log.Fatal("Error updating event: ", err)
		http.Error(w, "Failed to subscribe", http.StatusInternalServerError)
		return
	}

	log.Println("Email successfully subscribed to event!")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Email successfully subscribed to event!"))

}

func (app * application) unsubscribe(w http.ResponseWriter, r *http.Request) {
	userEmail := r.Header.Get("Email")

	if userEmail == "" {
		http.Error(w, "Email header is required", http.StatusBadRequest)
		return
	}


	res, err := generalMailingList.DeleteOne(context.Background(), bson.M{"email": userEmail})

	if err != nil {
		log.Fatal(err)
	}

	if res.DeletedCount == 0 {
		http.Error(w, "Email not found in the mailing list or already unsubscribed", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Successfully unsubscribed"))
	log.Printf("Email '%s' has been unsubscribed", userEmail)
	app.background(func() {
		err := app.mailer.Send(userEmail,"UnsubscribeTemplate.tmpl",nil)
		if err !=nil{
			app.logger.Println(err)

		}
	})
}

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
		config: cfg,
		logger: logger,
		mailer: mailer.New(cfg.smtp.host, cfg.smtp.port, cfg.smtp.username, cfg.smtp.password, cfg.smtp.sender),
	}

	log.Printf("starting user service on %s\n", cfg.port)

	

	







	r := chi.NewRouter()
	connectToDb()

	r.Route("/subscribe", func(r chi.Router) {
		
		r.Post("/general", app.subscribeGeneral)
		r.Post("/event/{id}", subscribeEvent)
	})

	r.Post("/unsubscribe", app.unsubscribe)

	err=http.ListenAndServe(app.config.port, r)

	log.Fatal(err)
}
