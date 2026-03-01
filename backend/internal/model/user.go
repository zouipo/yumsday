package model

import (
	"time"

	"github.com/zouipo/yumsday/backend/internal/model/enum"
)

type User struct {
	ID               int64         `json:"id"`
	Username         string        `json:"username"`
	Password         string        `json:"password"`
	AppAdmin         bool          `json:"app_admin"`
	CreatedAt        time.Time     `json:"created_at"`
	Avatar           *enum.Avatar  `json:"avatar"`
	Language         enum.Language `json:"language"`
	AppTheme         enum.AppTheme `json:"theme"`
	LastVisitedGroup *int64        `json:"last_visited_group"`
}
