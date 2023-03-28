package entities

import "time"

type JwtToken struct {
	UserID    int64
	Exp       int64
	Token     string
	IsActive  bool   `gorm:"default:true"`
	Device    string `json:"device"`
	CreatedAt time.Time
	UpdatedAt time.Time
}
