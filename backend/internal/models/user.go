package models

import (
	"time"

	"github.com/zouipo/yumsday/backend/internal/models/enums"
)

type User struct {
	ID       int64  `json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
	AppAdmin bool   `json:"app_admin"`
	// handle both TIMESTAMP or DATETIME SQL types
	CreatedAt time.Time      `json:"created_at"`
	Avatar    enums.Avatar   `json:"avatar"`
	Language  enums.Language `json:"language"`
	AppTheme  enums.AppTheme `json:"app_theme"`
	//lastVisitedGroup
}
