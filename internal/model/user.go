package model

type User struct {
	ID        string `json:"id"`
	FirstName string `json:"first_name" validate:"required"`
	LastName  string `json:"last_name" validate:"required"`
	Nickname  string `json:"nickname" validate:"required"`
	Password  string `json:"password" validate:"required"`
	Email     string `json:"email" validate:"email"`
	Country   string `json:"country"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}
