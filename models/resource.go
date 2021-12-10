package models

import (
	"github.com/google/uuid"
	"go.uber.org/zap/zapcore"
)

type Resource struct {
	ID          uuid.UUID
	VampireID   uuid.UUID
	Description string
	Stationary  bool
}

type CreateResourceParams struct {
	Description string `form:"description"`
	Stationary  bool   `form:"stationary"`
}

func (p CreateResourceParams) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	enc.AddString("description", p.Description)
	enc.AddBool("stationary", p.Stationary)
	return nil
}
