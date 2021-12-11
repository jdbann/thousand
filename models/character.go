package models

import (
	"github.com/google/uuid"
	"go.uber.org/zap/zapcore"
)

type Character struct {
	ID   uuid.UUID
	Name string
	Type string
}

type CreateCharacterParams struct {
	Name string `form:"name"`
	Type string `form:"type"`
}

func (p CreateCharacterParams) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	enc.AddString("name", p.Name)
	enc.AddString("type", p.Type)
	return nil
}
