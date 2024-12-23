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

func TestApplyToEventHandler(t *testing.T) {
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
				mockTokenExtractor.On("extractTokenData", mock.Anything).Return("test@example.com", true, true, nil)
				mockEventAppModel.On("GetEventApp", mock.Anything, mock.AnythingOfType("primitive.ObjectID")).Return(&data.EventApps{}, data.ErrNoRecords)
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
				mockTokenExtractor.On("extractTokenData", mock.Anything).Return("", false, false, errors.New("Invalid token"))
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
				mockEventAppModel.On("GetEventApp", mock.Anything, mock.AnythingOfType("primitive.ObjectID")).Return(&data.EventApps{Attendee: []string{"test@example.com"}}, nil)
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
				mockTokenExtractor.On("extractTokenData", mock.Anything).Return("test@example.com", true, true, nil)
				mockEventAppModel.On("GetEventApp", mock.Anything, mock.AnythingOfType("primitive.ObjectID")).Return(&data.EventApps{Attendee: []string{}}, nil)
				mockEventModel.On("GetEventByID", mock.AnythingOfType("primitive.ObjectID")).Return(&data.Event{Name: "Test Event", Date: time.Now().Add(1 * time.Hour), Location: data.Location{Address: "123 Test St", City: "Test City", State: "Test State", Country: "Test Country"}}, nil)
				mockEventAppModel.On("AddAttendeeToEvent", "test@example.com", mock.AnythingOfType("primitive.ObjectID")).Return(nil)
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
				mockTokenExtractor.On("extractTokenData", mock.Anything).Return("test@example.com", true, true, nil)
				mockEventAppModel.On("GetEventApp", mock.Anything, mock.AnythingOfType("primitive.ObjectID")).Return(&data.EventApps{Attendee: []string{}}, nil)
				mockEventModel.On("GetEventByID", mock.AnythingOfType("primitive.ObjectID")).Return(&data.Event{Name: "Test Event", Date: time.Now().Add(-1 * time.Hour), Location: data.Location{Address: "123 Test St", City: "Test City", State: "Test State", Country: "Test Country"}}, nil)
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
