package repository

import (
	"database/sql"
	"fmt"
	"log/slog"
	"strings"

	customErrors "github.com/zouipo/yumsday/backend/internal/error"
	"github.com/zouipo/yumsday/backend/internal/model"
	"github.com/zouipo/yumsday/backend/internal/pkg/utils"
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
	opt := &utils.SelectFilteringOptions{
		Where: []utils.WhereClause{
			{Column: "groups.id", Values: []any{id}},
		},
	}
	groups, err := r.fetchGroups(opt)
	if err != nil {
		return nil, err
	}
	return &groups[0], nil
}

func (r *GroupRepository) fetchGroups(opt *utils.SelectFilteringOptions) ([]model.Group, error) {
	query := fmt.Sprintf(`
	SELECT groups.id, groups.name, groups.image_url, groups.created_at, group_members.user_id, group_members.admin, group_members.joined_at
	FROM groups
	LEFT JOIN group_members ON groups.id = group_members.group_id
	%s;`, utils.MakeSelectFiltering(opt))

	slog.Debug("fetching groups", "query", query)

	rows, err := r.db.Query(query, opt.WhereValues()...)
	if err != nil {
		return nil, customErrors.NewInternalError("failed to fetch groups", err)
	}

	m := make(map[int64]*model.Group)
	seenUsers := make(map[int64]map[int64]bool)
	tmpGroup := &model.Group{}
	tmpUser := &model.GroupMember{}

	for rows.Next() {
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
		if _, exists := m[id]; !exists {
			m[id] = tmpGroup
			seenUsers[id] = make(map[int64]bool)
		}

		if !seenUsers[id][tmpUser.UserID] {
			m[id].Members = append(m[id].Members, *tmpUser)
			seenUsers[id][tmpUser.UserID] = true
		}
	}

	if err := rows.Err(); err != nil {
		return nil, customErrors.NewInternalError("failed to fetch groups", err)
	}

	if len(m) == 0 {
		return []model.Group{}, customErrors.NewNotFoundError("Group", strings.Join(opt.WhereColumns(), ","), err)
	}

	ret := make([]model.Group, 0, len(m))
	for _, group := range m {
		ret = append(ret, *group)
	}

	return ret, nil
}
