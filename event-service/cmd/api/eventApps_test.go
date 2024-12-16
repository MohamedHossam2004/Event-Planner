package main

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

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

func TestCreateEventAppHandler(t *testing.T) {
	mockEventAppModel := new(MockEventAppModel)
	mockEventModel := new(MockEventModel)

	rabbitConn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		t.Fatalf("Failed to connect to RabbitMQ: %v", err)
	}

	app := &application{
		Logger: log.New(io.Discard, "", 0),
		config: config{port: "80", env: "development"},
		models: data.Models{
			EventApps: mockEventAppModel,
			Event:     mockEventModel,
		},
		Rabbit: rabbitConn,
	}

	tests := []struct {
		name           string
		eventApp       interface{}
		expectedStatus int
		setupMock      func(mockEventAppModel *MockEventAppModel, mockEventModel *MockEventModel)
	}{
		{
			name: "CreateEventApp Returns Error",
			eventApp: struct {
				ID       primitive.ObjectID
				EventID  primitive.ObjectID
				Attendee []string
			}{
				ID:       primitive.NewObjectID(),
				EventID:  primitive.NewObjectID(),
				Attendee: []string{"user1", "user2"},
			},
			expectedStatus: http.StatusInternalServerError,
			setupMock: func(mockEventAppModel *MockEventAppModel, mockEventModel *MockEventModel) {
				mockEventAppModel.On("CreateEventApp", mock.Anything, mock.AnythingOfType("*data.EventApps")).Return(errors.New("database error"))
			},
		},
		{
			name: "Valid Event App",
			eventApp: struct {
				ID       primitive.ObjectID
				EventID  primitive.ObjectID
				Attendee []string
			}{
				ID:       primitive.NewObjectID(),
				EventID:  primitive.NewObjectID(),
				Attendee: []string{"user1", "user2"},
			},
			expectedStatus: http.StatusCreated,
			setupMock: func(mockEventAppModel *MockEventAppModel, mockEventModel *MockEventModel) {
				mockEventAppModel.On("CreateEventApp", mock.Anything, mock.AnythingOfType("*data.EventApps")).Return(nil)
				mockEventModel.On("GetEventByID", mock.AnythingOfType("primitive.ObjectID")).Return(&data.Event{Name: "Test Event"}, nil)
			},
		},
		{
			name: "Missing Event ID",
			eventApp: struct {
				ID       primitive.ObjectID
				Attendee []string
			}{
				ID:       primitive.NewObjectID(),
				Attendee: []string{"user1", "user2"},
			},
			expectedStatus: http.StatusBadRequest,
			setupMock: func(mockEventAppModel *MockEventAppModel, mockEventModel *MockEventModel) {
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockEventAppModel := new(MockEventAppModel)
			mockEventModel := new(MockEventModel)

			app.models = data.Models{
				EventApps: mockEventAppModel,
				Event:     mockEventModel,
			}

			tt.setupMock(mockEventAppModel, mockEventModel)

			eventAppJSON, _ := json.Marshal(tt.eventApp)

			req := httptest.NewRequest(http.MethodPost, "/v1/eventapps", bytes.NewBuffer(eventAppJSON))
			req.Header.Set("Content-Type", "application/json")

			rr := httptest.NewRecorder()

			handler := http.HandlerFunc(app.createEventAppHandler)
			handler.ServeHTTP(rr, req)

			assert.Equal(t, tt.expectedStatus, rr.Code)

			mockEventAppModel.AssertExpectations(t)
			mockEventModel.AssertExpectations(t)
		})
	}
}
