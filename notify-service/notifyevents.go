package main

import (
	"github.com/google/uuid"
	"sync"
)

//
// Handle events and APIs
//

type NotifyEvent struct {
	Id      uuid.UUID `json:"id"`
	Message string    `json:"message"`
}

type NotifyEvents []NotifyEvent

type NotifyStorage interface {
	Add(event NotifyEvent) NotifyEvent
	GetLatest() NotifyEvent
	GetNRecent(n int) NotifyEvents
}

//
// In-memory storage
//

type InMemoryNotifyStorage struct {
	events NotifyEvents
	mutex  *sync.Mutex
}

func NewInMemoryNotifyStorage() *InMemoryNotifyStorage {
	return &InMemoryNotifyStorage{
		events: make(NotifyEvents, 0),
		mutex:  &sync.Mutex{},
	}
}

func (s *InMemoryNotifyStorage) Add(event NotifyEvent) NotifyEvent {
	// Add it to our list
	event.Id = uuid.New()

	s.mutex.Lock()
	s.events = append(s.events, event)
	// TODO: At some point, remove elements to keep the size of this array down
	s.mutex.Unlock()

	return event
}

func (s *InMemoryNotifyStorage) GetLatest() *NotifyEvent {
	if len(s.events) > 0 {
		evt := s.events[len(s.events)-1]
		return &evt
	} else {
		// No events yet
		return nil
	}
}

func (s *InMemoryNotifyStorage) GetNRecent(n int) NotifyEvents {
	size := len(s.events)
	if size < n {
		return s.events[:]
	} else {
		return s.events[size-n:]
	}
}
