package models

import (
	"time"
)

// Bug represents a submitted bug report.
type Bug struct {
	ID          int64          `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	Title       string         `gorm:"column:title;not null" json:"title"`
	Description string         `gorm:"column:description;not null" json:"description"`
	ReporterID  int64          `gorm:"column:reporter_id;not null" json:"reporter_id"`
	Status      string         `gorm:"column:status;not null" json:"status"`
	Priority    string         `gorm:"column:priority;not null" json:"priority"`
	Category    string         `gorm:"column:category;not null" json:"category"`
	CreatedAt   time.Time      `gorm:"column:created_at;autoCreateTime" json:"created_at"`
}