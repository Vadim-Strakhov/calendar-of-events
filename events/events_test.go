package events

import (
	"testing"
	"time"
)

func TestNewEvent_Valid(t *testing.T) {
	_, err := NewEvent("Meeting 123", "2025-09-18 10:00", PriorityHigh)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func TestNewEvent_InvalidTitle(t *testing.T) {
	_, err := NewEvent("!", time.Now().Format(time.RFC3339), PriorityLow)
	if err == nil {
		t.Fatalf("expected validation error for title, got nil")
	}
}

func TestPriorityValidate(t *testing.T) {
	if err := Priority("weird").Validate(); err == nil {
		t.Fatalf("expected error for invalid priority")
	}
	if err := PriorityHigh.Validate(); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}
