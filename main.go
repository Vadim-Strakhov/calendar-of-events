package main

import (
	"fmt"
	"time"

	"github.com/Vadim-Strakhov/calendar-of-events/calendar"
	"github.com/Vadim-Strakhov/calendar-of-events/events"
)

func main() {
	e := events.Event{
		Title:   "Встреча",
		StartAt: time.Now(),
	}
	calendar.AddEvent("event1", e)
	fmt.Println("Календарь обновлён")
}
