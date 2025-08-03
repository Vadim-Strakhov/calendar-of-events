package calendar

import (
	"fmt"
	"github.com/Vadim-Strakhov/calendar-of-events/events"
)

var eventsMap = make(map[string]events.Event)

func AddEvent(key string, e events.Event) {
	eventsMap[key] = e
	fmt.Println("Событие добавлено:", e.Title)
}
