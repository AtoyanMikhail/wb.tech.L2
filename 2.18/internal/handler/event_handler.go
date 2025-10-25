package handler

import (
	"calendar/internal/model"
	"calendar/internal/service"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
)

// EventHandler handles HTTP requests to the events API
type EventHandler struct {
	service *service.EventService
}

// NewEventHandler creates a new event handler
func NewEventHandler(service *service.EventService) *EventHandler {
	return &EventHandler{
		service: service,
	}
}

// CreateEvent handles POST /create_event
func (h *EventHandler) CreateEvent(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		sendError(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req model.CreateEventRequest
	if err := h.parseRequest(r, &req); err != nil {
		sendError(w, err.Error(), http.StatusBadRequest)
		return
	}

	event, err := h.service.CreateEvent(req.UserID, req.Date, req.EventText)
	if err != nil {
		h.handleServiceError(w, err)
		return
	}

	sendSuccess(w, event, http.StatusOK)
}

// UpdateEvent handles POST /update_event
func (h *EventHandler) UpdateEvent(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		sendError(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req model.UpdateEventRequest
	if err := h.parseRequest(r, &req); err != nil {
		sendError(w, err.Error(), http.StatusBadRequest)
		return
	}

	event, err := h.service.UpdateEvent(req.ID, req.UserID, req.Date, req.EventText)
	if err != nil {
		h.handleServiceError(w, err)
		return
	}

	sendSuccess(w, event, http.StatusOK)
}

// DeleteEvent handles POST /delete_event
func (h *EventHandler) DeleteEvent(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		sendError(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req model.DeleteEventRequest
	if err := h.parseRequest(r, &req); err != nil {
		sendError(w, err.Error(), http.StatusBadRequest)
		return
	}

	err := h.service.DeleteEvent(req.ID)
	if err != nil {
		h.handleServiceError(w, err)
		return
	}

	sendSuccess(w, "event deleted successfully", http.StatusOK)
}

// GetEventsForDay handles GET /events_for_day
func (h *EventHandler) GetEventsForDay(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		sendError(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	userID, date, err := h.parseQueryParams(r)
	if err != nil {
		sendError(w, err.Error(), http.StatusBadRequest)
		return
	}

	events, err := h.service.GetEventsForDay(userID, date)
	if err != nil {
		h.handleServiceError(w, err)
		return
	}

	sendSuccess(w, events, http.StatusOK)
}

// GetEventsForWeek handles GET /events_for_week
func (h *EventHandler) GetEventsForWeek(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		sendError(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	userID, date, err := h.parseQueryParams(r)
	if err != nil {
		sendError(w, err.Error(), http.StatusBadRequest)
		return
	}

	events, err := h.service.GetEventsForWeek(userID, date)
	if err != nil {
		h.handleServiceError(w, err)
		return
	}

	sendSuccess(w, events, http.StatusOK)
}

// GetEventsForMonth handles GET /events_for_month
func (h *EventHandler) GetEventsForMonth(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		sendError(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	userID, date, err := h.parseQueryParams(r)
	if err != nil {
		sendError(w, err.Error(), http.StatusBadRequest)
		return
	}

	events, err := h.service.GetEventsForMonth(userID, date)
	if err != nil {
		h.handleServiceError(w, err)
		return
	}

	sendSuccess(w, events, http.StatusOK)
}

func (h *EventHandler) parseRequest(r *http.Request, v interface{}) error {
	contentType := r.Header.Get("Content-Type")

	if contentType == "application/json" {
		if err := json.NewDecoder(r.Body).Decode(v); err != nil {
			return errors.New("invalid JSON format")
		}
		return nil
	}

	if err := r.ParseForm(); err != nil {
		return errors.New("failed to parse form data")
	}

	switch req := v.(type) {
	case *model.CreateEventRequest:
		userID, err := strconv.Atoi(r.FormValue("user_id"))
		if err != nil {
			return errors.New("invalid user_id")
		}
		req.UserID = userID
		req.Date = r.FormValue("date")
		req.EventText = r.FormValue("event")

	case *model.UpdateEventRequest:
		id, err := strconv.Atoi(r.FormValue("id"))
		if err != nil {
			return errors.New("invalid id")
		}
		userID, err := strconv.Atoi(r.FormValue("user_id"))
		if err != nil {
			return errors.New("invalid user_id")
		}
		req.ID = id
		req.UserID = userID
		req.Date = r.FormValue("date")
		req.EventText = r.FormValue("event")

	case *model.DeleteEventRequest:
		id, err := strconv.Atoi(r.FormValue("id"))
		if err != nil {
			return errors.New("invalid id")
		}
		req.ID = id

	default:
		return errors.New("unsupported request type")
	}

	return nil
}

func (h *EventHandler) parseQueryParams(r *http.Request) (int, string, error) {
	userIDStr := r.URL.Query().Get("user_id")
	if userIDStr == "" {
		return 0, "", errors.New("user_id is required")
	}

	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		return 0, "", errors.New("invalid user_id")
	}

	date := r.URL.Query().Get("date")
	if date == "" {
		return 0, "", errors.New("date is required")
	}

	return userID, date, nil
}

func (h *EventHandler) handleServiceError(w http.ResponseWriter, err error) {
	switch err {
	case service.ErrInvalidDate, service.ErrInvalidUserID, service.ErrInvalidEventText:
		sendError(w, err.Error(), http.StatusBadRequest)
	case service.ErrEventNotFound:
		sendError(w, err.Error(), http.StatusServiceUnavailable)
	default:
		sendError(w, "internal server error", http.StatusInternalServerError)
	}
}

func sendSuccess(w http.ResponseWriter, result interface{}, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(model.Response{Result: result})
}

func sendError(w http.ResponseWriter, errorMsg string, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(model.Response{Error: errorMsg})
}
