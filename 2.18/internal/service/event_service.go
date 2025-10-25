package service

import (
	"calendar/internal/model"
	"errors"
	"sync"
	"time"
)

var (
	// ErrEventNotFound is returned when event is not found
	ErrEventNotFound = errors.New("event not found")
	// ErrInvalidDate is returned when date format is invalid
	ErrInvalidDate = errors.New("invalid date format, expected YYYY-MM-DD")
	// ErrInvalidUserID is returned when user_id is invalid
	ErrInvalidUserID = errors.New("invalid user_id")
	// ErrInvalidEventText is returned when event text is empty
	ErrInvalidEventText = errors.New("event text cannot be empty")
)

// EventService implements business logic for working with events
type EventService struct {
	mu         sync.RWMutex
	events     map[int]*model.Event
	nextID     int
	userEvents map[int][]*model.Event
}

// NewEventService creates a new instance of event service
func NewEventService() *EventService {
	return &EventService{
		events:     make(map[int]*model.Event),
		nextID:     1,
		userEvents: make(map[int][]*model.Event),
	}
}

// CreateEvent creates a new event
func (s *EventService) CreateEvent(userID int, dateStr, eventText string) (*model.Event, error) {
	if userID <= 0 {
		return nil, ErrInvalidUserID
	}
	if eventText == "" {
		return nil, ErrInvalidEventText
	}

	date, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		return nil, ErrInvalidDate
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	event := &model.Event{
		ID:        s.nextID,
		UserID:    userID,
		Date:      date,
		EventText: eventText,
	}
	s.nextID++

	s.events[event.ID] = event
	s.userEvents[userID] = append(s.userEvents[userID], event)

	return event, nil
}

// UpdateEvent updates an existing event
func (s *EventService) UpdateEvent(id, userID int, dateStr, eventText string) (*model.Event, error) {
	if userID <= 0 {
		return nil, ErrInvalidUserID
	}
	if eventText == "" {
		return nil, ErrInvalidEventText
	}

	date, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		return nil, ErrInvalidDate
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	event, exists := s.events[id]
	if !exists {
		return nil, ErrEventNotFound
	}

	s.removeFromUserIndex(event)

	event.UserID = userID
	event.Date = date
	event.EventText = eventText

	s.userEvents[userID] = append(s.userEvents[userID], event)

	return event, nil
}

// DeleteEvent deletes an event
func (s *EventService) DeleteEvent(id int) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	event, exists := s.events[id]
	if !exists {
		return ErrEventNotFound
	}

	s.removeFromUserIndex(event)

	delete(s.events, id)

	return nil
}

// GetEventsForDay returns all events for a user on the specified day
func (s *EventService) GetEventsForDay(userID int, dateStr string) ([]*model.Event, error) {
	date, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		return nil, ErrInvalidDate
	}

	s.mu.RLock()
	defer s.mu.RUnlock()

	var result []*model.Event
	for _, event := range s.userEvents[userID] {
		if isSameDay(event.Date, date) {
			result = append(result, event)
		}
	}

	return result, nil
}

// GetEventsForWeek returns all events for a user for the week (7 days from the specified date)
func (s *EventService) GetEventsForWeek(userID int, dateStr string) ([]*model.Event, error) {
	startDate, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		return nil, ErrInvalidDate
	}

	endDate := startDate.AddDate(0, 0, 7)

	s.mu.RLock()
	defer s.mu.RUnlock()

	var result []*model.Event
	for _, event := range s.userEvents[userID] {
		if (event.Date.Equal(startDate) || event.Date.After(startDate)) && event.Date.Before(endDate) {
			result = append(result, event)
		}
	}

	return result, nil
}

// GetEventsForMonth returns all events for a user for the month
func (s *EventService) GetEventsForMonth(userID int, dateStr string) ([]*model.Event, error) {
	date, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		return nil, ErrInvalidDate
	}

	year, month, _ := date.Date()
	startDate := time.Date(year, month, 1, 0, 0, 0, 0, date.Location())
	endDate := startDate.AddDate(0, 1, 0)

	s.mu.RLock()
	defer s.mu.RUnlock()

	var result []*model.Event
	for _, event := range s.userEvents[userID] {
		if (event.Date.Equal(startDate) || event.Date.After(startDate)) && event.Date.Before(endDate) {
			result = append(result, event)
		}
	}

	return result, nil
}

func (s *EventService) removeFromUserIndex(event *model.Event) {
	events := s.userEvents[event.UserID]
	for i, e := range events {
		if e.ID == event.ID {
			s.userEvents[event.UserID] = append(events[:i], events[i+1:]...)
			break
		}
	}
}

func isSameDay(date1, date2 time.Time) bool {
	y1, m1, d1 := date1.Date()
	y2, m2, d2 := date2.Date()
	return y1 == y2 && m1 == m2 && d1 == d2
}
