package models

import "gorm.io/gorm"

type Group struct {
	gorm.Model
	Title       string `json:"title"`
	CurrentWeek int    `json:"current_week"`
	TotalWeeks  int    `json:"total_weeks"`
	IsFinished  bool   `json:"is_finished"`
}

type CreateGroupRequest struct {
	Title       string `json:"title"`
	CurrentWeek int    `json:"current_week"`
	TotalWeeks  int    `json:"total_weeks"`
	IsFinished  bool   `json:"is_finished"`
}

type UpdateGroupRequest struct {
	Title       *string `json:"title"`
	CurrentWeek *int    `json:"current_week"`
	TotalWeeks  *int    `json:"total_weeks"`
	IsFinished  *bool   `json:"is_finished"`
}
