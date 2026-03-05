package models

import "time"

// Bug represents a submitted bug report.

type Bug struct {
    ID          int64     `db:"id" json:"id"`
    Title       string    `db:"title" json:"title"`
    Description string    `db:"description" json:"description"`
    ReporterID  int64     `db:"reporter_id" json:"reporter_id"`
    Status      string    `db:"status" json:"status"`
    Priority    string    `db:"priority" json:"priority"`
    Category    string    `db:"category" json:"category"`
    CreatedAt   time.Time `db:"created_at" json:"created_at"`
}