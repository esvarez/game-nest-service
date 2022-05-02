package dto

type User struct {
	Email string `json:"email" validate:"required"`
	User  string `json:"user" validate:"required"`
}
