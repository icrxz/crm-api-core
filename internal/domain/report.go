package domain

import "time"

type Report struct {
	ReportID       string
	ReportName     string
	ReportTemplate string
	Active         bool
	CreatedBy      string
	CreatedAt      time.Time
	UpdatedBy      string
	UpdatedAt      time.Time
}
