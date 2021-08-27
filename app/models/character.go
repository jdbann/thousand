package models

// Character holds everything related to a character.
type Character struct {
	Details *Details
}

// Details holds the important details specific to a character.
type Details struct {
	Name string `form:"name"`
}
