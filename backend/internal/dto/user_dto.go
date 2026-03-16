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
	Avatar    *enum.Avatar  `json:"avatar" swaggertype:"string"`
	Language  enum.Language `json:"language" swaggertype:"string"`
	AppTheme  enum.AppTheme `json:"app_theme" swaggertype:"string"`
	//lastVisitedGroup
}

type NewUserDto struct {
	Username string        `json:"username" binding:"required"`
	Password string        `json:"password"`
	AppAdmin bool          `json:"app_admin"`
	Avatar   *enum.Avatar  `json:"avatar" swaggertype:"string"`
	Language enum.Language `json:"language" swaggertype:"string"`
	AppTheme enum.AppTheme `json:"app_theme" swaggertype:"string"`
	//lastVisitedGroup
}
