package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"time"

	"fmt"
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
type payload struct{
	Topic string `json:"topic"`
	Data map[string]any `json:"data"`
}


type removeStruct struct{
	Name string
	date time.Time
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

	generalMailingList = mongoClient.Database("Notification-Service").Collection("generalMailingList")
	eventMailingList = mongoClient.Database("Notification-Service").Collection("eventMailingList")
	log.Println("Connected to MongoDB!")
}

// func (app * application) subscribeGeneral(w http.ResponseWriter, r *http.Request) {
// 	userEmail := r.Header.Get("Email")

// 	var result bson.M
// 	err := generalMailingList.FindOne(context.Background(), bson.M{"email": userEmail}).Decode(&result)

// 	if err == mongo.ErrNoDocuments {
// 		_, insertErr := generalMailingList.InsertOne(context.Background(), bson.M{"email": userEmail})
// 			if insertErr != nil {
// 				log.Println("Error inserting document: ", insertErr)
// 				http.Error(w, "Failed to subscribe", http.StatusInternalServerError)
// 				return
// 			}
	
// 			log.Println("Email inserted successfully!")
// 			w.WriteHeader(http.StatusOK)
// 			w.Write([]byte("Email successfully subscribed!"))
// 			app.background(func() {
// 				err := app.Mailer.Send(userEmail,"SubscribeTemplate.tmpl",nil)
// 				if err !=nil{
// 					app.Logger.Println(err)

// 				}
// 			})
			



// 	} else {
// 		log.Println("Email already subscribed")
// 		w.WriteHeader(http.StatusBadRequest)
// 		w.Write([]byte("Email already subscribed"))
// 	}	
// }

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

func (app *application) getEmails(Payload payload) ([]string, error) {
    emailsInterface, ok := Payload.Data["emails"].([]interface{})
    if !ok {
        app.Logger.Println("emails is not a slice of interfaces")
        return nil, fmt.Errorf("emails is not a slice of interfaces")
    }

    var emails []string

    for _, email := range emailsInterface {
        emailStr, ok := email.(string)

        if !ok {
            app.Logger.Println("Email Parse Error")
            return nil, fmt.Errorf("email is not a valid string")
        }
        emails = append(emails, emailStr)
    }

    return emails, nil
}
func (app *application) getEventName(Payload payload) (string, error) {
    eventName, ok := Payload.Data["event_name"].(string)
    if !ok {
        app.Logger.Println("Name Parse Error")
        return "", fmt.Errorf("event_name is not a valid string")
    }
    return eventName, nil
}
func (app *application) getEventDate(Payload payload) (string, error) {
    eventDate, ok := Payload.Data["event_date"].(string)
    if !ok {
        app.Logger.Println("Date Parse Error")
        return "", fmt.Errorf("date is not a valid string")
    }
    return eventDate, nil
}
func (app *application) getEventDescription(Payload payload) (string, error) {
    eventDesc, ok := Payload.Data["event_description"].(string)
    if !ok {
        app.Logger.Println("Description Parse Error")
        return "", fmt.Errorf("description is not a valid string")
    }
    return eventDesc, nil
}
func (app *application) getEventLocation(Payload payload) (string, error){
	eventLocation, ok := Payload.Data["event_location"].(string)
    if !ok {
        app.Logger.Println("Name Parse Error")
        return "", fmt.Errorf("event location is not a valid string")
    }
    return eventLocation, nil

}



func (app * application) notify(w http.ResponseWriter, r *http.Request){
	var Payload payload
	err:=app.readJSON(w,r,&Payload)

	if err !=nil{
		app.serverErrorResponse(w,r,err)
		return
	}

	topic:=Payload.Topic

	// if err == nil{
	// 	return
	// }
	switch topic {
    case "event_add":
        app.eventAdd(Payload)
	case "event_remove":
        app.eventRemove(Payload)
	case "event_update":
        app.eventUpdate(Payload)
	case "event_register":
        app.eventRegister(Payload)
    }


}

func (app * application) eventAdd(Payload payload){}
func (app * application) eventRemove(Payload payload){
emails,err:=app.getEmails(Payload)
	if err != nil{
		app.Logger.Println("Emails Parse Error")
	}
eventName,err:=app.getEventName(Payload)
	if err != nil{
		app.Logger.Println("event Name Parse Error")
	}
eventDate, err := app.getEventDate(Payload)
	if err !=nil {
		app.Logger.Println("Date Parse Error")
		return
	}

type removeStruct struct{
		Name string
		Date string
}
data:=removeStruct{
	Name:eventName, 
	Date: eventDate,
}

app.background(func() {
		err := app.Mailer.Send(emails,"EventDeleteTemplate.tmpl",data)
		if err !=nil{
			app.Logger.Println(err)

		}
})
}
func (app * application) eventUpdate(Payload payload){
emails,err:=app.getEmails(Payload)
	if err != nil{
		app.Logger.Println("Emails Parse Error")
	}
eventName,err:=app.getEventName(Payload)
	if err != nil{
		app.Logger.Println("event Name Parse Error")
	}
eventDate, err := app.getEventDate(Payload)
	if err !=nil {
		app.Logger.Println("Date Parse Error")
		return
	}
updateDesc,err:=app.getEventDescription(Payload)
if err !=nil {
	app.Logger.Println("description Parse Error")
	return
}
type updateStruct struct{
	Name string
	Date string
	Description string
}
data:=updateStruct{
Name:eventName, 
Date: eventDate,
Description:updateDesc,
}
app.background(func() {
	err := app.Mailer.Send(emails,"EventUpdateTemplate.tmpl",data)
	if err !=nil{
		app.Logger.Println(err)

	}
})
}
func (app * application) eventRegister(Payload payload){
emails,err:=app.getEmails(Payload)
	if err != nil{
		app.Logger.Println("Emails Parse Error")
	}
eventName,err:=app.getEventName(Payload)
	if err != nil{
		app.Logger.Println("event Name Parse Error")
	}
eventDate, err := app.getEventDate(Payload)
	if err !=nil {
		app.Logger.Println("Date Parse Error")
		return
	}
eventLocation, err := app.getEventLocation(Payload)
	if err !=nil {
		app.Logger.Println("Date Parse Error")
		return
	}

type appliedStruct struct{
		Name string
		Date string
		Location string
}
data:=appliedStruct{
	Name:eventName, 
	Date: eventDate,
	Location: eventLocation,
}

app.background(func() {
	err := app.Mailer.Send(emails,"EventRegisterTemplate.tmpl",data)
	if err !=nil{
		app.Logger.Println(err)

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
		Config: cfg,
		Logger: logger,
		Mailer: mailer.New(cfg.smtp.host, cfg.smtp.port, cfg.smtp.username, cfg.smtp.password, cfg.smtp.sender),
	}

	log.Printf("starting user service on %s\n", cfg.port)

	

	







	r := chi.NewRouter()
	connectToDb()

	
		
	// r.Post("/subscribe/{category}", app.subscribeGeneral)
	// r.Post("/unsubscribe", app.unsubscribeGeneral)
	r.Post("/notify",app.notify)

	err=http.ListenAndServe(app.Config.port, r)

	log.Fatal(err)
}
