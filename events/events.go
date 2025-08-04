package events

import (
	"errors"
	"regexp"
	"strings"
	"time"

	"github.com/araddon/dateparse"
	"github.com/google/uuid"
)

type Event struct {
	ID      string
	Title   string
	StartAt time.Time
}

var titleRegex = regexp.MustCompile(`^[a-zA-Z0-9 ,\.]{3,50}$`)

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

func createEventFromData(title string, startAt time.Time) Event {
	return Event{
		ID:      uuid.New().String(),
		Title:   strings.TrimSpace(title),
		StartAt: startAt,
	}
}

func NewEvent(title string, dateStr string) (Event, error) {
	if err := validateTitle(title); err != nil {
		return Event{}, err
	}

	parsedTime, err := parseDate(dateStr)
	if err != nil {
		return Event{}, err
	}

	return createEventFromData(title, parsedTime), nil
}

func UpdateEvent(event *Event, newTitle string, newDateStr string) error {
	if err := validateTitle(newTitle); err != nil {
		return err
	}

	parsedTime, err := parseDate(newDateStr)
	if err != nil {
		return err
	}

	*event = createEventFromData(newTitle, parsedTime)
	return nil
}

func ValidateEvent(event Event) error {
	if err := validateTitle(event.Title); err != nil {
		return err
	}

	if event.StartAt.IsZero() {
		return errors.New("дата события не может быть пустой")
	}

	return nil
}
