package main

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"net/http"
)

// Notify that a new event has just happened
func NewEventNotification(w http.ResponseWriter, r *http.Request,
	events NotifyStorage, notifyChannel chan<- NotifyEvent) {

	var event NotifyEvent

	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
	if err != nil {
		panic(err)
	}
	if err := r.Body.Close(); err != nil {
		panic(err)
	}
	if err := json.Unmarshal(body, &event); err != nil {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusUnprocessableEntity)
		if err := json.NewEncoder(w).Encode(err); err != nil {
			panic(err)
		}
	}

	// Add it to our list
	event = events.Add(event)

	// Notify
	log.Printf("About to dump event to channel.")
	notifyChannel <- event
	log.Printf("Dumped event.")
	//NOTIFIER.NotifyOfEvent(event)

	// Return it
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(event); err != nil {
		panic(err)
	}
}

func GetLatestEvent(w http.ResponseWriter, r *http.Request,
	events NotifyStorage) {

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	evt := events.GetLatest()

	if evt != nil {
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(evt); err != nil {
			panic(err)
		}
	} else {
		// No events yet
		w.WriteHeader(http.StatusNotFound)
	}
}

func GetRecentEvents(w http.ResponseWriter, r *http.Request,
	events NotifyStorage) {

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(events.GetNRecent(10)); err != nil {
		panic(err)
	}
}
