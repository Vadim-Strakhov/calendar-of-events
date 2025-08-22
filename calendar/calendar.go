package calendar

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/Vadim-Strakhov/calendar-of-events/events"
	"github.com/Vadim-Strakhov/calendar-of-events/storage"
)

// Calendar представляет календарь с возможностью сохранения и загрузки
type Calendar struct {
	eventsMap map[string]events.Event
	storage   storage.Store
}

// NewCalendar создает новый календарь с указанным хранилищем
func NewCalendar(s storage.Store) *Calendar {
	return &Calendar{
		eventsMap: make(map[string]events.Event),
		storage:   s,
	}
}

// checkEventExists проверяет существование события - вынесено в отдельную функцию
func (c *Calendar) checkEventExists(key string) (events.Event, error) {
	event, exists := c.eventsMap[key]
	if !exists {
		return events.Event{}, errors.New("событие с таким ID не найдено")
	}
	return event, nil
}

// wrapError оборачивает ошибку с контекстом - утилитарная функция
func (c *Calendar) wrapError(err error, context string) error {
	return fmt.Errorf("%s: %w", context, err)
}

// AddEvent добавляет событие в календарь
func (c *Calendar) AddEvent(key string, e events.Event) error {
	// Валидируем событие перед добавлением
	if err := events.ValidateEvent(e); err != nil {
		return c.wrapError(err, "ошибка валидации события")
	}

	c.eventsMap[key] = e
	fmt.Println("Событие добавлено:", e.Title)
	return nil
}

// ShowEvents отображает все события в календаре
func (c *Calendar) ShowEvents() {
	fmt.Println("\n=== Все события в календаре ===")
	if len(c.eventsMap) == 0 {
		fmt.Println("В календаре нет событий")
		return
	}

	for key, event := range c.eventsMap {
		fmt.Printf("ID: %s | Название: %s | Дата: %s | Приоритет: %s\n",
			key,
			event.Title,
			event.StartAt.Format("02.01.2006 15:04"),
			event.Priority)
	}
	fmt.Println("================================")
}

// UpdateEvent обновляет существующее событие
func (c *Calendar) UpdateEvent(key string, newTitle string, newDateStr string, newPriority events.Priority) error {
	event, err := c.checkEventExists(key)
	if err != nil {
		return err
	}

	// Используем функцию UpdateEvent из пакета events
	if err := events.UpdateEvent(&event, newTitle, newDateStr, newPriority); err != nil {
		return c.wrapError(err, "ошибка обновления события")
	}

	c.eventsMap[key] = event
	fmt.Printf("Событие '%s' обновлено\n", key)
	return nil
}

// DeleteEvent удаляет событие по ключу
func (c *Calendar) DeleteEvent(key string) error {
	_, err := c.checkEventExists(key)
	if err != nil {
		return err
	}

	delete(c.eventsMap, key)
	fmt.Printf("Событие '%s' удалено\n", key)
	return nil
}

// GetEvent возвращает событие по ключу
func (c *Calendar) GetEvent(key string) (events.Event, error) {
	return c.checkEventExists(key)
}

// Save сохраняет календарь в JSON файл
func (c *Calendar) Save() error {
	// Маршалинг в JSON с обработкой ошибок
	data, err := json.MarshalIndent(c.eventsMap, "", "  ")
	if err != nil {
		return c.wrapError(err, "ошибка маршалинга в JSON")
	}

	// Сохраняем в хранилище
	err = c.storage.Save(data)
	if err != nil {
		return c.wrapError(err, "ошибка сохранения в файл")
	}

	fmt.Println("Календарь успешно сохранен")
	return nil
}

// Load загружает календарь из JSON файла
func (c *Calendar) Load() error {
	// Загружаем данные из хранилища
	data, err := c.storage.Load()
	if err != nil {
		return c.wrapError(err, "ошибка загрузки из файла")
	}

	// Десериализация в структуру
	err = json.Unmarshal(data, &c.eventsMap)
	if err != nil {
		return c.wrapError(err, "ошибка десериализации JSON")
	}

	fmt.Println("Календарь успешно загружен")
	return nil
}

// GetEventsMap возвращает карту событий (для тестирования)
func (c *Calendar) GetEventsMap() map[string]events.Event {
	return c.eventsMap
}

// Методы календаря для напоминаний
func (c *Calendar) SetEventReminder(key string, message string, at time.Time) error {
	event, err := c.checkEventExists(key)
	if err != nil {
		return err
	}
	event.AddReminder(message, at)
	c.eventsMap[key] = event
	fmt.Printf("Напоминание установлено для события '%s'\n", key)
	return nil
}

func (c *Calendar) RemoveEventReminder(key string) error {
	event, err := c.checkEventExists(key)
	if err != nil {
		return err
	}
	event.RemoveReminder()
	c.eventsMap[key] = event
	fmt.Printf("Напоминание удалено для события '%s'\n", key)
	return nil
}
