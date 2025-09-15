package calendar

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/Vadim-Strakhov/calendar-of-events/events"
	"github.com/Vadim-Strakhov/calendar-of-events/storage"
)

type Calendar struct {
	eventsMap    map[string]events.Event
	storage      storage.Store
	Notification chan string
}

func NewCalendar(s storage.Store) *Calendar {
	return &Calendar{
		eventsMap:    make(map[string]events.Event),
		storage:      s,
		Notification: make(chan string),
	}
}

func (c *Calendar) checkEventExists(key string) (events.Event, error) {
	event, exists := c.eventsMap[key]
	if !exists {
		return events.Event{}, errors.New("событие с таким ID не найдено")
	}
	return event, nil
}

func (c *Calendar) wrapError(err error, context string) error {
	return fmt.Errorf("%s: %w", context, err)
}

func (c *Calendar) AddEvent(key string, e events.Event) error {
	if err := events.ValidateEvent(e); err != nil {
		return c.wrapError(err, "ошибка валидации события")
	}

	c.eventsMap[key] = e
	return nil
}

func (c *Calendar) GetEvents() map[string]events.Event {
	return c.eventsMap
}

func (c *Calendar) UpdateEvent(key string, newTitle string, newDateStr string, newPriority events.Priority) error {
	event, err := c.checkEventExists(key)
	if err != nil {
		return err
	}

	if err := events.UpdateEvent(&event, newTitle, newDateStr, newPriority); err != nil {
		return c.wrapError(err, "ошибка обновления события")
	}

	c.eventsMap[key] = event
	return nil
}

func (c *Calendar) DeleteEvent(key string) error {
	event, err := c.checkEventExists(key)
	if err != nil {
		return err
	}

	if event.Reminder != nil {
		event.Reminder.Stop()
	}

	delete(c.eventsMap, key)
	return nil
}

func (c *Calendar) GetEvent(key string) (events.Event, error) {
	return c.checkEventExists(key)
}

func (c *Calendar) Save() error {
	data, err := json.MarshalIndent(c.eventsMap, "", "  ")
	if err != nil {
		return c.wrapError(err, "ошибка маршалинга в JSON")
	}

	err = c.storage.Save(data)
	if err != nil {
		return c.wrapError(err, "ошибка сохранения в файл")
	}
	return nil
}

func (c *Calendar) Load() error {
	data, err := c.storage.Load()
	if err != nil {
		return c.wrapError(err, "ошибка загрузки из файла")
	}

	err = json.Unmarshal(data, &c.eventsMap)
	if err != nil {
		return c.wrapError(err, "ошибка десериализации JSON")
	}
	return nil
}

func (c *Calendar) GetEventsMap() map[string]events.Event {
	return c.eventsMap
}

func (c *Calendar) SetEventReminder(key string, message string, at time.Time) error {
	event, err := c.checkEventExists(key)
	if err != nil {
		return err
	}
	event.AddReminder(message, at, c.Notify)
	c.eventsMap[key] = event
	return nil
}

func (c *Calendar) RemoveEventReminder(key string) error {
	event, err := c.checkEventExists(key)
	if err != nil {
		return err
	}
	event.RemoveReminder()
	c.eventsMap[key] = event
	return nil
}

func (c *Calendar) SetEventReminderBefore(key string, message string, durationBefore time.Duration) error {
	event, err := c.checkEventExists(key)
	if err != nil {
		return err
	}

	event.AddReminderWithDuration(message, durationBefore, c.Notify)
	c.eventsMap[key] = event
	return nil
}

func (c *Calendar) CancelEventReminder(key string) error {
	event, err := c.checkEventExists(key)
	if err != nil {
		return err
	}

	if event.Reminder == nil {
		return errors.New("у события нет напоминания")
	}

	event.CancelReminder()
	c.eventsMap[key] = event
	return nil
}

func (c *Calendar) Notify(msg string) {
	c.Notification <- msg
}

func (c *Calendar) Close() {
	close(c.Notification)
}
