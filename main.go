package main

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/Vadim-Strakhov/calendar-of-events/calendar"
	"github.com/Vadim-Strakhov/calendar-of-events/events"
	"github.com/Vadim-Strakhov/calendar-of-events/storage"
)

type CustomStruct struct {
	Name        string `json:"name"`
	Age         int    `json:"age"`
	Email       string `json:"email,omitempty"`
	IsActive    bool   `json:"is_active"`
	Description string `json:"description,omitempty"`
}

func main() {
	fmt.Println("=== Демонстрация сериализации данных ===")

	// 1. Демонстрация работы с собственной структурой и JSON тегами
	fmt.Println("\n--- Тестирование JSON тегов ---")

	person1 := CustomStruct{
		Name:        "Иван Петров",
		Age:         30,
		Email:       "ivan@example.com",
		IsActive:    true,
		Description: "Активный пользователь",
	}

	person2 := CustomStruct{
		Name:     "Мария Сидорова",
		Age:      25,
		IsActive: false,
	}

	data1, err := json.MarshalIndent(person1, "", "  ")
	if err != nil {
		log.Fatal("Ошибка маршалинга person1:", err)
	}

	data2, err := json.MarshalIndent(person2, "", "  ")
	if err != nil {
		log.Fatal("Ошибка маршалинга person2:", err)
	}

	fmt.Println("Person1 (все поля):")
	fmt.Println(string(data1))

	fmt.Println("\nPerson2 (с omitempty):")
	fmt.Println(string(data2))

	// 2. Демонстрация работы с календарем и разными storage
	fmt.Println("\n--- Тестирование календаря с разными storage ---")

	// Выберите нужный стор: JSON или ZIP
	// s := storage.NewJsonStorage("calendar.json")
	s := storage.NewZipStorage("calendar.zip")

	c := calendar.NewCalendar(s)

	// Загружаем существующие данные при старте
	err = c.Load()
	if err != nil {
		fmt.Println("Файл календаря не найден или не прочитан, создаем новый календарь")
	} else {
		fmt.Println("Календарь загружен из хранилища")
		c.ShowEvents()
	}

	defer func() {
		err := c.Save()
		if err != nil {
			fmt.Println("Ошибка при сохранении календаря:", err)
		} else {
			fmt.Println("Календарь автоматически сохранен при завершении")
		}
	}()

	fmt.Println("\n--- Создание новых событий ---")
	if len(c.GetEventsMap()) == 0 {
		event1, err := events.NewEvent("Meeting with client", "2024-01-15 14:30", events.PriorityHigh)
		if err != nil {
			log.Fatal("Ошибка создания события 1:", err)
		}

		event2, err := events.NewEvent("Project presentation", "2024-01-16 10:00", events.PriorityMedium)
		if err != nil {
			log.Fatal("Ошибка создания события 2:", err)
		}

		event3, err := events.NewEvent("Team lunch", "2024-01-15 12:00", events.PriorityLow)
		if err != nil {
			log.Fatal("Ошибка создания события 3:", err)
		}

		if err := c.AddEvent("meeting1", event1); err != nil {
			log.Fatal("Ошибка добавления события 1:", err)
		}
		if err := c.AddEvent("presentation", event2); err != nil {
			log.Fatal("Ошибка добавления события 2:", err)
		}
		if err := c.AddEvent("lunch", event3); err != nil {
			log.Fatal("Ошибка добавления события 3:", err)
		}
	}

	c.ShowEvents()

	fmt.Println("\n--- Тестирование редактирования ---")
	if err := c.UpdateEvent("meeting1", "Updated meeting with client", "2024-01-15 16:00", events.PriorityMedium); err != nil {
		log.Fatal("Ошибка обновления события:", err)
	}

	c.ShowEvents()

	fmt.Println("\n--- Тестирование напоминаний ---")
	if err := c.SetEventReminder("meeting1", "Подготовить документы", time.Now().Add(24*time.Hour)); err != nil {
		log.Println("Ошибка установки напоминания:", err)
	}
	if err := c.RemoveEventReminder("presentation"); err != nil {
		log.Println("Ошибка удаления напоминания:", err)
	}

	fmt.Println("\n--- Тестирование удаления ---")
	if _, err := c.GetEvent("presentation"); err != nil {
		fmt.Println("Событие 'presentation' не найдено, пропускаем удаление")
	} else {
		if err := c.DeleteEvent("presentation"); err != nil {
			log.Fatal("Ошибка удаления события:", err)
		}
	}

	c.ShowEvents()

	fmt.Println("\n--- Демонстрация JSON тегов в Event ---")
	event4, err := events.NewEvent("Test event for JSON", "2024-01-20 15:00", events.PriorityMedium)
	if err != nil {
		log.Fatal("Ошибка создания тестового события:", err)
	}

	eventJSON, err := json.MarshalIndent(event4, "", "  ")
	if err != nil {
		log.Fatal("Ошибка маршалинга события:", err)
	}

	fmt.Println("JSON представление события:")
	fmt.Println(string(eventJSON))

	fmt.Println("\n=== Демонстрация завершена ===")
	fmt.Println("Календарь будет автоматически сохранен при завершении программы (defer)")
}
