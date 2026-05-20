package repository

import (
	"database/sql"
	"log/slog"

	customErrors "github.com/zouipo/yumsday/backend/internal/error"
	"github.com/zouipo/yumsday/backend/internal/model"
)

type GroupRepositoryInterface interface {
	GetByID(id int64) (*model.Group, error)
}

type GroupRepository struct {
	db *sql.DB
}

// NewGroupRepository constructs a new GroupRepository using the provided database.
func NewGroupRepository(db *sql.DB) GroupRepositoryInterface {
	return &GroupRepository{
		db: db,
	}
}

// GetByID retrieves a group from the database by its ID, including its members.
func (r *GroupRepository) GetByID(id int64) (*model.Group, error) {
	groups, err := r.fetchGroups("WHERE groups.id = ?", id)
	if err != nil {
		return nil, err
	}

	if len(groups) == 0 {
		return nil, customErrors.NewNotFoundError("groups", "id", nil)
	}

	return &groups[0], nil
}

func (r *GroupRepository) fetchGroups(clauses string, values ...any) ([]model.Group, error) {
	query := `SELECT 
	groups.id, groups.name, groups.image_url, groups.created_at, 
	group_members.user_id, group_members.admin, group_members.joined_at
	FROM groups
	LEFT JOIN group_members ON groups.id = group_members.group_id ` + clauses

	slog.Debug("fetching groups", "query", query)

	rows, err := r.db.Query(query, values...)
	if err != nil {
		return nil, customErrors.NewInternalError("failed to fetch groups", err)
	}

	ret := []model.Group{}
	seenUsers := make(map[int64]map[int64]bool)
	stateMap := make(map[int64]int64)

	for rows.Next() {
		tmpGroup := &model.Group{}
		tmpUser := &model.GroupMember{}

		err := rows.Scan(
			&tmpGroup.ID,
			&tmpGroup.Name,
			&tmpGroup.ImageURL,
			&tmpGroup.CreatedAt,
			&tmpUser.UserID,
			&tmpUser.Admin,
			&tmpUser.JoinedAt,
		)

		if err != nil {
			return nil, customErrors.NewInternalError("failed to fetch groups", err)
		}

		id := tmpGroup.ID
		if _, exists := stateMap[id]; !exists {
			ret = append(ret, *tmpGroup)
			stateMap[id] = int64(len(ret) - 1)
			seenUsers[id] = make(map[int64]bool)
		}

		idx := stateMap[id]
		if !seenUsers[id][tmpUser.UserID] {
			ret[idx].Members = append(ret[idx].Members, *tmpUser)
			seenUsers[id][tmpUser.UserID] = true
		}
	}

	if err := rows.Err(); err != nil {
		return nil, customErrors.NewInternalError("failed to fetch groups", err)
	}

	return ret, nil
}
