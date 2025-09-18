package reminder

import (
	"fmt"
	"time"
)

func testReminders() {
	fmt.Print("\n=== Тест напоминаний ===\n")

	eventTime := time.Now().Add(10 * time.Second)
	fmt.Printf("Событие запланировано на: %s\n", eventTime.Format("15:04:05"))

	reminderTime := eventTime.Add(-5 * time.Second)
	fmt.Printf("Напоминание установлено на: %s\n", reminderTime.Format("15:04:05"))

	notifier := func(msg string) { fmt.Println(msg) }

	r := NewReminder("Проверить подготовку к событию!", reminderTime, notifier)
	r.Start()

	fmt.Println("Ждем срабатывания напоминания...")
	time.Sleep(6 * time.Second)

	fmt.Print("=== Тест завершен ===\n")
}
