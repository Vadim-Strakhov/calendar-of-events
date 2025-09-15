package main

import (
	"fmt"
	"log"
	"time"

	"github.com/Vadim-Strakhov/calendar-of-events/calendar"
	"github.com/Vadim-Strakhov/calendar-of-events/cmd"
	"github.com/Vadim-Strakhov/calendar-of-events/storage"
)

func main() {
	s := storage.NewJsonStorage("calendar.json")

	c := calendar.NewCalendar(s)

	if err := c.Load(); err != nil {
		fmt.Println("Файл календаря не найден или не прочитан, создаем новый календарь")
	} else {
		fmt.Println("Календарь загружен из хранилища")
	}

	defer func() {
		if err := c.Save(); err != nil {
			log.Println("Ошибка при сохранении календаря:", err)
		} else {
			fmt.Println("Календарь автоматически сохранен при завершении")
		}
		c.Close()
	}()

	demonstrateTimers()

	cli := cmd.NewCmd(c)
	cli.Run()
}

func demonstrateTimers() {
	fmt.Println("\n=== Демонстрация таймеров ===")

	fmt.Println("Ждем 2 секунды...")
	time.Sleep(2 * time.Second)
	fmt.Println("Время вышло!")

	fmt.Println("\nЗапускаем таймер на 3 секунды...")
	timer1 := time.AfterFunc(3*time.Second, func() {
		fmt.Println("Таймер 1 сработал!")
	})

	time.Sleep(1 * time.Second)
	if stopped := timer1.Stop(); stopped {
		fmt.Println("Таймер 1 остановлен до срабатывания")
	} else {
		fmt.Println("Таймер 1 уже сработал или остановлен")
	}

	fmt.Println("\nЗапускаем таймер на 2 секунды...")
	timer2 := time.AfterFunc(2*time.Second, func() {
		fmt.Println("Таймер 2 сработал!")
	})

	time.Sleep(3 * time.Second)

	if stopped := timer2.Stop(); stopped {
		fmt.Println("Таймер 2 остановлен")
	} else {
		fmt.Println("Таймер 2 уже сработал")
	}

	fmt.Println("=== Демонстрация завершена ===\n")
}
