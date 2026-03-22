package model

import "time"

type MemberGroup struct {
	UserID      int64     `json:"user_id"`
	UserGroupID int64     `json:"user_group_id"`
	Admin       bool      `json:"admin"`
	JoinedAt    time.Time `json:"joined_at"`
}
