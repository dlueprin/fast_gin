package model

import "time"

// Model model模板，让其他model继承，可以少写一点
type Model struct {
	ID        int `gorm:"primaryKey"`
	CreatedAt time.Time
	UpdatedAt time.Time
}
