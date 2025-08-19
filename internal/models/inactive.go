package models

import "time"

type Inactive struct {
	ID        int       `json:"id"`
	UserID    string    `json:"user_id"`
	UserName  string    `json:"user_name"`
	Reason    string    `json:"reason"`
	CreatedAt time.Time `json:"created_at"`
	EndAt     time.Time `json:"end_at"`
}
