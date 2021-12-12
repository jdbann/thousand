package health

import (
	"context"
	"time"
)

const (
	_okLabel     = "ok"
	_failedLabel = "failed"
)

type Health struct {
	components []*Component
}

func New() *Health {
	return &Health{}
}

type Component struct {
	Name  string
	Check func(context.Context) error
}

func (h *Health) Register(c *Component) {
	h.components = append(h.components, c)
}

type Result struct {
	Details []ComponentResult `json:"details"`
	Status  string            `json:"status"`
}

type ComponentResult struct {
	Name      string    `json:"name"`
	Status    string    `json:"status"`
	Error     string    `json:"error,omitempty"`
	Timestamp time.Time `json:"timestamp"`
}

func (h *Health) Check(ctx context.Context) (Result, bool) {
	healthy := true

	result := Result{
		Details: make([]ComponentResult, len(h.components)),
		Status:  _okLabel,
	}

	for i := range h.components {
		cr, ok := h.components[i].PerformCheck(ctx)
		result.Details[i] = cr

		if !ok {
			healthy = false
			result.Status = _failedLabel
		}
	}

	return result, healthy
}

func (c *Component) PerformCheck(ctx context.Context) (ComponentResult, bool) {
	healthy := true

	result := ComponentResult{
		Name:      c.Name,
		Status:    _okLabel,
		Timestamp: time.Now(),
	}

	if err := c.Check(ctx); err != nil {
		healthy = false

		result.Status = _failedLabel
		result.Error = err.Error()
	}

	return result, healthy
}
