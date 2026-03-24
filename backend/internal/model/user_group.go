package model

import "time"

type UserGroup struct {
	ID        int64     `json:"id"`
	Name      string    `json:"name"`
	ImageURL  *string   `json:"image_url"`
	CreatedAt time.Time `json:"created_at"`
	Users     []struct {
		UserID   int64     `json:"user_id"`
		Admin    bool      `json:"admin"`
		JoinedAt time.Time `json:"joined_at"`
	}
}
