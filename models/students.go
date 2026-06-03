package models

import "gorm.io/gorm"

type Student struct {
	gorm.Model
	FullName      string        `json:"full_name"`
	Email         string        `json:"email"`
	Telegram      string        `json:"telegram"`
	GroupID       uint          `json:"group_id"`
	Group         Group         `json:"group,omitempty"`
	TuitionTotal  int           `json:"tuition_total"`
	TuitionPaid   int           `json:"tuition_paid"`
	PaymentStatus PaymentStatus `json:"payment_status"`
	StudyStatus   StudyStatus   `json:"study_status"`
}

type PaymentStatus string

const (
	PaymentPaid    PaymentStatus = "paid"
	PaymentUnpaid  PaymentStatus = "unpaid"
	PaymentPartial PaymentStatus = "partial"
)

type StudyStatus string

const (
	StudyLearning  StudyStatus = "learning"
	StudyJobSearch StudyStatus = "job_search"
	StudyOffer     StudyStatus = "offer"
	StudyWorking   StudyStatus = "working"
)

type CreateStudentRequest struct {
	FullName      string        `json:"full_name"`
	Email         string        `json:"email"`
	Telegram      string        `json:"telegram"`
	GroupID       uint          `json:"group_id"`
	TuitionTotal  int           `json:"tuition_total"`
	TuitionPaid   int           `json:"tuition_paid"`
	PaymentStatus PaymentStatus `json:"payment_status"`
	StudyStatus   StudyStatus   `json:"study_status"`
}

type UpdateStudentRequest struct {
	FullName      *string        `json:"full_name"`
	Email         *string        `json:"email"`
	Telegram      *string        `json:"telegram"`
	GroupID       *uint          `json:"group_id"`
	TuitionTotal  *int           `json:"tuition_total"`
	TuitionPaid   *int           `json:"tuition_paid"`
	PaymentStatus *PaymentStatus `json:"payment_status"`
	StudyStatus   *StudyStatus   `json:"study_status"`
}
