package database

import (
	"context"
	"strconv"

	"github.com/icrxz/crm-api-core/internal/domain"
	"github.com/jmoiron/sqlx"
)

type userDatabase struct {
	client *sqlx.DB
}

func NewUserDatabase(client *sqlx.DB) domain.UserRepository {
	return &userDatabase{
		client: client,
	}
}

func (db *userDatabase) Create(ctx context.Context, user domain.User) (string, error) {
	userDTO := mapUserToUserDTO(user)

	createdUser, err := db.client.NamedExecContext(
		ctx,
		"INSERT INTO user (first_name, last_name, email, password) VALUES (:first_name, :last_name, :email, :password)",
		userDTO,
	)
	if err != nil {
		return "", err
	}

	createdID, err := createdUser.LastInsertId()
	if err != nil {
		return "", err
	}

	userID := strconv.FormatInt(createdID, 10)

	return userID, nil
}

func (db *userDatabase) GetByID(ctx context.Context, userID string) (*domain.User, error) {
	db.client.NamedQueryContext(ctx, "SELECT * FROM crm_users WHERE user_id = :user_id", userID)
	return nil, nil
}

func (db *userDatabase) Search(ctx context.Context, filters map[string]string) ([]domain.User, error) {
	return nil, nil
}

func (db *userDatabase) Update(ctx context.Context, userToUpdate domain.User) error {
	return nil
}

func (db *userDatabase) Delete(ctx context.Context, userID string) error {
	return nil
}
