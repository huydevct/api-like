package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// ReportService : Lưu số luợng report service do mobile bắn lên
type ReportService struct {
	ID           primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	ServiceCode  string             `json:"service_code" bson:"service_code"`
	NumberReport int                `json:"number_report" bson:"number_report"` // Tổng sổ lượt
	// trace
	CreatedDate *time.Time `json:"created_date,omitempty" bson:"created_date,omitempty"`
	UpdatedDate *time.Time `json:"updated_date,omitempty" bson:"updated_date,omitempty"`
	ResumeDate  *time.Time `json:"resume_date,omitempty" bson:"resume_date,omitempty"` // thời gian userautolike resume lại gói
}

// IsExists ..
func (m ReportService) IsExists() (ok bool) {
	if !m.ID.IsZero() {
		ok = true
	}
	return
}
