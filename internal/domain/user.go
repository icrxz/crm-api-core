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
	UserID       string
	Username     string
	FirstName    string
	LastName     string
	Email        string
	Role         UserRole
	Region       int
	Password     string
	LastLoggedIP string
	Active       bool
	Cases        []Case
	CreatedBy    string
	CreatedAt    time.Time
	UpdatedBy    string
	UpdatedAt    time.Time
}

type UserFilters struct {
	UserID    []string
	Username  []string
	FirstName []string
	Email     []string
	Role      []string
	Region    []string
	Active    *bool
}

type UserUpdate struct {
	FirstName    *string
	LastName     *string
	Email        *string
	Role         *UserRole
	Region       *int
	Password     *string
	LastLoggedIP *string
	Active       *bool
}

type UserRole string

const (
	THAVANNA_ADMIN UserRole = "thavanna_admin"
	ADMIN          UserRole = "admin"
	OPERATOR       UserRole = "operator"
)

func NewUser(firstName, lastName, email, password, author, username string, role UserRole, region int) (User, error) {
	now := time.Now().UTC()
	userID, err := uuid.NewRandom()
	if err != nil {
		return User{}, err
	}

	userPassword, err := encryptPassword(password)
	if err != nil {
		return User{}, err
	}

	return User{
		UserID:    userID.String(),
		Username:  username,
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
		Active:    true,
	}, nil
}

func encryptPassword(password string) (string, error) {
	encryptedPass, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	return string(encryptedPass), nil
}

func (u *User) ComparePassword(passwordInput string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(passwordInput))
	return err == nil
}

func (u *User) MergeUpdate(userUpdate UserUpdate, author string) {
	u.UpdatedAt = time.Now().UTC()

	if author != "" {
		u.UpdatedBy = author
	}

	if userUpdate.Active != nil {
		u.Active = *userUpdate.Active
	}

	if userUpdate.Email != nil {
		u.Email = *userUpdate.Email
	}

	if userUpdate.FirstName != nil {
		u.FirstName = *userUpdate.FirstName
	}

	if userUpdate.LastName != nil {
		u.LastName = *userUpdate.LastName
	}

	if userUpdate.LastLoggedIP != nil {
		u.LastLoggedIP = *userUpdate.LastLoggedIP
	}

	if userUpdate.Region != nil {
		u.Region = *userUpdate.Region
	}

	if userUpdate.Role != nil {
		u.Role = *userUpdate.Role
	}
}
