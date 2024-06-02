package models

type Todo struct {
	ID        uint   `gorm:"primaryKey"`
	Title     string `gorm:"not null"`
	Completed bool   `gorm:"not null"`
}
