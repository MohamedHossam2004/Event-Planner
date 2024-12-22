package main

import (
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/MohamedHossam2004/Event-Planner/event-service/internal/data"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson/primitive"
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

func TestDeleteEvent(t *testing.T) {
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

			name: "Valid delete case",
			eventApp: struct {
				EventID string
			}{
				EventID: primitive.NewObjectID().Hex(),
			},
			expectedStatus: http.StatusOK,
			expectedBody:   `{"message":"Event deleted successfully"}`,
			setupMock: func(mockEventAppModel *MockEventAppModel, mockEventModel *MockEventModel, mockTokenExtractor *MockTokenExtractor) {
				mockEventModel.On("GetEventByID", mock.AnythingOfType("primitive.ObjectID")).Return(&data.Event{}, nil)
				mockEventModel.On("DeleteEvent", mock.Anything).Return(nil)

			},
		},
		{

			name: "Invalid ID Format",
			eventApp: struct {
				EventID string
			}{
				EventID: "",
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `{"error":"Invalid ID format"}`,
			setupMock: func(mockEventAppModel *MockEventAppModel, mockEventModel *MockEventModel, mockTokenExtractor *MockTokenExtractor) {

			},
		},
		{

			name: "Event Not Found Error",
			eventApp: struct {
				EventID string
			}{
				EventID: primitive.NewObjectID().Hex(),
			},
			expectedStatus: http.StatusNotFound,
			expectedBody:   `{"error":"Failed to fetch event"}`,
			setupMock: func(mockEventAppModel *MockEventAppModel, mockEventModel *MockEventModel, mockTokenExtractor *MockTokenExtractor) {
				mockEventModel.On("GetEventByID", mock.AnythingOfType("primitive.ObjectID")).Return(&data.Event{}, data.ErrNoRecords)
			},
		},
		{

			name: "Error Deleting Event",
			eventApp: struct {
				EventID string
			}{
				EventID: primitive.NewObjectID().Hex(),
			},
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   `{"error":"Failed to delete event"}`,
			setupMock: func(mockEventAppModel *MockEventAppModel, mockEventModel *MockEventModel, mockTokenExtractor *MockTokenExtractor) {
				mockEventModel.On("GetEventByID", mock.AnythingOfType("primitive.ObjectID")).Return(&data.Event{}, nil)
				mockEventModel.On("DeleteEvent", mock.AnythingOfType("primitive.ObjectID")).Return(data.ErrNoRecords)
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

			url := "/v1/events/{id}"
			req := httptest.NewRequest(http.MethodDelete, url, nil)
			req.SetPathValue("id", tt.eventApp.(struct{ EventID string }).EventID)
			req.Header.Set("Content-Type", "application/json")

			rr := httptest.NewRecorder()

			handler := http.HandlerFunc(app.deleteEventHandler)
			handler.ServeHTTP(rr, req)

			assert.Equal(t, tt.expectedStatus, rr.Code)
			assert.JSONEq(t, tt.expectedBody, rr.Body.String())

			mockEventAppModel.AssertExpectations(t)
			mockEventModel.AssertExpectations(t)
			mockTokenExtractor.AssertExpectations(t)
		})
	}
}
