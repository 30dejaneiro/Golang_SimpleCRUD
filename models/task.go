package models

type Task struct {
	ID        int    `json:"id" gorm:"primaryKey" `
	Title     string `json:"title" validate:"required"`
	Duration  int    `json:"duration" validate:"required,gt=10"`
	CreatedBy string `json:"created_by"`
}
