package main

import (
	"context"
	"log"
	"net/http"
	"os"

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
	webPort = ":80"
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
	jwt struct {
		secret string
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
	clientOptions := options.Client().ApplyURI("mongodb://mongo:27017")
	clientOptions.SetAuth(options.Credential{
		Username: "admin",
		Password: "password",
	})
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

func (app *application) isSubscribed(EventType, email string) bool {
	filter := bson.M{
		"Event_Type": EventType,
		"emails":     email,
	}

	count, err := mailingList.CountDocuments(context.Background(), filter)
	if err != nil {
		return false
	}

	if count > 0 {
		return true
	}

	return false
}

func (app *application) subscribe(w http.ResponseWriter, r *http.Request) {
	userEmail, _, _, err := app.extractTokenData(r)
	if err != nil {
		app.writeJSON(w, http.StatusUnauthorized, envelope{"error": "Invalid token"}, nil)
		return
	}
	eventType := chi.URLParam(r, "eventType")

	isSubs := app.isSubscribed(eventType, userEmail)

	if isSubs {
		log.Println("Email Already subscribed to Mailing List!")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Email Already subscribed to Mailing List!"))
		return

	}

	_, err = mailingList.UpdateOne(context.Background(), bson.M{"Event_Type": eventType}, bson.M{"$addToSet": bson.M{"emails": userEmail}}, options.Update().SetUpsert(true))

	if err != nil {
		log.Fatal("Error updating Mailling List: ", err)
		http.Error(w, "Failed to subscribe", http.StatusInternalServerError)
		return
	}

	log.Println("Email successfully subscribed to Mailing List!")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Email successfully subscribed to Mailing List!"))
	eventTypeText := "All"
	if eventType != "general" {
		eventTypeText = eventType
	}
	type subStruct struct {
		Email string
		Type  string
	}

	data := subStruct{
		Email: userEmail,
		Type:  eventTypeText,
	}

	app.background(func() {
		err := app.Mailer.Send([]string{userEmail}, "SubscribeTemplate.tmpl", data)
		if err != nil {
			app.Logger.Println(err)

		}
	})
}
func (app *application) unsubscribe(w http.ResponseWriter, r *http.Request) {
	userEmail := r.Header.Get("Email")
	eventType := chi.URLParam(r, "eventType")

	isSubs := app.isSubscribed(eventType, userEmail)

	if !isSubs {
		log.Println("Email Not subscribed to Mailing List!")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Email Not subscribed to Mailing List!"))
		return

	}

	_, err := mailingList.UpdateOne(
		context.Background(),
		bson.M{"Event_Type": eventType},
		bson.M{"$pull": bson.M{"emails": userEmail}},
	)

	if err != nil {
		log.Fatal("Error updating Mailling List: ", err)
		http.Error(w, "Failed to unsubscribe", http.StatusInternalServerError)
		return
	}

	log.Println("Email successfully unsubscribed to Mailing List!")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Email successfully unsubscribed to Mailing List!"))

	eventTypeText := "All"
	if eventType != "general" {
		eventTypeText = eventType
	}
	type subStruct struct {
		Email string
		Type  string
	}

	data := subStruct{
		Email: userEmail,
		Type:  eventTypeText,
	}

	app.background(func() {
		err := app.Mailer.Send([]string{userEmail}, "UnsubscribeTemplate.tmpl", data)
		if err != nil {
			app.Logger.Println(err)

		}
	})
}

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
func (app *application) getEventLocation(Payload payload) (string, error) {
	eventLocation, ok := Payload.Data["event_location"].(string)
	if !ok {
		app.Logger.Println("Name Parse Error")
		return "", fmt.Errorf("event location is not a valid string")
	}
	return eventLocation, nil
}
func (app *application) getEventType(Payload payload) (string, error) {
	eventType, ok := Payload.Data["event_type"].(string)
	if !ok {
		app.Logger.Println("Type Parse Error")
		return "", fmt.Errorf("event_type is not a valid string")
	}
	return eventType, nil

}

func (app *application) notify(w http.ResponseWriter, r *http.Request) {
	var Payload payload
	err := app.readJSON(w, r, &Payload)

	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	topic := Payload.Topic

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
	default:
		app.Logger.Println("Unknown Topic")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Unknown Topic"))

	}

}

func (app *application) eventAdd(Payload payload) {
	event_type, err := app.getEventType(Payload)
	if err != nil {
		app.Logger.Println("Type Parse Error")
	}
	event_name, err := app.getEventName(Payload)
	if err != nil {
		app.Logger.Println("event Name Parse Error")
	}
	event_location, err := app.getEventLocation(Payload)
	if err != nil {
		app.Logger.Println("event Name Parse Error")
	}
	event_date, err := app.getEventDate(Payload)
	if err != nil {
		app.Logger.Println("Date Parse Error")
		return
	}
	Desc, err := app.getEventDescription(Payload)
	if err != nil {
		app.Logger.Println("description Parse Error")
		return
	}

	filter := bson.M{
		"Event_Type": bson.M{"$in": []interface{}{"general", event_type}},
	}

	emails, err := mailingList.Find(
		context.Background(),
		filter,
		options.Find().SetProjection(bson.M{"emails": 1, "_id": 0}),
	)
	if err != nil {
		app.Logger.Printf("Find error: %v", err)
		return
	}

	type EmailStruct struct {
		Emails []string `bson:"emails"`
	}

	var results []string
	for emails.Next(context.TODO()) {
		var result EmailStruct
		if decodeErr := emails.Decode(&result); decodeErr != nil {
			app.Logger.Printf("Decode error: %v", decodeErr)
			continue
		}
		results = append(results, result.Emails...)
	}

	if err := emails.Err(); err != nil {
		app.Logger.Printf("Cursor error: %v", err)
	}

	type eventAddStruct struct {
		Name        string
		Date        string
		Location    string
		Description string
	}

	data := eventAddStruct{
		Name:        event_name,
		Date:        event_date,
		Location:    event_location,
		Description: Desc,
	}
	fmt.Printf(data.Date)

	app.background(func() {
		err := app.Mailer.Send(results, "EventAddTemplate.tmpl", data)
		if err != nil {
			app.Logger.Println(err)

		}
	})
}
func (app *application) eventRemove(Payload payload) {
	emails, err := app.getEmails(Payload)
	if err != nil {
		app.Logger.Println("Emails Parse Error")
	}
	eventName, err := app.getEventName(Payload)
	if err != nil {
		app.Logger.Println("event Name Parse Error")
	}
	eventDate, err := app.getEventDate(Payload)
	if err != nil {
		app.Logger.Println("Date Parse Error")
		return
	}

	type removeStruct struct {
		Name string
		Date string
	}
	data := removeStruct{
		Name: eventName,
		Date: eventDate,
	}

	app.background(func() {
		err := app.Mailer.Send(emails, "EventDeleteTemplate.tmpl", data)
		if err != nil {
			app.Logger.Println(err)

		}
	})
}
func (app *application) eventUpdate(Payload payload) {
	emails, err := app.getEmails(Payload)
	if err != nil {
		app.Logger.Println("Emails Parse Error")
	}
	eventName, err := app.getEventName(Payload)
	if err != nil {
		app.Logger.Println("event Name Parse Error")
	}
	eventDate, err := app.getEventDate(Payload)
	if err != nil {
		app.Logger.Println("Date Parse Error")
		return
	}
	updateDesc, err := app.getEventDescription(Payload)
	if err != nil {
		app.Logger.Println("description Parse Error")
		return
	}
	type updateStruct struct {
		Name        string
		Date        string
		Description string
	}
	data := updateStruct{
		Name:        eventName,
		Date:        eventDate,
		Description: updateDesc,
	}
	app.background(func() {
		err := app.Mailer.Send(emails, "EventUpdateTemplate.tmpl", data)
		if err != nil {
			app.Logger.Println(err)

		}
	})
}

func (app *application) eventRegister(Payload payload) {

	emails, err := app.getEmails(Payload)
	if err != nil {
		app.Logger.Println("Emails Parse Error")
	}
	eventName, err := app.getEventName(Payload)
	if err != nil {
		app.Logger.Println("event Name Parse Error")
	}
	eventDate, err := app.getEventDate(Payload)
	if err != nil {
		app.Logger.Println("Date Parse Error")
		return
	}
	eventLocation, err := app.getEventLocation(Payload)
	if err != nil {
		app.Logger.Println("Date Parse Error")
		return
	}

	type appliedStruct struct {
		Name     string
		Date     string
		Location string
	}
	data := appliedStruct{
		Name:     eventName,
		Date:     eventDate,
		Location: eventLocation,
	}

	app.background(func() {
		err := app.Mailer.Send(emails, "EventRegisterTemplate.tmpl", data)
		if err != nil {
			app.Logger.Println(err)

		}
	})
}

func main() {
	var cfg config
	cfg.port = webPort
	cfg.env = webEnv
	cfg.smtp.host = os.Getenv("MAILHOG_HOST")
	portStr := os.Getenv("MAILHOG_PORT")
	cfg.smtp.username = os.Getenv("MAILHOG_USERNAME")
	cfg.smtp.password = os.Getenv("MAILHOG_PASSWORD")
	cfg.smtp.sender = os.Getenv("SENDER_EMAIL")

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

	r.Post("/subscribe/{eventType}", app.subscribe)
	r.Post("/unsubscribe/{eventType}", app.unsubscribe)
	r.Post("/notify", app.notify)

	err = http.ListenAndServe(app.Config.port, r)

	log.Fatal(err)
}
