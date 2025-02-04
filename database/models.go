package database

import (
	"time"

	"gorm.io/gorm"
)

type URL struct {
	gorm.Model
	OriginalURL string    `gorm:"not null"`
	ShortCode   string    `gorm:"uniqueIndex;not null"`
	Visits      int       `gorm:"default:0"`
	LastVisit   time.Time `gorm:"default:null"`
}
