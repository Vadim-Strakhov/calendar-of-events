package main

import (
	"fmt"
	"log"

	"github.com/Vadim-Strakhov/calendar-of-events/calendar"
	"github.com/Vadim-Strakhov/calendar-of-events/events"
)

func main() {
	fmt.Println("=== Тестирование валидации заголовков с regexp ===")

	fmt.Println("\n--- Создание валидных событий ---")
	event1, err := events.NewEvent("Meeting with client", "2024-01-15 14:30")
	if err != nil {
		log.Fatal("Ошибка создания события 1:", err)
	}
	fmt.Println("✓ Событие 1 создано успешно")

	event2, err := events.NewEvent("Project presentation", "2024-01-16 10:00")
	if err != nil {
		log.Fatal("Ошибка создания события 2:", err)
	}
	fmt.Println("✓ Событие 2 создано успешно")

	event3, err := events.NewEvent("Team lunch", "2024-01-15 12:00")
	if err != nil {
		log.Fatal("Ошибка создания события 3:", err)
	}
	fmt.Println("✓ Событие 3 создано успешно")

	fmt.Println("\n--- Тестирование невалидных заголовков ---")

	_, err = events.NewEvent("", "2024-01-15 14:30")
	if err != nil {
		fmt.Printf("✓ Ожидаемая ошибка для пустого заголовка: %v\n", err)
	}

	_, err = events.NewEvent("Hi", "2024-01-15 14:30")
	if err != nil {
		fmt.Printf("✓ Ожидаемая ошибка для короткого заголовка: %v\n", err)
	}

	longTitle := "This is a very long title that exceeds the maximum allowed length of fifty characters"
	_, err = events.NewEvent(longTitle, "2024-01-15 14:30")
	if err != nil {
		fmt.Printf("✓ Ожидаемая ошибка для длинного заголовка: %v\n", err)
	}

	_, err = events.NewEvent("Встреча с клиентом", "2024-01-15 14:30")
	if err != nil {
		fmt.Printf("✓ Ожидаемая ошибка для заголовка с кириллицей: %v\n", err)
	}

	_, err = events.NewEvent("Meeting!@#$%", "2024-01-15 14:30")
	if err != nil {
		fmt.Printf("✓ Ожидаемая ошибка для заголовка со спецсимволами: %v\n", err)
	}

	_, err = events.NewEvent("   ", "2024-01-15 14:30")
	if err != nil {
		fmt.Printf("✓ Ожидаемая ошибка для заголовка из пробелов: %v\n", err)
	}

	fmt.Println("\n--- Добавление событий в календарь ---")
	if err := calendar.AddEvent("meeting1", event1); err != nil {
		log.Fatal("Ошибка добавления события 1:", err)
	}

	if err := calendar.AddEvent("presentation", event2); err != nil {
		log.Fatal("Ошибка добавления события 2:", err)
	}

	if err := calendar.AddEvent("lunch", event3); err != nil {
		log.Fatal("Ошибка добавления события 3:", err)
	}

	calendar.ShowEvents()

	fmt.Println("\n--- Тестирование редактирования событий ---")

	if err := calendar.UpdateEvent("meeting1", "Updated meeting with client", "2024-01-15 16:00"); err != nil {
		log.Fatal("Ошибка обновления события:", err)
	}

	if err := calendar.UpdateEvent("presentation", "Презентация!", "2024-01-16 10:00"); err != nil {
		fmt.Printf("✓ Ожидаемая ошибка при обновлении с недопустимыми символами: %v\n", err)
	}

	if err := calendar.UpdateEvent("lunch", "Hi", "2024-01-15 12:00"); err != nil {
		fmt.Printf("✓ Ожидаемая ошибка при обновлении с коротким заголовком: %v\n", err)
	}

	calendar.ShowEvents()

	fmt.Println("\n--- Тестирование удаления событий ---")

	if err := calendar.DeleteEvent("presentation"); err != nil {
		log.Fatal("Ошибка удаления события:", err)
	}

	if err := calendar.DeleteEvent("nonexistent"); err != nil {
		fmt.Printf("✓ Ожидаемая ошибка при удалении несуществующего события: %v\n", err)
	}

	calendar.ShowEvents()

	fmt.Println("\n--- Тестирование получения события ---")
	if event, err := calendar.GetEvent("meeting1"); err != nil {
		log.Fatal("Ошибка получения события:", err)
	} else {
		fmt.Printf("Получено событие: %s на %s\n",
			event.Title,
			event.StartAt.Format("02.01.2006 15:04"))
	}

	if _, err := calendar.GetEvent("nonexistent"); err != nil {
		fmt.Printf("✓ Ожидаемая ошибка при получении несуществующего события: %v\n", err)
	}

	fmt.Println("\n--- Тестирование допустимых символов ---")

	event4, err := events.NewEvent("Meeting 123", "2024-01-17 09:00")
	if err != nil {
		log.Fatal("Ошибка создания события с цифрами:", err)
	}
	fmt.Println("✓ Событие с цифрами создано успешно")

	event5, err := events.NewEvent("Team meeting, project discussion.", "2024-01-18 11:00")
	if err != nil {
		log.Fatal("Ошибка создания события с запятыми и точками:", err)
	}
	fmt.Println("✓ Событие с запятыми и точками создано успешно")

	if err := calendar.AddEvent("meeting123", event4); err != nil {
		log.Fatal("Ошибка добавления события с цифрами:", err)
	}

	if err := calendar.AddEvent("team_meeting", event5); err != nil {
		log.Fatal("Ошибка добавления события с запятыми и точками:", err)
	}

	calendar.ShowEvents()
}
