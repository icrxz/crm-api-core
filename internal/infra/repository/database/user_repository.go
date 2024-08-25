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
		"INSERT INTO users "+
			"(user_id, username, first_name, last_name, email, password, role, created_at, created_by, updated_at, updated_by, active, region, last_logged_ip) "+
			"VALUES "+
			"(:user_id, :username, :first_name, :last_name, :email, :password, :role, :created_at, :created_by, :updated_at, :updated_by, :active, :region, :last_logged_ip)",
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

func (db *userDatabase) Search(ctx context.Context, filters domain.UserFilters) (domain.PagingResult[domain.User], error) {
	whereQuery := []string{"1=1"}
	whereArgs := make([]any, 0)
	limitArgs := make([]any, 0, 2)

	whereQuery, whereArgs = prepareInQuery(filters.Email, whereQuery, whereArgs, "email")
	whereQuery, whereArgs = prepareInQuery(filters.FirstName, whereQuery, whereArgs, "first_name")
	whereQuery, whereArgs = prepareInQuery(filters.Role, whereQuery, whereArgs, "role")
	whereQuery, whereArgs = prepareInQuery(filters.UserID, whereQuery, whereArgs, "user_id")
	whereQuery, whereArgs = prepareInQuery(filters.Region, whereQuery, whereArgs, "region")
	whereQuery, whereArgs = prepareInQuery(filters.Username, whereQuery, whereArgs, "username")
	if filters.Active != nil {
		whereQuery = append(whereQuery, fmt.Sprintf("active = $%d", len(whereArgs)+1))
		whereArgs = append(whereArgs, filters.Active)
	}

	limitQuery := fmt.Sprintf("LIMIT $%d OFFSET $%d", len(whereArgs)+1, len(whereArgs)+2)
	limitArgs = append(whereArgs, filters.Limit, filters.Offset)

	query := fmt.Sprintf("SELECT * FROM users WHERE %s %s", strings.Join(whereQuery, " AND "), limitQuery)
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM users WHERE %s", strings.Join(whereQuery, " AND "))

	var foundUsers []UserDTO
	err := db.client.SelectContext(ctx, &foundUsers, query, limitArgs...)
	if err != nil {
		return domain.PagingResult[domain.User]{}, err
	}

	var countResult int
	err = db.client.GetContext(ctx, &countResult, countQuery, whereArgs...)
	if err != nil {
		return domain.PagingResult[domain.User]{}, err
	}

	users := mapUserDTOsToUsers(foundUsers)

	result := domain.PagingResult[domain.User]{
		Result: users,
		Paging: domain.Paging{
			Total:  countResult,
			Limit:  filters.Limit,
			Offset: filters.Offset,
		},
	}

	return result, nil
}

func (db *userDatabase) Update(ctx context.Context, userToUpdate domain.User) error {
	userDTO := mapUserToUserDTO(userToUpdate)

	_, err := db.client.NamedExecContext(
		ctx,
		"UPDATE users "+
			"SET first_name = :first_name, "+
			"last_name = :last_name, "+
			"email = :email, "+
			"role = :role, "+
			"updated_at = :updated_at, "+
			"updated_by = updated_by, "+
			"active = :active, "+
			"region = :region, "+
			"last_logged_ip = :last_logged_ip "+
			"WHERE user_id = :user_id",
		userDTO,
	)
	if err != nil {
		return err
	}

	return nil
}

func (db *userDatabase) Delete(ctx context.Context, userID string) error {
	if userID == "" {
		return domain.NewValidationError("userID cannot be empty", nil)
	}

	_, err := db.client.ExecContext(ctx, "UPDATE users SET active = false WHERE user_id = $1", userID)
	if err != nil {
		return err
	}

	return nil
}
