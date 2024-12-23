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
	
	"github.com/stretchr/testify/assert"
	amqp "github.com/rabbitmq/amqp091-go"
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

func TestGetEventByID(t *testing.T) {
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
		event       interface{}
		expectedStatus int
		expectedBody   string
		setupMock      func(mockEventAppModel *MockEventAppModel, mockEventModel *MockEventModel, mockTokenExtractor *MockTokenExtractor)
	}{
		{
			name: "Valid get event by ID",
			event: struct {
				EventID string
			}{
				EventID: "67473b35332e9a9361e03fef",
			},
			expectedStatus: http.StatusOK,
			expectedBody: `{
				"event": {
					"_id": "67473b35332e9a9361e03fef",
					"created_at": "2024-11-27T15:31:01.393Z",
					"date": "2024-07-15T18:00:00Z",
					"description": "",
					"location": {
						"address": "123 Main Street",
						"city": "San Francisco",
						"country": "USA",
						"state": "CA"
					},
					"max_capacity": 1000,
					"min_capacity": 100,
					"name": "Tech Conference 2024",
					"number_of_applications": 0,
					"organizers": [
						{
							"email": "john.doe@example.com",
							"id": "000000000000000000000000",
							"name": "John Doe",
							"phone": "123-456-7890",
							"role": "Lead Organizer"
						}
					],
					"status": "PENDING",
					"type": "CONFERENCE",
					"updated_at": "2024-11-27T15:31:01.393Z",
					"ushers": [
						"Alice Smith",
						"Bob Johnson"
					]
				}
			}`,
			setupMock: func(mockEventAppModel *MockEventAppModel, mockEventModel *MockEventModel, mockTokenExtractor *MockTokenExtractor) {
				id, _ := primitive.ObjectIDFromHex("67473b35332e9a9361e03fef")
				organizerId, _ := primitive.ObjectIDFromHex("000000000000000000000000")
				createdAt, _ := time.Parse(time.RFC3339, "2024-11-27T15:31:01.393Z")
				date, _ := time.Parse(time.RFC3339, "2024-07-15T18:00:00Z")
				updatedAt, _ := time.Parse(time.RFC3339, "2024-11-27T15:31:01.393Z")
				mockEventModel.On("GetEventByID", mock.AnythingOfType("primitive.ObjectID")).Return(&data.Event{
					ID:                 id,
					CreatedAt:          createdAt,
					Date:               date,
					Description:        "",
					Location:           data.Location{
						Address: "123 Main Street",
						City:    "San Francisco",
						Country: "USA",
						State:   "CA",
					},
					MaxCapacity:        1000,
					MinCapacity:        100,
					Name:               "Tech Conference 2024",
					NumberOfApplications: 0,
					Organizers: []data.Organizer{
						{
							Email: "john.doe@example.com",
							ID:    organizerId,
							Name:  "John Doe",
							Phone: "123-456-7890",
							Role:  "Lead Organizer",
						},
					},
					Status:    "PENDING",
					Type:      "CONFERENCE",
					UpdatedAt: updatedAt,
					Ushers:    []string{"Alice Smith", "Bob Johnson"},
				}, nil)
			},
		},
		{
			name: "Event not found",
			event: struct {
				EventID string
			}{
				EventID: "67473b35332e9a9361e03fef",
			},
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   `{"error": "Failed to fetch event"}`,
			setupMock: func(mockEventAppModel *MockEventAppModel, mockEventModel *MockEventModel, mockTokenExtractor *MockTokenExtractor) {
				mockEventModel.On("GetEventByID", mock.AnythingOfType("primitive.ObjectID")).Return(&data.Event{}, errors.New("Failed to fetch event"))
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
	
			req := httptest.NewRequest(http.MethodGet, url, nil)
			req.SetPathValue("id", tt.event.(struct{ EventID string }).EventID)
			req.Header.Set("Content-Type", "application/json")
	
			rr := httptest.NewRecorder()
	
			handler := http.HandlerFunc(app.getEventByIDHandler)
			handler.ServeHTTP(rr, req)
	
			assert.Equal(t, tt.expectedStatus, rr.Code)
			assert.JSONEq(t, tt.expectedBody, rr.Body.String())
	
			mockEventAppModel.AssertExpectations(t)
			mockEventModel.AssertExpectations(t)
			mockTokenExtractor.AssertExpectations(t)
		})
	}
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






