package model

// Code ..
type Code struct {
	ID  string `json:"_id" bson:"_id,omitempty"`
	Num int    `json:"num,omitempty" bson:"num,omitempty"`
}

// IsExists struct
func (m Code) IsExists() (ok bool) {
	if m.ID != "" {
		ok = true
	}
	return
}
