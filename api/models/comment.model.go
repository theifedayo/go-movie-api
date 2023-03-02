package models

import "time"

type Comment struct {
    ID        uint      `json:"id" gorm:"primary_key"`
    MovieID   uint      `json:"-"`
    Comment   string    `json:"comment"`
    CreatedAt time.Time `json:"created_at"`
}