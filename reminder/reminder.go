package reminder

import (
	"time"
)

type Notifier func(string)

type Reminder struct {
	Message string    `json:"message"`
	At      time.Time `json:"at"`
	Sent    bool      `json:"sent"`
	timer   *time.Timer
	notify  Notifier
}

func NewReminder(message string, at time.Time, notify Notifier) *Reminder {
	return &Reminder{
		Message: message,
		At:      at,
		Sent:    false,
		notify:  notify,
	}
}

func (r *Reminder) Start() {
	if r.Sent {
		return
	}

	delay := time.Until(r.At)
	if delay <= 0 {
		r.Send()
		return
	}

	r.timer = time.AfterFunc(delay, r.Send)
}

func (r *Reminder) Send() {
	if r.Sent {
		return
	}
	if r.notify != nil {
		r.notify("🔔 НАПОМИНАНИЕ: " + r.Message)
	}
	r.Sent = true
	if r.timer != nil {
		r.timer.Stop()
	}
}

func (r *Reminder) Stop() {
	if r.timer != nil && !r.Sent {
		_ = r.timer.Stop()
		r.timer = nil
	}
}

func (r *Reminder) Cancel() {
	if !r.Sent {
		r.Stop()
		r.Sent = true
		if r.notify != nil {
			r.notify("Отменено напоминание: " + r.Message)
		}
	}
}
