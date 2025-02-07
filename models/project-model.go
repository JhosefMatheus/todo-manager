package models

import "time"

type ProjectModel struct {
	Id              int
	ParentProjectId *int
	Name            string
	CreatedAt       *time.Time
	UpdatedAt       *time.Time
}
