package models

// Mark holds the details of a mark which betrays the player as a Vampire.
type Mark struct {
	ID          int
	Description string `form:"description"`
}
