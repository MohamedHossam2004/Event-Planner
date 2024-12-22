package main

import (
	"bytes"
	"encoding/json"
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

func TestCreateEvent(t *testing.T) {
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
		event          interface{}
		ID             primitive.ObjectID
		expectedStatus int
		expectedBody   string
		setupMock      func(mockEventAppModel *MockEventAppModel, mockEventModel *MockEventModel, mockTokenExtractor *MockTokenExtractor)
	}{
		{
			name: "Event not created",
			event: struct {
				Name        string
				Description string
				Location    data.Location
				Date        time.Time
				Type        string
			}{
				Name:        "Test Event",
				Description: "Test Description",
				Location: data.Location{
					Address: "123 Main St",
					City:    "New York",
					State:   "NY",
					Country: "USA",
				},
				Date: time.Now(),
				Type: "SOCIAL",
			},
			ID:             primitive.NewObjectID(),
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   `{"error": "Failed to create event"}`,
			setupMock: func(mockEventAppModel *MockEventAppModel, mockEventModel *MockEventModel, mockTokenExtractor *MockTokenExtractor) {
				mockEventModel.On("CreateEvent", mock.AnythingOfType("*data.Event")).Return(&data.Event{}, errors.New("something went wrong"))
			},
		},
		{
			name: "Event already exists",
			event: struct {
				Name        string
				Description string
				Location    data.Location
				Date        time.Time
				Type        string
			}{
				Name:        "Test Event",
				Description: "Test Description",
				Location: data.Location{
					Address: "123 Main St",
					City:    "New York",
					State:   "NY",
					Country: "USA",
				},
				Date: time.Now(),
				Type: "SOCIAL",
			},
			ID:             primitive.NewObjectID(),
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   `{"error": "Failed to create event"}`,
			setupMock: func(mockEventAppModel *MockEventAppModel, mockEventModel *MockEventModel, mockTokenExtractor *MockTokenExtractor) {
				mockEventModel.On("CreateEvent", mock.AnythingOfType("*data.Event")).Return(&data.Event{}, errors.New("event already exists"))
			},
		},
		{
			name: "Event created",
			event: struct {
				Name        string
				Description string
				Location    data.Location
				Date        time.Time
				Type        string
			}{
				Name:        "Test Event",
				Description: "Test Description",
				Location: data.Location{
					Address: "123 Main St",
					City:    "New York",
					State:   "NY",
					Country: "USA",
				},
				Date: time.Now(),
				Type: "SOCIAL",
			},
			ID:             primitive.NewObjectID(),
			expectedStatus: http.StatusCreated,
			expectedBody: `{
        "event": {
            "_id": "67689dee645ad5b4bbbe122d",
            "created_at": "0001-01-01T00:00:00Z",
            "date": "2024-12-23T01:17:02Z",
            "description": "Test Description",
            "location": {
                "address": "123 Main St",
                "city": "New York",
                "country": "USA",
                "state": "NY"
            },
            "max_capacity": 0,
            "min_capacity": 0,
            "name": "Test Event",
            "number_of_applications": 0,
            "organizers": null,
            "status": "",
            "type": "SOCIAL",
            "updated_at": "0001-01-01T00:00:00Z",
            "ushers": null
        }
    }`,
			setupMock: func(mockEventAppModel *MockEventAppModel, mockEventModel *MockEventModel, mockTokenExtractor *MockTokenExtractor) {
				id, _ := primitive.ObjectIDFromHex("67689dee645ad5b4bbbe122d")
				mockEventModel.On("CreateEvent", mock.AnythingOfType("*data.Event")).Return(&data.Event{
					ID:          id,
					Name:        "Test Event",
					Description: "Test Description",
					Location: data.Location{
						Address: "123 Main St",
						City:    "New York",
						State:   "NY",
						Country: "USA",
					},
					Date: time.Date(2024, time.December, 23, 1, 17, 2, 0, time.UTC),
					Type: "SOCIAL",
				}, nil)
				mockEventAppModel.On("CreateEventApp", mock.Anything, mock.AnythingOfType("*data.EventApps")).Return(nil)
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

			url := "/v1/events"

			body, err := json.Marshal(tt.event)
			if err != nil {
				t.Fatal(err)
			}
			req := httptest.NewRequest(http.MethodPost, url, bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")

			rr := httptest.NewRecorder()

			handler := http.HandlerFunc(app.createEventHandler)
			handler.ServeHTTP(rr, req)

			assert.Equal(t, tt.expectedStatus, rr.Code)
			assert.JSONEq(t, tt.expectedBody, rr.Body.String())

			mockEventAppModel.AssertExpectations(t)
			mockEventModel.AssertExpectations(t)
			mockTokenExtractor.AssertExpectations(t)
		})
	}
}
