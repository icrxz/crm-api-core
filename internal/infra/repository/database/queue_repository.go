package database

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/icrxz/crm-api-core/internal/domain"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

type queueRepository struct {
	client *sqlx.DB
}

func NewQueueRepository(client *sqlx.DB) domain.QueueRepository {
	return &queueRepository{
		client: client,
	}
}

func (r *queueRepository) Create(ctx context.Context, queue domain.Queue) (string, error) {
	queueDTO := mapQueueToQueueDTO(queue)

	_, err := executor(ctx, r.client).NamedExecContext(
		ctx,
		"INSERT INTO queues "+
			"(queue_id, name, category, states, active, created_at, created_by, updated_at, updated_by) "+
			"VALUES "+
			"(:queue_id, :name, :category, :states, :active, :created_at, :created_by, :updated_at, :updated_by)",
		queueDTO,
	)
	if err != nil {
		return "", err
	}

	return queue.QueueID, nil
}

func (r *queueRepository) Update(ctx context.Context, queue domain.Queue) error {
	queueDTO := mapQueueToQueueDTO(queue)

	_, err := executor(ctx, r.client).NamedExecContext(
		ctx,
		"UPDATE queues SET "+
			"name = :name, "+
			"category = :category, "+
			"states = :states, "+
			"active = :active, "+
			"updated_at = :updated_at, "+
			"updated_by = :updated_by "+
			"WHERE queue_id = :queue_id",
		queueDTO,
	)

	return err
}

func (r *queueRepository) Delete(ctx context.Context, queueID string) error {
	if queueID == "" {
		return domain.NewValidationError("queueID is required", map[string]any{"queue_id": queueID})
	}

	_, err := executor(ctx, r.client).ExecContext(ctx, "UPDATE queues SET active = false WHERE queue_id = $1", queueID)

	return err
}

func (r *queueRepository) GetByID(ctx context.Context, queueID string) (*domain.Queue, error) {
	var queueDTO QueueDTO
	err := executor(ctx, r.client).GetContext(ctx, &queueDTO, "SELECT * FROM queues WHERE queue_id = $1", queueID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.NewNotFoundError("no queue found with this id", map[string]any{"queue_id": queueID})
		}
		return nil, err
	}

	queue := mapQueueDTOToQueue(queueDTO)

	return &queue, nil
}

func (r *queueRepository) Search(ctx context.Context, filters domain.QueueFilters) (domain.PagingResult[domain.Queue], error) {
	whereQuery := []string{"1=1"}
	whereArgs := make([]any, 0)

	whereQuery, whereArgs = prepareInQuery(filters.QueueID, whereQuery, whereArgs, "queue_id")
	whereQuery, whereArgs = prepareInQuery(filters.Category, whereQuery, whereArgs, "category")
	if len(filters.State) > 0 {
		whereQuery = append(whereQuery, fmt.Sprintf("states && $%d", len(whereArgs)+1))
		whereArgs = append(whereArgs, pq.StringArray(filters.State))
	}
	if filters.Active != nil {
		whereQuery = append(whereQuery, fmt.Sprintf("active = $%d", len(whereArgs)+1))
		whereArgs = append(whereArgs, strconv.FormatBool(*filters.Active))
	}

	limitArgs := append([]any{}, whereArgs...)
	limitArgs = append(limitArgs, filters.Limit, filters.Offset)
	limitQuery := fmt.Sprintf("LIMIT $%d OFFSET $%d", len(whereArgs)+1, len(whereArgs)+2)

	query := fmt.Sprintf("SELECT * FROM queues WHERE %s %s", strings.Join(whereQuery, " AND "), limitQuery)
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM queues WHERE %s", strings.Join(whereQuery, " AND "))

	var foundQueues []QueueDTO
	err := executor(ctx, r.client).SelectContext(ctx, &foundQueues, query, limitArgs...)
	if err != nil {
		return domain.PagingResult[domain.Queue]{}, err
	}

	var countResult int
	err = executor(ctx, r.client).GetContext(ctx, &countResult, countQuery, whereArgs...)
	if err != nil {
		return domain.PagingResult[domain.Queue]{}, err
	}

	result := domain.PagingResult[domain.Queue]{
		Result: mapQueueDTOsToQueues(foundQueues),
		Paging: domain.Paging{
			Total:  countResult,
			Limit:  filters.Limit,
			Offset: filters.Offset,
		},
	}

	return result, nil
}

func (r *queueRepository) AddMember(ctx context.Context, queueID string, userID string) error {
	_, err := executor(ctx, r.client).ExecContext(
		ctx,
		"INSERT INTO user_queues (user_id, queue_id) VALUES ($1, $2) ON CONFLICT DO NOTHING",
		userID,
		queueID,
	)

	return err
}

func (r *queueRepository) RemoveMember(ctx context.Context, queueID string, userID string) error {
	_, err := executor(ctx, r.client).ExecContext(
		ctx,
		"DELETE FROM user_queues WHERE queue_id = $1 AND user_id = $2",
		queueID,
		userID,
	)

	return err
}

func (r *queueRepository) GetMembers(ctx context.Context, queueID string) ([]domain.User, error) {
	var userDTOs []UserDTO
	err := executor(ctx, r.client).SelectContext(
		ctx,
		&userDTOs,
		"SELECT u.* FROM users u "+
			"INNER JOIN user_queues uq ON uq.user_id = u.user_id "+
			"WHERE uq.queue_id = $1",
		queueID,
	)
	if err != nil {
		return nil, err
	}

	return mapUserDTOsToUsers(userDTOs), nil
}

func (r *queueRepository) GetQueuesByUser(ctx context.Context, userID string) ([]domain.Queue, error) {
	var queueDTOs []QueueDTO
	err := executor(ctx, r.client).SelectContext(
		ctx,
		&queueDTOs,
		"SELECT q.* FROM queues q "+
			"INNER JOIN user_queues uq ON uq.queue_id = q.queue_id "+
			"WHERE uq.user_id = $1",
		userID,
	)
	if err != nil {
		return nil, err
	}

	return mapQueueDTOsToQueues(queueDTOs), nil
}
