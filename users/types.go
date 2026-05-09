package user

import (
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type User struct {
	UserRequest
	id         string
	created_at time.Time
	admin      bool
}

type UserResponse struct {
	Id    string `json:"id"`
	Name  string `json:"name"`
	Age   int    `json:"age"`
	Email string `json:"email"`
	Admin bool   `json:"admin"`
}

type UserRequest struct {
	Name     string `json:"name"`
	Age      int    `json:"age"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UserHandler struct {
	Pool *pgxpool.Pool
}
