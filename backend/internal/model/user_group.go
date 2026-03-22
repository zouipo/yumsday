package model

import "time"

type UserGroup struct {
	ID        int64     `json:"id"`
	Name      string    `json:"name"`
	ImageURL  string    `json:"image_url"`
	CreatedAt time.Time `json:"created_at"`
}
