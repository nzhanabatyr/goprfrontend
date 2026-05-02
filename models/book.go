package models

type Book struct {
	ID       uint   `json:"id" gorm:"primaryKey"`
	Title    string `json:"title"`
	AuthorID uint   `json:"author_id"`
}
