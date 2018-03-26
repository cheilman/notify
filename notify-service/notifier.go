package main

import "fmt"

//
// Notify the user that something happened.
//

type Notifier interface {
	NotifyOfEvent(event NotifyEvent)
}

//
// Simple implementation
//

type StdoutNotifier struct{}

func NewStdoutNotifier() *StdoutNotifier {
	return &StdoutNotifier{}
}

func (s *StdoutNotifier) NotifyOfEvent(event NotifyEvent) {
	fmt.Printf("*** NEW EVENT: %v ***\n", event)
}
