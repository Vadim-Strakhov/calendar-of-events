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

// Event представляет событие в календаре
// JSON теги позволяют настроить сериализацию:
// - omitempty: поле не включается в JSON если оно пустое
// - json:"custom_name": изменяет имя поля в JSON
type Event struct {
	ID       string             `json:"id,omitempty"`
	Title    string             `json:"title"`
	StartAt  time.Time          `json:"start_at"`
	Priority Priority           `json:"priority"`
	Reminder *reminder.Reminder `json:"reminder,omitempty"`
}

// Компилируем регулярное выражение один раз при инициализации пакета
var titleRegex = regexp.MustCompile(`^[a-zA-Z0-9 ,\.]{3,50}$`)

// validateTitle проверяет валидность заголовка события с помощью регулярного выражения
func validateTitle(title string) error {
	title = strings.TrimSpace(title)

	if title == "" {
		return errors.New("заголовок не может быть пустым")
	}

	// Проверяем соответствие регулярному выражению
	if !titleRegex.MatchString(title) {
		return errors.New("заголовок должен содержать только буквы, цифры, пробелы, запятые и точки, длиной от 3 до 50 символов")
	}

	return nil
}

// parseDate парсит строку даты - вынесено в отдельную функцию для переиспользования
func parseDate(dateStr string) (time.Time, error) {
	return dateparse.ParseAny(dateStr)
}

// createEventFromData создает событие из валидированных данных - общая логика
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
	// Валидируем заголовок
	if err := validateTitle(title); err != nil {
		return Event{}, err
	}

	// Валидируем приоритет
	if err := priority.Validate(); err != nil {
		return Event{}, err
	}

	// Парсим дату
	parsedTime, err := parseDate(dateStr)
	if err != nil {
		return Event{}, err
	}

	return createEventFromData(title, parsedTime, priority), nil
}

// UpdateEvent обновляет существующее событие
func UpdateEvent(event *Event, newTitle string, newDateStr string, newPriority Priority) error {
	// Валидируем новый заголовок
	if err := validateTitle(newTitle); err != nil {
		return err
	}

	// Валидируем приоритет
	if err := newPriority.Validate(); err != nil {
		return err
	}

	// Парсим новую дату
	parsedTime, err := parseDate(newDateStr)
	if err != nil {
		return err
	}

	// Обновляем поля события, сохраняя ID и напоминание
	event.Title = strings.TrimSpace(newTitle)
	event.StartAt = parsedTime
	event.Priority = newPriority
	return nil
}

// ValidateEvent проверяет валидность всего события
func ValidateEvent(event Event) error {
	if err := validateTitle(event.Title); err != nil {
		return err
	}

	// Можно добавить дополнительные проверки для даты
	if event.StartAt.IsZero() {
		return errors.New("дата события не может быть пустой")
	}

	// Проверяем приоритет
	if err := event.Priority.Validate(); err != nil {
		return err
	}

	return nil
}

// Методы для управления напоминанием
func (e *Event) AddReminder(message string, at time.Time) {
	e.Reminder = reminder.NewReminder(message, at)
}

func (e *Event) RemoveReminder() {
	e.Reminder = nil
}
