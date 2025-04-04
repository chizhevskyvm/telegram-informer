package domain

import (
	"time"
)

type Event struct {
	ID           int
	UserID       int
	TypeID       string
	Title        string
	Notification time.Time
	TimeToNotify time.Time
}
