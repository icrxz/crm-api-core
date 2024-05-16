package domain

import (
	"context"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type UserRepository interface {
	Create(ctx context.Context, user User) (string, error)
	GetByID(ctx context.Context, userID string) (*User, error)
	Search(ctx context.Context, filters UserFilters) ([]User, error)
	Update(ctx context.Context, userToUpdate User) error
	Delete(ctx context.Context, userID string) error
}

type User struct {
	UserID    string
	FirstName string
	LastName  string
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
	UserID    []string
	FirstName []string
	Email     []string
	Role      []string
	Region    []string
}

type UserRole string

const (
	ADMIN    UserRole = "admin"
	OPERATOR UserRole = "operator"
)

func NewUser(firstName, lastName, email, password, author string, role UserRole, region int) (User, error) {
	now := time.Now().UTC()
	userID, err := uuid.NewRandom()
	if err != nil {
		return User{}, err
	}

	userPassword, err := generatePassword(password)
	if err != nil {
		return User{}, err
	}

	return User{
		UserID:    userID.String(),
		FirstName: firstName,
		LastName:  lastName,
		Email:     email,
		Role:      role,
		Password:  userPassword,
		Region:    region,
		CreatedBy: author,
		CreatedAt: now,
		UpdatedBy: author,
		UpdatedAt: now,
	}, nil
}

func generatePassword(password string) (string, error) {
	encryptedPass, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	return string(encryptedPass), nil
}

func (u *User) ComparePassword(passwordInput string) error {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(passwordInput))
	if err != nil {
		return err
	}

	return nil
}
