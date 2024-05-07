package domain

import (
	"context"
	"time"
)

type UserRepository interface {
	Create(ctx context.Context, user User) (string, error)
	GetByID(ctx context.Context, userID string) (*User, error)
	Search(ctx context.Context, filters map[string]string) ([]User, error)
	Update(ctx context.Context, userToUpdate User) error
	Delete(ctx context.Context, userID string) error
}

type User struct {
	UserID    string
	Name      string
	Email     string
	Role      UserRole
	Region    int
	Password  string
	Orders    []Order
	CreatedBy string
	CreatedAt time.Time
	UpdatedBy string
	UpdatedAt time.Time
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
