// A server that handles notifications
package main

import (
	"log"
	"net/http"

	"context"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"io"
	"io/ioutil"
	"os"
	"os/signal"
	"time"
)

type NotifyEvent struct {
	Id      uuid.UUID `json:"id"`
	Message string    `json:"message"`
}

type NotifyEvents []NotifyEvent

func main() {
	host := "localhost"
	port := 20035

	addr := fmt.Sprintf("%s:%d", host, port)

	router := mux.NewRouter().StrictSlash(true)

	router.Methods("GET").Path("/").Name("RootIndex").HandlerFunc(RootIndex)
	router.Methods("POST").Path("/event").Name("EventPOST").HandlerFunc(NewEventNotification)
	router.Methods("GET").Path("/event").Name("EventGET").HandlerFunc(GetLatestEvent)
	router.Methods("GET").Path("/events").Name("EventsGET").HandlerFunc(GetRecentEvents)

	router.Use(loggingMiddleware)

	srv := &http.Server{
		Addr: addr,
		// Good practice to set timeouts to avoid Slowloris attacks.
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
		Handler:      router, // Pass our instance of gorilla/mux in.
	}

	// Run our server in a goroutine so that it doesn't block.
	go func() {
		log.Printf("Listening on %s...\n", addr)
		if err := srv.ListenAndServe(); err != nil {
			log.Println(err)
		}
	}()

	c := make(chan os.Signal, 1)
	// We'll accept graceful shutdowns when quit via SIGINT (Ctrl+C)
	// SIGKILL, SIGQUIT or SIGTERM (Ctrl+/) will not be caught.
	signal.Notify(c, os.Interrupt)

	// Block until we receive our signal.
	<-c

	// Create a deadline to wait for.
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*15)
	defer cancel()
	// Doesn't block if no connections, but will otherwise wait
	// until the timeout deadline.
	srv.Shutdown(ctx)
	// Optionally, you could run srv.Shutdown in a goroutine and block on
	// <-ctx.Done() if your application should wait for other services
	// to finalize based on context cancellation.
	log.Println("shutting down")
	os.Exit(0)
}

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		next.ServeHTTP(w, r)

		log.Printf(
			"%s\t%s\t%s",
			r.Method,
			r.RequestURI,
			time.Since(start),
		)
	})
}

func RootIndex(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=UTF-8")
	w.WriteHeader(http.StatusOK)

	fmt.Fprintln(w, "Hello World!")
	fmt.Fprintf(w, "We have %d events.\n", len(EVENTS))
}

var EVENTS = make(NotifyEvents, 0)

// Notify that a new event has just happened
func NewEventNotification(w http.ResponseWriter, r *http.Request) {
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
	event.Id = uuid.New()
	EVENTS = append(EVENTS, event)

	// Return it
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(event); err != nil {
		panic(err)
	}
}

func GetLatestEvent(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	if len(EVENTS) > 0 {
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(EVENTS[len(EVENTS)-1]); err != nil {
			panic(err)
		}
	} else {
		// No events yet
		w.WriteHeader(http.StatusNotFound)
	}
}

func GetRecentEvents(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)

	n := len(EVENTS)
	if n < 10 {
		if err := json.NewEncoder(w).Encode(EVENTS); err != nil {
			panic(err)
		}
	} else {
		if err := json.NewEncoder(w).Encode(EVENTS[n-10:]); err != nil {
			panic(err)
		}
	}
}
