package repository

import (
	"database/sql"
	"errors"
	"strconv"
	"time"

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
func NewGroupRepository(db *sql.DB) *GroupRepository {
	return &GroupRepository{
		db: db,
	}
}

// GetByID retrieves a group from the database by its ID, including its members.
func (r *GroupRepository) GetByID(id int64) (*model.Group, error) {
	rows, err := r.db.Query(`
	SELECT g.id, g.name, g.ImageURL, g.created_at, gm.user_id, gm.admin, gm.joined_at
	FROM groups
	JOIN group_members gm ON g.id = gm.group_id
	WHERE g.id = ?`, id)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, customErrors.NewNotFoundError("Group", strconv.FormatInt(id, 10), err)
		}
		return nil, customErrors.NewInternalError("Failed to fetch group", err)
	}
	defer rows.Close()

	var group *model.Group

	for rows.Next() {
		var member model.GroupMember
		var groupID int64
		var groupName string
		var groupImageURL *string
		var groupCreatedAt time.Time

		err := rows.Scan(
			&groupID,
			&groupName,
			&groupImageURL,
			&groupCreatedAt,
			&member.UserID,
			&member.Admin,
			&member.JoinedAt,
		)
		if err != nil {
			return nil, customErrors.NewInternalError("Failed to scan group data", err)
		}

		if group == nil {
			group = &model.Group{
				ID:        groupID,
				Name:      groupName,
				ImageURL:  groupImageURL,
				CreatedAt: groupCreatedAt,
				Members:   []model.GroupMember{},
			}
		}

		group.Members = append(group.Members, member)
	}

	if err := rows.Err(); err != nil {
		return nil, customErrors.NewInternalError("Failed to iterate group rows", err)
	}
	return group, nil
}
