package models

import "time"

type FavoriteBook struct {
	ID        uint      `gorm:"primaryKey"`
	UserID    uint      `json:"user_id"`
	BookID    uint      `json:"book_id"`
	CreatedAt time.Time `json:"created_at"`
}
