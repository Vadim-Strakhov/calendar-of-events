package calendar

import (
	"errors"
	"fmt"

	"github.com/Vadim-Strakhov/calendar-of-events/events"
)

var eventsMap = make(map[string]events.Event)

func checkEventExists(key string) (events.Event, error) {
	event, exists := eventsMap[key]
	if !exists {
		return events.Event{}, errors.New("событие с таким ID не найдено")
	}
	return event, nil
}

func wrapError(err error, context string) error {
	return fmt.Errorf("%s: %w", context, err)
}

func AddEvent(key string, e events.Event) error {
	if err := events.ValidateEvent(e); err != nil {
		return wrapError(err, "ошибка валидации события")
	}

	eventsMap[key] = e
	fmt.Println("Событие добавлено:", e.Title)
	return nil
}

func ShowEvents() {
	fmt.Println("\n=== Все события в календаре ===")
	if len(eventsMap) == 0 {
		fmt.Println("В календаре нет событий")
		return
	}

	for key, event := range eventsMap {
		fmt.Printf("ID: %s | Название: %s | Дата: %s\n",
			key,
			event.Title,
			event.StartAt.Format("02.01.2006 15:04"))
	}
	fmt.Println("================================")
}

func UpdateEvent(key string, newTitle string, newDateStr string) error {
	event, err := checkEventExists(key)
	if err != nil {
		return err
	}

	if err := events.UpdateEvent(&event, newTitle, newDateStr); err != nil {
		return wrapError(err, "ошибка обновления события")
	}

	eventsMap[key] = event
	fmt.Printf("Событие '%s' обновлено\n", key)
	return nil
}

func DeleteEvent(key string) error {
	_, err := checkEventExists(key)
	if err != nil {
		return err
	}

	delete(eventsMap, key)
	fmt.Printf("Событие '%s' удалено\n", key)
	return nil
}

func GetEvent(key string) (events.Event, error) {
	return checkEventExists(key)
}
