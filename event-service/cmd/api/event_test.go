package main

import (
	"errors"
	"log"
	"os"
	"testing"
	"time"

	"github.com/MohamedHossam2004/Event-Planner/event-service/internal/data"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson/primitive"

	// "fmt"
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"

	"github.com/stretchr/testify/assert"
)

type MockEventModel struct {
	mock.Mock
}

func (m *MockEventModel) GetEventByID(id primitive.ObjectID) (*data.Event, error) {
	args := m.Called(id)
	return args.Get(0).(*data.Event), args.Error(1)
}

func (m *MockEventModel) CreateEvent(event *data.Event) (*data.Event, error) {
	args := m.Called(event)
	return args.Get(0).(*data.Event), args.Error(1)
}

func (m *MockEventModel) UpdateEvent(id primitive.ObjectID, event *data.Event) (*data.Event, error) {
	args := m.Called(id, event)
	return args.Get(0).(*data.Event), args.Error(1)
	
}

func (m *MockEventModel) DeleteEvent(id primitive.ObjectID) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockEventModel) GetAllEvents() ([]data.Event, error) {
	args := m.Called()
	return args.Get(0).([]data.Event), args.Error(1)
}

func TestUpdateEventHandler(t *testing.T) {
	mockEventAppModel := new(MockEventAppModel)
	mockEventModel := new(MockEventModel)
	mockTokenExtractor := new(MockTokenExtractor)

	rabbitConn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		t.Fatalf("Failed to connect to RabbitMQ: %v", err)
	}

	app := &application{
		Logger: log.New(os.Stdout, "", 0),
		config: config{port: "80", env: "development"},
		models: data.Models{
			EventApps: mockEventAppModel,
			Event:     mockEventModel,
		},
		Rabbit:         rabbitConn,
		tokenExtractor: mockTokenExtractor,
	}

	tests := []struct {
		name           string
		eventId        string
		eventData      data.Event
		expectedStatus int
		expectedBody   string
		setupMock      func(mockEventAppModel *MockEventAppModel, mockEventModel *MockEventModel, mockTokenExtractor *MockTokenExtractor)
	}{
		{
			name:    "Successful Event Update",
			eventId:  primitive.NewObjectID().Hex(),
			eventData: data.Event{
				ID:   primitive.NewObjectID(),
				Date: time.Date(2025, 7, 15, 18, 0, 0, 0, time.UTC),
				Type: "CONFERENCE",
				Name: "Tech Conference 2025",
				Location: data.Location{
					Address: "123 Main Street",
					City:    "San Francisco",
					State:   "CA",
					Country: "USA",
				},
				MaxCapacity:          1000,
				MinCapacity:          100,
				NumberOfApplications: -1,
				Ushers:               []string{"Alice Smith", "Bob Johnson"},
				Organizers: []data.Organizer{
					{
						ID:    primitive.NewObjectID(),
						Name:  "John Doe",
						Email: "johndoe@example.com",
						Phone: "123-456-7890",
						Role:  "Event Manager",
					},
				},
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
				Status:    "Open",
			},
			expectedStatus: http.StatusOK,
			expectedBody:`{"message":"Success"}`,
			
			setupMock: func(mockEventAppModel *MockEventAppModel, mockEventModel *MockEventModel, mockTokenExtractor *MockTokenExtractor) {
				// mockTokenExtractor.On("extractTokenData", mock.Anything).Return("test@example.com", true, true, nil)
				mockEventModel.On("UpdateEvent", mock.AnythingOfType("primitive.ObjectID"), mock.Anything).Return(&data.Event{}, nil)
				mockEventAppModel.On("GetEventApp", mock.Anything, mock.AnythingOfType("primitive.ObjectID")).Return(&data.EventApps{}, nil)
			},
		},
		{
			name:    "Invalid Id format",
			eventId:  "",
			eventData: data.Event{
				ID:   primitive.NewObjectID(),
				Date: time.Date(2025, 7, 15, 18, 0, 0, 0, time.UTC),
				Type: "CONFERENCE",
				Name: "Tech Conference 2025",
				Location: data.Location{
					Address: "123 Main Street",
					City:    "San Francisco",
					State:   "CA",
					Country: "USA",
				},
				MaxCapacity:          1000,
				MinCapacity:          100,
				NumberOfApplications: -1,
				Ushers:               []string{"Alice Smith", "Bob Johnson"},
				Organizers: []data.Organizer{
					{
						ID:    primitive.NewObjectID(),
						Name:  "John Doe",
						Email: "johndoe@example.com",
						Phone: "123-456-7890",
						Role:  "Event Manager",
					},
				},
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
				Status:    "Open",
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody:`{"error": "Invalid ID format"}`,
			
			setupMock: func(mockEventAppModel *MockEventAppModel, mockEventModel *MockEventModel, mockTokenExtractor *MockTokenExtractor) {
				
			},
		},
		{
			name:    "Error Updating Event",
			eventId:  primitive.NewObjectID().Hex(),
			eventData: data.Event{
				ID:   primitive.NewObjectID(),
				Date: time.Date(2025, 7, 15, 18, 0, 0, 0, time.UTC),
				Type: "CONFERENCE",
				Name: "Tech Conference 2025",
				Location: data.Location{
					Address: "123 Main Street",
					City:    "San Francisco",
					State:   "CA",
					Country: "USA",
				},
				MaxCapacity:          1000,
				MinCapacity:          100,
				NumberOfApplications: -1,
				Ushers:               []string{"Alice Smith", "Bob Johnson"},
				Organizers: []data.Organizer{
					{
						ID:    primitive.NewObjectID(),
						Name:  "John Doe",
						Email: "johndoe@example.com",
						Phone: "123-456-7890",
						Role:  "Event Manager",
					},
				},
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
				Status:    "Open",
			},
			expectedStatus: http.StatusInternalServerError,
			expectedBody:`{"error": "Failed to update event"}`,
			
			setupMock: func(mockEventAppModel *MockEventAppModel, mockEventModel *MockEventModel, mockTokenExtractor *MockTokenExtractor) {
				// mockTokenExtractor.On("extractTokenData", mock.Anything).Return("test@example.com", true, true, nil)
				mockEventModel.On("UpdateEvent", mock.AnythingOfType("primitive.ObjectID"), mock.AnythingOfType("*data.Event")).Return(&data.Event{}, errors.New("ERROR"))
				
			},
		},
		{
			name:    "Error fetching event apps",
			eventId:  primitive.NewObjectID().Hex(),
			eventData: data.Event{
				ID:   primitive.NewObjectID(),
				Date: time.Date(2025, 7, 15, 18, 0, 0, 0, time.UTC),
				Type: "CONFERENCE",
				Name: "Tech Conference 2025",
				Location: data.Location{
					Address: "123 Main Street",
					City:    "San Francisco",
					State:   "CA",
					Country: "USA",
				},
				MaxCapacity:          1000,
				MinCapacity:          100,
				NumberOfApplications: -1,
				Ushers:               []string{"Alice Smith", "Bob Johnson"},
				Organizers: []data.Organizer{
					{
						ID:    primitive.NewObjectID(),
						Name:  "John Doe",
						Email: "johndoe@example.com",
						Phone: "123-456-7890",
						Role:  "Event Manager",
					},
				},
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
				Status:    "Open",
			},
			expectedStatus: http.StatusInternalServerError,
			expectedBody:`{"error": "the server encountered a problem and could not process your request"}`,
			
			setupMock: func(mockEventAppModel *MockEventAppModel, mockEventModel *MockEventModel, mockTokenExtractor *MockTokenExtractor) {
			// mockTokenExtractor.On("extractTokenData", mock.Anything).Return("test@example.com", true, true, nil)
				mockEventModel.On("UpdateEvent", mock.AnythingOfType("primitive.ObjectID"), mock.AnythingOfType("*data.Event")).Return(&data.Event{},nil)
				mockEventAppModel.On("GetEventApp", mock.Anything, mock.AnythingOfType("primitive.ObjectID")).Return(&data.EventApps{}, data.ErrNoRecords)
				
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockEventAppModel := new(MockEventAppModel)
			mockEventModel := new(MockEventModel)
			mockTokenExtractor := new(MockTokenExtractor)

			app.models = data.Models{
				EventApps: mockEventAppModel,
				Event:     mockEventModel,
			}
			app.tokenExtractor = mockTokenExtractor

			tt.setupMock(mockEventAppModel, mockEventModel, mockTokenExtractor)



			
			url:="/v1/events/{id}"
			reqBody, _ := json.Marshal(tt.eventData)
			req := httptest.NewRequest(http.MethodPut, url,  bytes.NewBuffer(reqBody))
			req.SetPathValue("id", tt.eventId)
			
			
			
			req.Header.Set("Content-Type", "application/json")

			rr := httptest.NewRecorder()

			// Call the handler
			handler := http.HandlerFunc(app.updateEventHandler)
			handler.ServeHTTP(rr, req)

			// Assertions
			assert.Equal(t, tt.expectedStatus, rr.Code)
			assert.JSONEq(t, tt.expectedBody, rr.Body.String())

			// Verify mocks
			mockEventAppModel.AssertExpectations(t)
			mockEventModel.AssertExpectations(t)
			mockTokenExtractor.AssertExpectations(t)
		})
	}
}

// updatedEvent, err := app.models.Event.UpdateEvent(id, &event)
// 	if err != nil {
// 		app.Logger.Printf("Error updating event: %v", err)
// 		app.writeJSON(w, http.StatusInternalServerError, envelope{"error": "Failed to update event"}, nil)
// 		return
// 	}






// func TestApplyToEventHandler(t *testing.T) {
// 	mockEventAppModel := new(MockEventAppModel)
// 	mockEventModel := new(MockEventModel)
// 	mockTokenExtractor := new(MockTokenExtractor)

// 	rabbitConn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
// 	if err != nil {
// 		t.Fatalf("Failed to connect to RabbitMQ: %v", err)
// 	}

// 	app := &application{
// 		Logger: log.New(os.Stdout, "", 0),
// 		config: config{port: "80", env: "development"},
// 		models: data.Models{
// 			EventApps: mockEventAppModel,
// 			Event:     mockEventModel,
// 		},
// 		Rabbit:         rabbitConn,
// 		tokenExtractor: mockTokenExtractor,
// 	}

// 	tests := []struct {
// 		name           string
// 		eventApp       interface{}
// 		expectedStatus int
// 		expectedBody   string
// 		setupMock      func(mockEventAppModel *MockEventAppModel, mockEventModel *MockEventModel, mockTokenExtractor *MockTokenExtractor)
// 	}{
// 		{
// 			name: "Not Found Event Error",
// 			eventApp: struct {
// 				EventID string
// 			}{
// 				EventID: primitive.NewObjectID().Hex(),
// 			},
// 			expectedStatus: http.StatusNotFound,
// 			expectedBody:   `{"error":"Event app not found"}`,
// 			setupMock: func(mockEventAppModel *MockEventAppModel, mockEventModel *MockEventModel, mockTokenExtractor *MockTokenExtractor) {
// 				mockTokenExtractor.On("extractTokenData", mock.Anything).Return("test@example.com", true, true, nil)
// 				mockEventAppModel.On("GetEventApp", mock.Anything, mock.AnythingOfType("primitive.ObjectID")).Return(&data.EventApps{}, data.ErrNoRecords)
// 			},
// 		},
// 		{
// 			name: "Invalid Token",
// 			eventApp: struct {
// 				EventID string
// 			}{
// 				EventID: primitive.NewObjectID().Hex(),
// 			},
// 			expectedStatus: http.StatusUnauthorized,
// 			expectedBody:   `{"error":"Invalid token"}`,
// 			setupMock: func(mockEventAppModel *MockEventAppModel, mockEventModel *MockEventModel, mockTokenExtractor *MockTokenExtractor) {
// 				mockTokenExtractor.On("extractTokenData", mock.Anything).Return("", false, false, errors.New("Invalid token"))
// 			},
// 		},
// 		{
// 			name: "Already Applied to Event",
// 			eventApp: struct {
// 				EventID string
// 			}{
// 				EventID: primitive.NewObjectID().Hex(),
// 			},
// 			expectedStatus: http.StatusBadRequest,
// 			expectedBody:   `{"error":"You already applied to this event"}`,
// 			setupMock: func(mockEventAppModel *MockEventAppModel, mockEventModel *MockEventModel, mockTokenExtractor *MockTokenExtractor) {
// 				mockTokenExtractor.On("extractTokenData", mock.Anything).Return("test@example.com", true, true, nil)
// 				mockEventAppModel.On("GetEventApp", mock.Anything, mock.AnythingOfType("primitive.ObjectID")).Return(&data.EventApps{Attendee: []string{"test@example.com"}}, nil)
// 			},
// 		},
// 		{
// 			name: "Successful Event Application",
// 			eventApp: struct {
// 				EventID string
// 			}{
// 				EventID: primitive.NewObjectID().Hex(),
// 			},
// 			expectedStatus: http.StatusOK,
// 			expectedBody:   `{"message":"Applied to event successfully"}`,
// 			setupMock: func(mockEventAppModel *MockEventAppModel, mockEventModel *MockEventModel, mockTokenExtractor *MockTokenExtractor) {
// 				mockTokenExtractor.On("extractTokenData", mock.Anything).Return("test@example.com", true, true, nil)
// 				mockEventAppModel.On("GetEventApp", mock.Anything, mock.AnythingOfType("primitive.ObjectID")).Return(&data.EventApps{Attendee: []string{}}, nil)
// 				mockEventModel.On("GetEventByID", mock.AnythingOfType("primitive.ObjectID")).Return(&data.Event{Name: "Test Event", Date: time.Now().Add(1 * time.Hour), Location: data.Location{Address: "123 Test St", City: "Test City", State: "Test State", Country: "Test Country"}}, nil)
// 				mockEventAppModel.On("AddAttendeeToEvent", "test@example.com", mock.AnythingOfType("primitive.ObjectID")).Return(nil)
// 				mockEventAppModel.On("AddAttendeeToEvent", mock.Anything, mock.Anything).Return(nil)
// 			},
// 		},
// 		{
// 			name: "Event Has Ended",
// 			eventApp: struct {
// 				EventID string
// 			}{
// 				EventID: primitive.NewObjectID().Hex(),
// 			},
// 			expectedStatus: http.StatusBadRequest,
// 			expectedBody:   `{"error":"Event has ended"}`,
// 			setupMock: func(mockEventAppModel *MockEventAppModel, mockEventModel *MockEventModel, mockTokenExtractor *MockTokenExtractor) {
// 				mockTokenExtractor.On("extractTokenData", mock.Anything).Return("test@example.com", true, true, nil)
// 				mockEventAppModel.On("GetEventApp", mock.Anything, mock.AnythingOfType("primitive.ObjectID")).Return(&data.EventApps{Attendee: []string{}}, nil)
// 				mockEventModel.On("GetEventByID", mock.AnythingOfType("primitive.ObjectID")).Return(&data.Event{Name: "Test Event", Date: time.Now().Add(-1 * time.Hour), Location: data.Location{Address: "123 Test St", City: "Test City", State: "Test State", Country: "Test Country"}}, nil)
// 			},
// 		},
// 	}

// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			mockEventAppModel := new(MockEventAppModel)
// 			mockEventModel := new(MockEventModel)
// 			mockTokenExtractor := new(MockTokenExtractor)

// 			app.models = data.Models{
// 				EventApps: mockEventAppModel,
// 				Event:     mockEventModel,
// 			}
// 			app.tokenExtractor = mockTokenExtractor

// 			tt.setupMock(mockEventAppModel, mockEventModel, mockTokenExtractor)

// 			url := "/v1/events/{id}/apply"
// 			req := httptest.NewRequest(http.MethodPost, url, nil)
// 			req.SetPathValue("id", tt.eventApp.(struct{ EventID string }).EventID)
// 			req.Header.Set("Content-Type", "application/json")

// 			rr := httptest.NewRecorder()

// 			handler := http.HandlerFunc(app.applyToEventHandler)
// 			handler.ServeHTTP(rr, req)

// 			assert.Equal(t, tt.expectedStatus, rr.Code)
// 			assert.JSONEq(t, tt.expectedBody, rr.Body.String())

// 			mockEventAppModel.AssertExpectations(t)
// 			mockEventModel.AssertExpectations(t)
// 			mockTokenExtractor.AssertExpectations(t)
// 		})
// 	}
// }
