package dtos

import (
	"time"

	"github.com/zouipo/yumsday/backend/internal/models/enums"
)

type UserDto struct {
	ID        int64          `json:"id"`
	Username  string         `json:"username" binding:"required"`
	AppAdmin  bool           `json:"app_admin"`
	CreatedAt time.Time      `json:"created_at"`
	Avatar    *enums.Avatar  `json:"avatar"`
	Language  enums.Language `json:"language"`
	AppTheme  enums.AppTheme `json:"app_theme"`
	//lastVisitedGroup
}

type NewUserDto struct {
	Username string         `json:"username" binding:"required"`
	Password string         `json:"password"`
	AppAdmin bool           `json:"app_admin"`
	Avatar   *enums.Avatar  `json:"avatar"`
	Language enums.Language `json:"language"`
	AppTheme enums.AppTheme `json:"app_theme"`
	//lastVisitedGroup
}
