package cmd

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/Vadim-Strakhov/calendar-of-events/calendar"
	"github.com/Vadim-Strakhov/calendar-of-events/events"

	"github.com/c-bata/go-prompt"
	"github.com/google/shlex"
)

type Cmd struct {
	calendar *calendar.Calendar
	log      *Log
}

func NewCmd(c *calendar.Calendar) *Cmd {
	return &Cmd{
		calendar: c,
		log:      NewLog(),
	}
}

func (c *Cmd) Run() {

	go func() {
		for msg := range c.calendar.Notification {
			c.outln(msg)
		}
	}()

	p := prompt.New(
		c.executor,
		c.completer,
		prompt.OptionPrefix("> "),
	)
	p.Run()
}

func (c *Cmd) completer(d prompt.Document) []prompt.Suggest {
	suggestions := []prompt.Suggest{
		{Text: "add", Description: "Добавить событие: add \"key\" \"title\" \"date\" \"priority\""},
		{Text: "list", Description: "Показать все события"},
		{Text: "remove", Description: "Удалить событие: remove \"key\""},
		{Text: "update", Description: "Обновить событие: update \"key\" \"title\" \"date\" \"priority\""},
		{Text: "reminder", Description: "Установить напоминание: reminder \"key\" \"message\" \"time\""},
		{Text: "reminder-before", Description: "Установить напоминание за время до события: reminder-before \"key\" \"message\" \"duration\""},
		{Text: "cancel-reminder", Description: "Отменить напоминание: cancel-reminder \"key\""},
		{Text: "show-reminders", Description: "Показать активные напоминания"},
		{Text: "log", Description: "Показать историю лога"},
		{Text: "help", Description: "Показать справку"},
		{Text: "exit", Description: "Сохранить и выйти"},
	}
	return prompt.FilterHasPrefix(suggestions, d.GetWordBeforeCursor(), true)
}

func (c *Cmd) executor(input string) {
	input = strings.TrimSpace(input)
	if input == "" {
		return
	}
	c.log.Add("INPUT: ", input)

	parts, err := shlex.Split(input)
	if err != nil || len(parts) == 0 {
		c.outln("Ошибка парсинга ввода. Используйте кавычки для аргументов.")
		return
	}

	cmd := strings.ToLower(parts[0])

	switch cmd {
	case "help":
		c.outln("Доступные команды:")
		c.outln("  add \"key\" \"title\" \"date\" \"priority\"  — добавить событие")
		c.outln("  list                                     — показать все события")
		c.outln("  remove \"key\"                             — удалить событие")
		c.outln("  update \"key\" \"title\" \"date\" \"priority\" — обновить событие")
		c.outln("  reminder \"key\" \"msg\" \"time\"            — установить напоминание")
		c.outln("  reminder-before \"key\" \"msg\" \"duration\" — напоминание за время до события")
		c.outln("  cancel-reminder \"key\"                    — отменить напоминание")
		c.outln("  show-reminders                           — показать активные напоминания")
		c.outln("  log                                      — показать историю лога")
		c.outln("  exit                                     — сохранить и выйти")

	case "list":
		eventsMap := c.calendar.GetEvents()
		c.outln("\n=== Все события в календаре ===")
		if len(eventsMap) == 0 {
			c.outln("В календаре нет событий")
			break
		}
		for key, event := range eventsMap {
			c.outf("ID: %s | Название: %s | Дата: %s | Приоритет: %s\n",
				key,
				event.Title,
				event.StartAt.Format("02.01.2006 15:04"),
				event.Priority)
		}
		c.outln("================================")

	case "add":

		if len(parts) < 5 {
			c.outln("Формат: add \"key\" \"название\" \"дата и время\" \"приоритет\"")
			return
		}
		key := parts[1]
		title := parts[2]
		date := parts[3]
		priority := events.Priority(strings.ToLower(parts[4]))

		e, err := events.NewEvent(title, date, priority)
		if err != nil {
			c.outln("Ошибка создания события:", err)
			return
		}
		if err := c.calendar.AddEvent(key, e); err != nil {
			c.outln("Ошибка добавления:", err)
			return
		}
		c.outln("Событие добавлено:", e.Title)

	case "remove":

		if len(parts) < 2 {
			c.outln("Формат: remove \"key\"")
			return
		}
		key := parts[1]
		if err := c.calendar.DeleteEvent(key); err != nil {
			c.outln("Ошибка удаления:", err)
			return
		}

	case "update":

		if len(parts) < 5 {
			c.outln("Формат: update \"key\" \"название\" \"дата и время\" \"приоритет\"")
			return
		}
		key := parts[1]
		title := parts[2]
		date := parts[3]
		priority := events.Priority(strings.ToLower(parts[4]))
		if err := c.calendar.UpdateEvent(key, title, date, priority); err != nil {
			c.outln("Ошибка обновления:", err)
			return
		}

	case "reminder":

		if len(parts) < 4 {
			c.outln("Формат: reminder \"key\" \"сообщение\" \"дата и время\"")
			return
		}
		key := parts[1]
		message := parts[2]
		dateStr := parts[3]

		reminderTime, err := time.Parse("2006-01-02 15:04", dateStr)
		if err != nil {
			c.outln("Ошибка парсинга даты. Используйте формат: 2025-06-16 15:50")
			return
		}

		if err := c.calendar.SetEventReminder(key, message, reminderTime); err != nil {
			c.outln("Ошибка установки напоминания:", err)
		}

	case "reminder-before":

		if len(parts) < 4 {
			c.outln("Формат: reminder-before \"key\" \"сообщение\" \"5m\" (5m, 1h, 30s)")
			return
		}
		key := parts[1]
		message := parts[2]
		durationStr := parts[3]

		duration, err := time.ParseDuration(durationStr)
		if err != nil {
			c.outln("Ошибка парсинга продолжительности. Примеры: 5m, 1h, 30s")
			return
		}

		if err := c.calendar.SetEventReminderBefore(key, message, duration); err != nil {
			c.outln("Ошибка установки напоминания:", err)
		}

	case "cancel-reminder":

		if len(parts) < 2 {
			c.outln("Формат: cancel-reminder \"key\"")
			return
		}
		key := parts[1]
		if err := c.calendar.CancelEventReminder(key); err != nil {
			c.outln("Ошибка отмены напоминания:", err)
		}

	case "show-reminders":
		c.showReminders()

	case "log":
		entries := c.log.All()
		if len(entries) == 0 {
			c.outln("(лог пуст)")
			return
		}
		c.outln("\n=== ЛОГ КОНСОЛИ ===")
		for _, e := range entries {
			c.outln(e)
		}
		c.outln("===================")

	case "exit":
		if err := c.calendar.Save(); err != nil {
			c.outln("Ошибка при сохранении:", err)

		} else {
			c.outln("Календарь сохранен. Выход.")
		}
		c.calendar.Close()
		os.Exit(0)

	default:
		c.outln("Неизвестная команда.")
		c.outln("Введите 'help' для списка команд.")
	}
}

func (c *Cmd) showReminders() {
	c.outln("\n=== Активные напоминания ===")
	hasReminders := false

	for _, event := range c.calendar.GetEventsMap() {
		if event.Reminder != nil && !event.Reminder.Sent {
			c.outf("Событие: %s | Напоминание: %s | Время: %s\n",
				event.Title,
				event.Reminder.Message,
				event.Reminder.At.Format("02.01.2006 15:04"))
			hasReminders = true
		}
	}

	if !hasReminders {
		c.outln("Активных напоминаний нет")
	}
	c.outln("=============================")
}

func (c *Cmd) outln(a ...any) {
	s := fmt.Sprintln(a...)
	fmt.Print(s)
	c.log.Add("OUT: ", strings.TrimRight(s, "\n"))
}

func (c *Cmd) outf(format string, a ...any) {
	s := fmt.Sprintf(format, a...)
	fmt.Print(s)
	if len(s) > 0 && s[len(s)-1] == '\n' {
		c.log.Add("OUT: ", strings.TrimRight(s, "\n"))
	} else {
		c.log.Add("OUT: ", s)
	}
}
