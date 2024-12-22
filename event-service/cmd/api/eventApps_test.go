package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/MohamedHossam2004/Event-Planner/event-service/internal/data"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MockEventAppModel struct {
	mock.Mock
}

func (m *MockEventAppModel) CreateEventApp(ctx context.Context, eventApp *data.EventApps) error {
	args := m.Called(ctx, eventApp)
	return args.Error(0)
}

func (m *MockEventAppModel) GetEventApp(ctx context.Context, id primitive.ObjectID) (*data.EventApps, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*data.EventApps), args.Error(1)
}

func (m *MockEventAppModel) UpdateEventApp(ctx context.Context, id primitive.ObjectID, update bson.M) error {
	args := m.Called(ctx, id, update)
	return args.Error(0)
}

func (m *MockEventAppModel) DeleteEventApp(ctx context.Context, id primitive.ObjectID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockEventAppModel) ListEventApps(ctx context.Context, filter bson.M, opts *options.FindOptions) ([]*data.EventApps, error) {
	args := m.Called(ctx, filter, opts)
	return args.Get(0).([]*data.EventApps), args.Error(1)
}

func (m *MockEventAppModel) AddAttendeeToEvent(name string, eventId primitive.ObjectID) error {
	args := m.Called(name, eventId)
	return args.Error(0)
}

func (m *MockEventAppModel) RemoveAttendeeFromEvent(name string, eventId primitive.ObjectID) error {
	args := m.Called(name, eventId)
	return args.Error(0)
}

func (m *MockEventAppModel) GetEventsByUserEmail(email string) ([]*data.Event, error) {
	args := m.Called(email)
	return args.Get(0).([]*data.Event), args.Error(1)
}

type MockTokenExtractor struct {
	mock.Mock
}

func (m *MockTokenExtractor) extractTokenData(r *http.Request) (string, bool, bool, error) {
	args := m.Called(r)
	return args.String(0), args.Bool(1), args.Bool(2), args.Error(3)
}

// func TestCreateEventAppHandler(t *testing.T) {
// 	mockEventAppModel := new(MockEventAppModel)
// 	mockEventModel := new(MockEventModel)

// 	rabbitConn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
// 	if err != nil {
// 		t.Fatalf("Failed to connect to RabbitMQ: %v", err)
// 	}

// 	app := &application{
// 		Logger: log.New(io.Discard, "", 0),
// 		config: config{port: "80", env: "development"},
// 		models: data.Models{
// 			EventApps: mockEventAppModel,
// 			Event:     mockEventModel,
// 		},
// 		Rabbit: rabbitConn,
// 	}

// 	tests := []struct {
// 		name           string
// 		eventApp       interface{}
// 		expectedStatus int
// 		setupMock      func(mockEventAppModel *MockEventAppModel, mockEventModel *MockEventModel)
// 	}{
// 		{
// 			name: "CreateEventApp Returns Error",
// 			eventApp: struct {
// 				ID       primitive.ObjectID
// 				EventID  primitive.ObjectID
// 				Attendee []string
// 			}{
// 				ID:       primitive.NewObjectID(),
// 				EventID:  primitive.NewObjectID(),
// 				Attendee: []string{"user1", "user2"},
// 			},
// 			expectedStatus: http.StatusInternalServerError,
// 			setupMock: func(mockEventAppModel *MockEventAppModel, mockEventModel *MockEventModel) {
// 				mockEventAppModel.On("CreateEventApp", mock.Anything, mock.AnythingOfType("*data.EventApps")).Return(errors.New("database error"))
// 			},
// 		},
// 		{
// 			name: "Valid Event App",
// 			eventApp: struct {
// 				ID       primitive.ObjectID
// 				EventID  primitive.ObjectID
// 				Attendee []string
// 			}{
// 				ID:       primitive.NewObjectID(),
// 				EventID:  primitive.NewObjectID(),
// 				Attendee: []string{"user1", "user2"},
// 			},
// 			expectedStatus: http.StatusCreated,
// 			setupMock: func(mockEventAppModel *MockEventAppModel, mockEventModel *MockEventModel) {
// 				mockEventAppModel.On("CreateEventApp", mock.Anything, mock.AnythingOfType("*data.EventApps")).Return(nil)
// 				mockEventModel.On("GetEventByID", mock.AnythingOfType("primitive.ObjectID")).Return(&data.Event{Name: "Test Event"}, nil)
// 			},
// 		},
// 		{
// 			name: "Missing Event ID",
// 			eventApp: struct {
// 				ID       primitive.ObjectID
// 				Attendee []string
// 			}{
// 				ID:       primitive.NewObjectID(),
// 				Attendee: []string{"user1", "user2"},
// 			},
// 			expectedStatus: http.StatusBadRequest,
// 			setupMock: func(mockEventAppModel *MockEventAppModel, mockEventModel *MockEventModel) {
// 			},
// 		},
// 	}

// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			mockEventAppModel := new(MockEventAppModel)
// 			mockEventModel := new(MockEventModel)

// 			app.models = data.Models{
// 				EventApps: mockEventAppModel,
// 				Event:     mockEventModel,
// 			}

// 			tt.setupMock(mockEventAppModel, mockEventModel)

// 			eventAppJSON, _ := json.Marshal(tt.eventApp)

// 			req := httptest.NewRequest(http.MethodPost, "/v1/eventapps", bytes.NewBuffer(eventAppJSON))
// 			req.Header.Set("Content-Type", "application/json")

// 			rr := httptest.NewRecorder()

// 			handler := http.HandlerFunc(app.createEventAppHandler)
// 			handler.ServeHTTP(rr, req)

// 			assert.Equal(t, tt.expectedStatus, rr.Code)

// 			mockEventAppModel.AssertExpectations(t)
// 			mockEventModel.AssertExpectations(t)
// 		})
// 	}
// }

func TestApplyToEventHandler(t *testing.T) {
	mockEventAppModel := new(MockEventAppModel)
	mockEventModel := new(MockEventModel)
	mockTokenExtractor := new(MockTokenExtractor) // Mock token extractor

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
		tokenExtractor: mockTokenExtractor, // Use the mock token extractor here
	}

	tests := []struct {
		name           string
		eventApp       interface{}
		expectedStatus int
		expectedBody   string
		setupMock      func(mockEventAppModel *MockEventAppModel, mockEventModel *MockEventModel, mockTokenExtractor *MockTokenExtractor)
	}{
		{
			name: "Not Found Event Error",
			eventApp: struct {
				EventID string
			}{
				EventID: primitive.NewObjectID().Hex(),
			},
			expectedStatus: http.StatusNotFound,
			expectedBody:   `{"error":"Event app not found"}`,
			setupMock: func(mockEventAppModel *MockEventAppModel, mockEventModel *MockEventModel, mockTokenExtractor *MockTokenExtractor) {
				mockTokenExtractor.On("extractTokenData", mock.Anything).Return("test@example.com", true, true, nil)                                       // Mock successful token extraction
				mockEventAppModel.On("GetEventApp", mock.Anything, mock.AnythingOfType("primitive.ObjectID")).Return(&data.EventApps{}, data.ErrNoRecords) // Simulate event app not found
			},
		},
		{
			name: "Invalid Token",
			eventApp: struct {
				EventID string
			}{
				EventID: primitive.NewObjectID().Hex(),
			},
			expectedStatus: http.StatusUnauthorized,
			expectedBody:   `{"error":"Invalid token"}`,
			setupMock: func(mockEventAppModel *MockEventAppModel, mockEventModel *MockEventModel, mockTokenExtractor *MockTokenExtractor) {
				mockTokenExtractor.On("extractTokenData", mock.Anything).Return("", false, false, errors.New("Invalid token")) // Simulate invalid token
			},
		},
		{
			name: "Already Applied to Event",
			eventApp: struct {
				EventID string
			}{
				EventID: primitive.NewObjectID().Hex(),
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `{"error":"You already applied to this event"}`,
			setupMock: func(mockEventAppModel *MockEventAppModel, mockEventModel *MockEventModel, mockTokenExtractor *MockTokenExtractor) {
				mockTokenExtractor.On("extractTokenData", mock.Anything).Return("test@example.com", true, true, nil)
				mockEventAppModel.On("GetEventApp", mock.Anything, mock.AnythingOfType("primitive.ObjectID")).Return(&data.EventApps{Attendee: []string{"test@example.com"}}, nil) // Simulate the user already applied
			},
		},
		{
			name: "Successful Event Application",
			eventApp: struct {
				EventID string
			}{
				EventID: primitive.NewObjectID().Hex(),
			},
			expectedStatus: http.StatusOK,
			expectedBody:   `{"message":"Applied to event successfully"}`,
			setupMock: func(mockEventAppModel *MockEventAppModel, mockEventModel *MockEventModel, mockTokenExtractor *MockTokenExtractor) {
				mockTokenExtractor.On("extractTokenData", mock.Anything).Return("test@example.com", true, true, nil)                                                                                                                                                                             // Mock successful token extraction
				mockEventAppModel.On("GetEventApp", mock.Anything, mock.AnythingOfType("primitive.ObjectID")).Return(&data.EventApps{Attendee: []string{}}, nil)                                                                                                                                 // Simulate no prior event application
				mockEventModel.On("GetEventByID", mock.AnythingOfType("primitive.ObjectID")).Return(&data.Event{Name: "Test Event", Date: time.Now().Add(1 * time.Hour), Location: data.Location{Address: "123 Test St", City: "Test City", State: "Test State", Country: "Test Country"}}, nil) // Simulate a valid event
				mockEventAppModel.On("AddAttendeeToEvent", "test@example.com", mock.AnythingOfType("primitive.ObjectID")).Return(nil)                                                                                                                                                            // Simulate adding attendee to event
				mockEventAppModel.On("AddAttendeeToEvent", mock.Anything, mock.Anything).Return(nil)
			},
		},
		{
			name: "Event Has Ended",
			eventApp: struct {
				EventID string
			}{
				EventID: primitive.NewObjectID().Hex(),
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `{"error":"Event has ended"}`,
			setupMock: func(mockEventAppModel *MockEventAppModel, mockEventModel *MockEventModel, mockTokenExtractor *MockTokenExtractor) {
				mockTokenExtractor.On("extractTokenData", mock.Anything).Return("test@example.com", true, true, nil)                                                                                                                                                                              // Mock successful token extraction
				mockEventAppModel.On("GetEventApp", mock.Anything, mock.AnythingOfType("primitive.ObjectID")).Return(&data.EventApps{Attendee: []string{}}, nil)                                                                                                                                  // Simulate no prior event application
				mockEventModel.On("GetEventByID", mock.AnythingOfType("primitive.ObjectID")).Return(&data.Event{Name: "Test Event", Date: time.Now().Add(-1 * time.Hour), Location: data.Location{Address: "123 Test St", City: "Test City", State: "Test State", Country: "Test Country"}}, nil) // Simulate past event
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

			url := "/v1/events/{id}/apply"
			req := httptest.NewRequest(http.MethodPost, url, nil)
			req.SetPathValue("id", tt.eventApp.(struct{ EventID string }).EventID)
			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("Authorization", "Bearer eyJhbGciOiJIUzI1NiJ9.eyJhdWQiOlsiZ2l1LWV2ZW50LWh1Yi5jb20iXSwiZW1haWwiOiJtQG0uY29tIiwiZXhwIjoxNzM0OTc3MjEzLjc2MjA1MywiaWF0IjoxNzM0ODkwODEzLjc2MjA1MiwiaXNBY3RpdmF0ZWQiOmZhbHNlLCJpc0FkbWluIjpmYWxzZSwiaXNzIjoiZ2l1LWV2ZW50LWh1Yi5jb20iLCJuYW1lIjoiTW9oYXJyYW0iLCJuYmYiOjE3MzQ4OTA4MTMuNzYyMDUzLCJzdWIiOiI2NyJ9.6NLXsUH4PxjvpR_OVyZjJElo8mHllaJm4yGPb96Fe0Q")

			rr := httptest.NewRecorder()

			handler := http.HandlerFunc(app.applyToEventHandler)
			handler.ServeHTTP(rr, req)

			assert.Equal(t, tt.expectedStatus, rr.Code)
			assert.JSONEq(t, tt.expectedBody, rr.Body.String())

			mockEventAppModel.AssertExpectations(t)
			mockEventModel.AssertExpectations(t)
			mockTokenExtractor.AssertExpectations(t)
		})
	}
}
