package models

// Character holds the details of a character known by a Vampire.
type Character struct {
	ID   int
	Name string `form:"name"`
	Type string `form:"type"`
}
