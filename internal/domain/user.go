package domain

import (
	"context"
	"time"
)

type UserRepository interface {
	Create(ctx context.Context, user User) (string, error)
	GetByID(ctx context.Context, userID string) (*User, error)
	Search(ctx context.Context, filters UserFilters) ([]User, error)
	Update(ctx context.Context, userToUpdate User) error
	Delete(ctx context.Context, userID string) error
}

type User struct {
	UserID    string    `json:"user_id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Role      UserRole  `json:"role"`
	Region    int       `json:"region"`
	Password  string    `json:"password"`
	Orders    []Order   `json:"orders"`
	CreatedBy string    `json:"created_by"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedBy string    `json:"updated_by"`
	UpdatedAt time.Time `json:"updated_at"`
}

type UserFilters struct {
	UserID    string
	FirstName string
	LastName  string
	Email     string
}

type UserRole string

const (
	ADMIN    UserRole = "admin"
	OPERATOR UserRole = "operator"
)

func NewUser(name, email, password, author string, role UserRole, region int) User {
	now := time.Now().UTC()

	return User{
		Name:      name,
		Email:     email,
		Role:      role,
		Password:  password,
		Region:    region,
		CreatedBy: author,
		CreatedAt: now,
		UpdatedBy: author,
		UpdatedAt: now,
	}
}
