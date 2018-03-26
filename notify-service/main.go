// A server that handles notifications
package main

import (
	"log"
	"net/http"

	"context"
	"fmt"
	"github.com/gorilla/mux"
	"os"
	"os/signal"
	"time"
)

var EVENTS = NewInMemoryNotifyStorage()
var NOTIFIER = NewDesktopNotifier()

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
}
