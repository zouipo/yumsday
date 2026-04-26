package model

import "time"

type GroupMember struct {
	UserID   int64     `json:"user_id"`
	GroupId  int64     `json:"group_id"`
	Admin    bool      `json:"admin"`
	JoinedAt time.Time `json:"joined_at"`
}
