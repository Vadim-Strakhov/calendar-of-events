package calendar

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/Vadim-Strakhov/calendar-of-events/events"
	"github.com/Vadim-Strakhov/calendar-of-events/storage"
)

func TestCalendar_SaveLoad(t *testing.T) {
	dir := t.TempDir()
	file := filepath.Join(dir, "calendar_test.json")

	store := storage.NewJsonStorage(file)
	cal := NewCalendar(store)

	ev, err := events.NewEvent("Meeting 123", "2025-09-18 10:00", events.PriorityMedium)
	if err != nil {
		t.Fatalf("unexpected error creating event: %v", err)
	}
	if err := cal.AddEvent(ev.ID, ev); err != nil {
		t.Fatalf("unexpected error adding event: %v", err)
	}

	if err := cal.Save(); err != nil {
		t.Fatalf("save failed: %v", err)
	}
	if _, err := os.Stat(file); err != nil {
		t.Fatalf("expected file to exist: %v", err)
	}

	cal2 := NewCalendar(storage.NewJsonStorage(file))
	if err := cal2.Load(); err != nil {
		t.Fatalf("load failed: %v", err)
	}

	if len(cal2.GetEventsMap()) != 1 {
		t.Fatalf("expected 1 event after load, got %d", len(cal2.GetEventsMap()))
	}
}

func TestCalendar_NotificationChannel(t *testing.T) {
	dir := t.TempDir()
	file := filepath.Join(dir, "calendar_test.json")
	cal := NewCalendar(storage.NewJsonStorage(file))

	ev := events.Event{
		ID:       "id-1",
		Title:    "Notify",
		StartAt:  time.Now().Add(time.Second),
		Priority: events.PriorityLow,
	}
	if err := cal.AddEvent(ev.ID, ev); err != nil {
		t.Fatalf("add event: %v", err)
	}

	go cal.Notify("hello")

	select {
	case <-cal.Notification:
		// ok
	case <-time.After(2 * time.Second):
		t.Fatal("expected to receive notification")
	}
}
