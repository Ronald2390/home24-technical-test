package model

import "time"

// Session represent user's session
type Session struct {
	ID        string
	Type      string
	ExpiredAt time.Time
	Info      map[string]interface{}
	User      *User
}
