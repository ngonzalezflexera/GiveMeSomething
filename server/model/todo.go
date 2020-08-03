package model

import "github.com/jinzhu/gorm"

type Priority int

const (
	Low = iota
	Medium
	High
)

func (p Priority) String() string {
	return [...]string{"Low", "Medium", "High"}[p]
}

type Todo struct {
	gorm.Model
	UserID      uint `json:"user_id"`
	Title       string
	URL         string
	Description string
	TimeToRead  int
	Priority    Priority
	// Check if an ID is necessary
}
