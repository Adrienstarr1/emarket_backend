package user

import (
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type User struct {
	Name       string
	Age        int
	Email      string
	Password   string
	id         string
	created_at time.Time
	admin      bool
}

type UserResponse struct {
	Name     string `json:"name"`
	Age      int    `json:"age"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UserHandler struct {
	Pool *pgxpool.Pool
}
