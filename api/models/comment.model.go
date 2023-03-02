package models

import "time"

type Comment struct {
	ID        uint      `json:"id" gorm:"primary_key"`
	MovieID   string    `json:"movie_id"`
	Comment   string    `json:"comment" gorm:"not null;size:500"`
	IP        string    `json:"ip"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at,omitempty"`
}

type CreateCommentRequest struct {
	MovieID   string    `json:"movie_id"`
	Comment   string    `json:"comment" binding:"required"`
	IP        string    `json:"ip"`
	CreatedAt time.Time `json:"created_at,omitempty"`
	UpdatedAt time.Time `json:"updated_at,omitempty"`
}

//TODO: Set comment max to 500 characters
