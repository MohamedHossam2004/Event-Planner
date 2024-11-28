package main

import (
	"context"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"

	// "github.com/go-chi/chi/v5/middleware"
	"go.mongodb.org/mongo-driver/mongo"

	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	mongoClient        *mongo.Client
	generalMailingList *mongo.Collection
	eventMailingList   *mongo.Collection
)

type EmailStruct struct {
	email string `bson:"email"`
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

func subscribeGeneral(w http.ResponseWriter, r *http.Request) {
	userEmail := r.Header.Get("Email")
	email := EmailStruct{email: userEmail}
	var result EmailStruct
	checkEmail := generalMailingList.FindOne(context.Background(), email).Decode(&result)

	if checkEmail == mongo.ErrNoDocuments {
		log.Println("Email not found in mailing list")
		_, err := generalMailingList.InsertOne(context.Background(), email)
		if err != nil {
			log.Fatal("Error inserting document: ", err)
			http.Error(w, "Failed to subscribe", http.StatusInternalServerError)
			return
		}

		log.Println("Email inserted successfully!")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Email successfully subscribed!"))
	} else {

		log.Println("Email already subscribed")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Email already subscribed"))
	}
}

func unsubscribe(w http.ResponseWriter, r *http.Request) {
	userEmail := r.Header.Get("Email")

	if userEmail == "" {
		http.Error(w, "Email header is required", http.StatusBadRequest)
		return
	}

	email := EmailStruct{email: userEmail}

	res, err := generalMailingList.DeleteOne(context.Background(), email)

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
}

func main() {
	r := chi.NewRouter()
	connectToDb()

	r.Route("/subscribe", func(r chi.Router) {
		// Add a POST handler for /subscribe/general
		r.Post("/general", subscribeGeneral)
	})

	r.Post("/unsubscribe", unsubscribe)

	http.ListenAndServe(":8080", r)
}
