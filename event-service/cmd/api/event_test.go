package main

import (
	"github.com/MohamedHossam2004/Event-Planner/event-service/internal/data"
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
