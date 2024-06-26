package database

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/icrxz/crm-api-core/internal/domain"
	"github.com/jmoiron/sqlx"
)

type userDatabase struct {
	client *sqlx.DB
}

func NewUserRepository(client *sqlx.DB) domain.UserRepository {
	return &userDatabase{
		client: client,
	}
}

func (db *userDatabase) Create(ctx context.Context, user domain.User) (string, error) {
	userDTO := mapUserToUserDTO(user)

	_, err := db.client.NamedExecContext(
		ctx,
		"INSERT INTO users (user_id, first_name, last_name, email, password, role, created_at, created_by, updated_at, updated_by) VALUES (:user_id, :first_name, :last_name, :email, :password, :role, :created_at, :created_by, :updated_at, :updated_by)",
		userDTO,
	)
	if err != nil {
		return "", err
	}

	return user.UserID, nil
}

func (db *userDatabase) GetByID(ctx context.Context, userID string) (*domain.User, error) {
	var userDTO UserDTO
	err := db.client.GetContext(ctx, &userDTO, "SELECT * FROM users WHERE user_id=$1", userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.NewNotFoundError("no user found with this id", map[string]any{"user_id": userID})
		}
		return nil, err
	}

	user := mapUserDTOToUser(userDTO)

	return &user, nil
}

func (db *userDatabase) Search(ctx context.Context, filters domain.UserFilters) ([]domain.User, error) {
	whereQuery := []string{"1=1"}
	whereArgs := make([]any, 0)

	whereQuery, whereArgs = prepareInQuery(filters.Email, whereQuery, whereArgs, "email")
	whereQuery, whereArgs = prepareInQuery(filters.FirstName, whereQuery, whereArgs, "first_name")
	whereQuery, whereArgs = prepareInQuery(filters.Role, whereQuery, whereArgs, "role")
	whereQuery, whereArgs = prepareInQuery(filters.UserID, whereQuery, whereArgs, "user_id")
	whereQuery, whereArgs = prepareInQuery(filters.Region, whereQuery, whereArgs, "region")

	query := fmt.Sprintf("SELECT * FROM users WHERE %s", strings.Join(whereQuery, " AND "))

	var foundUsers []UserDTO
	err := db.client.SelectContext(ctx, &foundUsers, query, whereArgs...)
	if err != nil {
		return nil, err
	}

	users := mapUserDTOsToUsers(foundUsers)

	return users, nil
}

func (db *userDatabase) Update(ctx context.Context, userToUpdate domain.User) error {
	return nil
}

func (db *userDatabase) Delete(ctx context.Context, userID string) error {
	_, err := db.client.ExecContext(ctx, "DELETE FROM users WHERE user_id = $1", userID)
	if err != nil {
		return err
	}

	return nil
}

func prepareInQuery[S comparable](filters []S, query []string, args []any, key string) ([]string, []any) {
	if len(filters) > 0 {
		parsedArray := make([]any, 0, len(filters))
		for _, filter := range filters {
			parsedArray = append(parsedArray, filter)
		}

		queryFormatted := fmt.Sprintf("%s IN (", key)
		for i := len(args) + 1; i < len(args)+1+len(filters); i++ {
			queryFormatted += fmt.Sprintf("$%d,", i)
		}
		queryFormatted = strings.TrimRight(queryFormatted, ",")
		queryFormatted += ")"

		query = append(query, queryFormatted)
		args = append(args, parsedArray...)
	}

	return query, args
}
