package model

import "time"

// JSONInfo ..
type JSONInfo struct {
	ID          string    `json:"_id" bson:"_id,omitempty"`
	Name        string    `json:"name,omitempty" bson:"name,omitempty"`
	Data        string    `json:"data,omitempty" bson:"data,omitempty"` // base64 encode
	CreatedDate time.Time `json:"created_date,omitempty" bson:"created_date,omitempty"`
}

// IsExists struct
func (m JSONInfo) IsExists() (ok bool) {
	if m.ID != "" {
		ok = true
	}
	return
}
