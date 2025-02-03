package main

import (
	"time"
	"gorm.io/gorm"
)

type URLMapping struct {
	gorm.Model
	ShortCode     string    `gorm:"unique_index"`
	OriginalURL  string    `gorm:"type:text;not null"`
	CreatedBy    string    `gorm:"type:varchar(100)"`
	Enabled      bool      `gorm:"default:true"`
	ExpiresAt    *time.Time
	ClickCounter int64     `gorm:"default:0"`
	LastAccessed time.Time
}

type CreateURLRequest struct {
	URL       string `json:"url" binding:"required,url"`
	CustomAlias string `json:"custom_alias,omitempty"`
}

type URLResponse struct {
	ShortURL     string    `json:"short_url"`
	OriginalURL  string    `json:"original_url"`
	CreatedAt    time.Time `json:"created_at"`
	Clicks       int64     `json:"clicks"`
}