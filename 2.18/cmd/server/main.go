package main

import (
	"calendar/internal/config"
	"calendar/internal/handler"
	"calendar/internal/middleware"
	"calendar/internal/service"
	"log"
	"net/http"
)

func main() {
	cfg := config.Load()

	eventService := service.NewEventService()

	eventHandler := handler.NewEventHandler(eventService)

	mux := http.NewServeMux()

	mux.HandleFunc("/create_event", eventHandler.CreateEvent)
	mux.HandleFunc("/update_event", eventHandler.UpdateEvent)
	mux.HandleFunc("/delete_event", eventHandler.DeleteEvent)
	mux.HandleFunc("/events_for_day", eventHandler.GetEventsForDay)
	mux.HandleFunc("/events_for_week", eventHandler.GetEventsForWeek)
	mux.HandleFunc("/events_for_month", eventHandler.GetEventsForMonth)

	loggedMux := middleware.Logger(mux)

	address := ":" + cfg.Port
	log.Printf("Starting server on %s", address)
	if err := http.ListenAndServe(address, loggedMux); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
