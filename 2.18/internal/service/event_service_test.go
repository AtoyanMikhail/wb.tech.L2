package service

import (
	"testing"
	"time"
)

func TestEventService_CreateEvent(t *testing.T) {
	service := NewEventService()

	tests := []struct {
		name      string
		userID    int
		date      string
		eventText string
		wantErr   error
	}{
		{
			name:      "valid event",
			userID:    1,
			date:      "2023-12-31",
			eventText: "New Year celebration",
			wantErr:   nil,
		},
		{
			name:      "invalid user_id",
			userID:    0,
			date:      "2023-12-31",
			eventText: "Event",
			wantErr:   ErrInvalidUserID,
		},
		{
			name:      "invalid date format",
			userID:    1,
			date:      "31-12-2023",
			eventText: "Event",
			wantErr:   ErrInvalidDate,
		},
		{
			name:      "empty event text",
			userID:    1,
			date:      "2023-12-31",
			eventText: "",
			wantErr:   ErrInvalidEventText,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			event, err := service.CreateEvent(tt.userID, tt.date, tt.eventText)

			if err != tt.wantErr {
				t.Errorf("CreateEvent() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr == nil {
				if event == nil {
					t.Error("CreateEvent() returned nil event")
					return
				}
				if event.UserID != tt.userID {
					t.Errorf("CreateEvent() userID = %v, want %v", event.UserID, tt.userID)
				}
				if event.EventText != tt.eventText {
					t.Errorf("CreateEvent() eventText = %v, want %v", event.EventText, tt.eventText)
				}
			}
		})
	}
}

func TestEventService_UpdateEvent(t *testing.T) {
	service := NewEventService()

	// Создаем событие для обновления
	event, _ := service.CreateEvent(1, "2023-12-31", "Original event")

	tests := []struct {
		name      string
		id        int
		userID    int
		date      string
		eventText string
		wantErr   error
	}{
		{
			name:      "valid update",
			id:        event.ID,
			userID:    1,
			date:      "2024-01-01",
			eventText: "Updated event",
			wantErr:   nil,
		},
		{
			name:      "event not found",
			id:        9999,
			userID:    1,
			date:      "2024-01-01",
			eventText: "Event",
			wantErr:   ErrEventNotFound,
		},
		{
			name:      "invalid date",
			id:        event.ID,
			userID:    1,
			date:      "invalid",
			eventText: "Event",
			wantErr:   ErrInvalidDate,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			updatedEvent, err := service.UpdateEvent(tt.id, tt.userID, tt.date, tt.eventText)

			if err != tt.wantErr {
				t.Errorf("UpdateEvent() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr == nil && updatedEvent.EventText != tt.eventText {
				t.Errorf("UpdateEvent() eventText = %v, want %v", updatedEvent.EventText, tt.eventText)
			}
		})
	}
}

func TestEventService_DeleteEvent(t *testing.T) {
	service := NewEventService()

	// Создаем событие для удаления
	event, _ := service.CreateEvent(1, "2023-12-31", "Event to delete")

	tests := []struct {
		name    string
		id      int
		wantErr error
	}{
		{
			name:    "valid delete",
			id:      event.ID,
			wantErr: nil,
		},
		{
			name:    "event not found",
			id:      9999,
			wantErr: ErrEventNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := service.DeleteEvent(tt.id)

			if err != tt.wantErr {
				t.Errorf("DeleteEvent() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestEventService_GetEventsForDay(t *testing.T) {
	service := NewEventService()

	// Создаем тестовые события
	service.CreateEvent(1, "2023-12-31", "Event 1")
	service.CreateEvent(1, "2023-12-31", "Event 2")
	service.CreateEvent(1, "2024-01-01", "Event 3")
	service.CreateEvent(2, "2023-12-31", "Event 4")

	tests := []struct {
		name      string
		userID    int
		date      string
		wantCount int
		wantErr   error
	}{
		{
			name:      "two events for day",
			userID:    1,
			date:      "2023-12-31",
			wantCount: 2,
			wantErr:   nil,
		},
		{
			name:      "one event for day",
			userID:    1,
			date:      "2024-01-01",
			wantCount: 1,
			wantErr:   nil,
		},
		{
			name:      "no events for day",
			userID:    1,
			date:      "2024-01-02",
			wantCount: 0,
			wantErr:   nil,
		},
		{
			name:      "invalid date",
			userID:    1,
			date:      "invalid",
			wantCount: 0,
			wantErr:   ErrInvalidDate,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			events, err := service.GetEventsForDay(tt.userID, tt.date)

			if err != tt.wantErr {
				t.Errorf("GetEventsForDay() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr == nil && len(events) != tt.wantCount {
				t.Errorf("GetEventsForDay() count = %v, want %v", len(events), tt.wantCount)
			}
		})
	}
}

func TestEventService_GetEventsForWeek(t *testing.T) {
	service := NewEventService()

	// Создаем тестовые события
	service.CreateEvent(1, "2023-12-31", "Event 1")
	service.CreateEvent(1, "2024-01-01", "Event 2")
	service.CreateEvent(1, "2024-01-05", "Event 3")
	service.CreateEvent(1, "2024-01-08", "Event 4") // За пределами недели

	events, err := service.GetEventsForWeek(1, "2023-12-31")
	if err != nil {
		t.Errorf("GetEventsForWeek() error = %v", err)
	}

	// Должно быть 3 события в неделю (31 декабря + 7 дней)
	if len(events) != 3 {
		t.Errorf("GetEventsForWeek() count = %v, want 3", len(events))
	}
}

func TestEventService_GetEventsForMonth(t *testing.T) {
	service := NewEventService()

	// Создаем тестовые события для января
	service.CreateEvent(1, "2024-01-01", "Event 1")
	service.CreateEvent(1, "2024-01-15", "Event 2")
	service.CreateEvent(1, "2024-01-31", "Event 3")
	service.CreateEvent(1, "2024-02-01", "Event 4") // Февраль

	events, err := service.GetEventsForMonth(1, "2024-01-15")
	if err != nil {
		t.Errorf("GetEventsForMonth() error = %v", err)
	}

	// Должно быть 3 события в январе
	if len(events) != 3 {
		t.Errorf("GetEventsForMonth() count = %v, want 3", len(events))
	}
}

func TestEventService_ConcurrentAccess(t *testing.T) {
	service := NewEventService()

	// Тест на data race
	done := make(chan bool)

	// Одновременное создание событий
	for i := 0; i < 10; i++ {
		go func(id int) {
			service.CreateEvent(id, "2024-01-01", "Concurrent event")
			done <- true
		}(i)
	}

	// Ждем завершения всех горутин
	for i := 0; i < 10; i++ {
		<-done
	}

	// Проверяем, что все события созданы
	events, _ := service.GetEventsForDay(1, "2024-01-01")
	if len(events) == 0 {
		t.Error("No events created in concurrent test")
	}
}

func TestIsSameDay(t *testing.T) {
	date1 := time.Date(2023, 12, 31, 10, 0, 0, 0, time.UTC)
	date2 := time.Date(2023, 12, 31, 20, 0, 0, 0, time.UTC)
	date3 := time.Date(2024, 1, 1, 10, 0, 0, 0, time.UTC)

	if !isSameDay(date1, date2) {
		t.Error("isSameDay() should return true for same day with different times")
	}

	if isSameDay(date1, date3) {
		t.Error("isSameDay() should return false for different days")
	}
}
