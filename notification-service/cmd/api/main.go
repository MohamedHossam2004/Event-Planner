package main

import(
	"net/http"
	"context"
	"log"
	//"github.com/go-chi/chi/v5"
	//"github.com/go-chi/chi/v5/middleware"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	mongoClient *mongo.Client
	generalMailingList *mongo.Collection
	eventMailingList *mongo.Collection
)



type email struct{
	email string `bson:"email"`
}



func connectToDb(){
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
	eventMailingList=mongoClient.Database("Notification-Service").Collection("eventMailingList")
	log.Println("Connected to MongoDB!")

}


func subscribe(w http.ResponseWriter, r *http.Request){
	
	
	
}

func main(){
	//r:=chi.NewRouter()
	connectToDb()
	// notification := Notification{ID: 10}
	// _, err := userCollection.InsertOne(context.Background(), notification)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// log.Println("Notification inserted successfully!")

}






