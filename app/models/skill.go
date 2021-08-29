package models

type Skill struct {
	ID          int
	Description string `form:"description"`
	Checked     bool   `form:"checked"`
}
