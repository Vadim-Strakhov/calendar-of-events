package main

import (
	"fmt"
	"log"

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

	cli := cmd.NewCmd(c)
	cli.Run()
}
