package events

import (
	"errors"
	"regexp"
	"strings"
	"time"

	"github.com/Vadim-Strakhov/calendar-of-events/reminder"

	"github.com/araddon/dateparse"
	"github.com/google/uuid"
)

type Event struct {
	ID       string             `json:"id,omitempty"`
	Title    string             `json:"title"`
	StartAt  time.Time          `json:"start_at"`
	Priority Priority           `json:"priority"`
	Reminder *reminder.Reminder `json:"reminder,omitempty"`
}

var titleRegex = regexp.MustCompile(`^[\p{L}\p{N} ,\.]{3,50}$`)

func validateTitle(title string) error {
	title = strings.TrimSpace(title)

	if title == "" {
		return errors.New("заголовок не может быть пустым")
	}

	if !titleRegex.MatchString(title) {
		return errors.New("заголовок должен содержать только буквы, цифры, пробелы, запятые и точки, длиной от 3 до 50 символов")
	}

	return nil
}

func parseDate(dateStr string) (time.Time, error) {
	return dateparse.ParseAny(dateStr)
}

func createEventFromData(title string, startAt time.Time, p Priority) Event {
	return Event{
		ID:       uuid.New().String(),
		Title:    strings.TrimSpace(title),
		StartAt:  startAt,
		Priority: p,
		Reminder: nil,
	}
}

func NewEvent(title string, dateStr string, priority Priority) (Event, error) {

	if err := validateTitle(title); err != nil {
		return Event{}, err
	}

	if err := priority.Validate(); err != nil {
		return Event{}, err
	}

	parsedTime, err := parseDate(dateStr)
	if err != nil {
		return Event{}, err
	}

	return createEventFromData(title, parsedTime, priority), nil
}

func UpdateEvent(event *Event, newTitle string, newDateStr string, newPriority Priority) error {

	if err := validateTitle(newTitle); err != nil {
		return err
	}

	if err := newPriority.Validate(); err != nil {
		return err
	}

	parsedTime, err := parseDate(newDateStr)
	if err != nil {
		return err
	}

	event.Title = strings.TrimSpace(newTitle)
	event.StartAt = parsedTime
	event.Priority = newPriority
	return nil
}

func ValidateEvent(event Event) error {
	if err := validateTitle(event.Title); err != nil {
		return err
	}

	if event.StartAt.IsZero() {
		return errors.New("дата события не может быть пустой")
	}

	if err := event.Priority.Validate(); err != nil {
		return err
	}

	return nil
}

func (e *Event) AddReminder(message string, at time.Time, notify reminder.Notifier) {
	e.Reminder = reminder.NewReminder(message, at, notify)

	e.Reminder.Start()
}

func (e *Event) AddReminderWithDuration(message string, durationBefore time.Duration, notify reminder.Notifier) {
	reminderTime := e.StartAt.Add(-durationBefore)
	e.AddReminder(message, reminderTime, notify)
}

func (e *Event) RemoveReminder() {
	if e.Reminder != nil {
		e.Reminder.Stop()
		e.Reminder = nil
	}
}

func (e *Event) CancelReminder() {
	if e.Reminder != nil && !e.Reminder.Sent {
		e.Reminder.Cancel()
	}
}
