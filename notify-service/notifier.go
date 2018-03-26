package main

import (
	"fmt"
	"github.com/gen2brain/beeep"
	"log"
	"os/exec"
)

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

//
// Fancier implementation
//

type DesktopNotifier struct{}

func NewDesktopNotifier() *DesktopNotifier {
	return &DesktopNotifier{}
}

func (s *DesktopNotifier) NotifyOfEvent(event NotifyEvent) {
	log.Printf("Hey, this shit got logged.")

	var err error = nil
	if event.isError() {
		log.Printf("Alerting.")
		err = beeep.Alert(event.Title, event.Message, event.Icon)
		log.Printf("Alerted: '%v'", err)
	} else {
		log.Printf("Notifying.")
		err = beeep.Notify(event.Title, event.Message, event.Icon)
		log.Printf("Notified: '%v'", err)
	}

	if err != nil {
		log.Printf("Everything is hosed! %v", err)
		panic(err)
	} else {
		log.Printf("Everything is fine.")
	}
}

//
// Dumbest implementation
//

type NotifySendNotifier struct{}

func NewNotifySendNotifier() *NotifySendNotifier {
	return &NotifySendNotifier{}
}

func alertLevelToUrgency(level AlertLevel) string {
	switch level {
	case ERROR:
		return "critical"
	case INFORMATION:
		return "low"
	case WARNING:
		return "normal"
	}

	return "normal"
}

func (s *NotifySendNotifier) NotifyOfEvent(event NotifyEvent) {
	log.Printf("Hey, this shit got logged.")

	send, err := exec.LookPath("notify-send")
	if err != nil {
		panic(err)
	}

	urgencyParam := fmt.Sprintf("--urgency=%s", alertLevelToUrgency(event.Level))

	var c *exec.Cmd = nil
	if event.Icon != "" {
		c = exec.Command(send, urgencyParam, fmt.Sprintf("--icon=%s", event.Icon), event.Title, event.Message)
	} else {
		c = exec.Command(send, urgencyParam, event.Title, event.Message)
	}

	err = c.Start()

	if err != nil {
		log.Printf("Everything is hosed! %v", err)
		panic(err)
	} else {
		log.Printf("Everything is fine.")
	}
}
