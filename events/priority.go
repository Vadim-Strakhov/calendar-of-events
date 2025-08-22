package events

import "errors"

type Priority string

const (
	PriorityLow    Priority = "low"
	PriorityMedium Priority = "medium"
	PriorityHigh   Priority = "high"
)

func (p Priority) Validate() error {
	switch p {
	case PriorityLow, PriorityMedium, PriorityHigh:
		return nil
	default:
		return errors.New("invalid priority")
	}
}

func (p Priority) IsValid() bool  { return p.Validate() == nil }
func (p Priority) IsLow() bool    { return p == PriorityLow }
func (p Priority) IsMedium() bool { return p == PriorityMedium }
func (p Priority) IsHigh() bool   { return p == PriorityHigh }
