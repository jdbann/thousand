package models

// Resource holds the details of a resource possessed by a Vampire.
type Resource struct {
	ID          int
	Description string `form:"description"`
	Stationary  bool   `form:"stationary"`
	Lost        bool   `form:"lost"`
}
