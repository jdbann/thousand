package models

// Skill holds the details of an ability possessed by a Character.
type Skill struct {
	ID          int
	Description string `form:"description"`
	Checked     bool   `form:"checked"`
}
