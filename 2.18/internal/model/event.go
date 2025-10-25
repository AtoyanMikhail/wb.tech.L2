package model

import "time"

// Event represents a calendar event
type Event struct {
	ID        int       `json:"id"`
	UserID    int       `json:"user_id"`
	Date      time.Time `json:"date"`
	EventText string    `json:"event"`
}

// CreateEventRequest is a request structure for creating an event
type CreateEventRequest struct {
	UserID    int    `json:"user_id"`
	Date      string `json:"date"`
	EventText string `json:"event"`
}

// UpdateEventRequest is a request structure for updating an event
type UpdateEventRequest struct {
	ID        int    `json:"id"`
	UserID    int    `json:"user_id"`
	Date      string `json:"date"`
	EventText string `json:"event"`
}

// DeleteEventRequest is a request structure for deleting an event
type DeleteEventRequest struct {
	ID int `json:"id"`
}

// Response is a standard server response
type Response struct {
	Result interface{} `json:"result,omitempty"`
	Error  string      `json:"error,omitempty"`
}
