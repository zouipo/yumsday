package dto

import (
	"time"

	"github.com/zouipo/yumsday/backend/internal/model/enum"
)

type UserDto struct {
	ID        int64         `json:"id"`
	Username  string        `json:"username" binding:"required"`
	AppAdmin  bool          `json:"app_admin"`
	CreatedAt time.Time     `json:"created_at"`
	Avatar    *enum.Avatar  `json:"avatar"`
	Language  enum.Language `json:"language"`
	AppTheme  enum.AppTheme `json:"app_theme"`
	//lastVisitedGroup
}

type NewUserDto struct {
	Username string        `json:"username" binding:"required"`
	Password string        `json:"password"`
	AppAdmin bool          `json:"app_admin"`
	Avatar   *enum.Avatar  `json:"avatar"`
	Language enum.Language `json:"language"`
	AppTheme enum.AppTheme `json:"app_theme"`
	//lastVisitedGroup
}
