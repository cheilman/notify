package main

import (
	"github.com/google/uuid"
	"sync"
)

//
// Handle events and APIs
//

type AlertLevel int

//go:generate stringer -type=Pill
const (
	ERROR       AlertLevel = iota // = 0
	WARNING                       // = 1
	INFORMATION                   // = 2
)

type NotifyEvent struct {
	Id              uuid.UUID  `json:"id"`
	OriginatingHost string     `json:"host"`
	Title           string     `json:"title"`
	Message         string     `json:"message"`
	Icon            string     `json:"icon"`
	Category        string     `json:"category"`
	SubCategory     string     `json:"subcategory"`
	Level           AlertLevel `json:"level"`
}

type NotifyEvents []NotifyEvent

func (e *NotifyEvent) isError() bool {
	return e.Level == ERROR
}

//
// Store stuff
//

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
	event.Icon = "/home/local/ANT/heilmanc/go/src/github.com/cheilman/notify/notify-service/info.png"

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
