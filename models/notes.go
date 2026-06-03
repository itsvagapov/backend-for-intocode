package models

import "gorm.io/gorm"

type Note struct {
	gorm.Model
	StudentID uint    `json:"student_id"`
	Student   Student `json:"student"`
	Author    string  `json:"author"`
	Text      string  `json:"text"`
}

type CreateNoteRequest struct {
	StudentID uint   `json:"student_id"`
	Author    string `json:"author"`
	Text      string `json:"text"`
}

type UpdateNoteRequest struct {
	StudentID *uint   `json:"student_id"`
	Author    *string `json:"author"`
	Text      *string `json:"text"`
}
