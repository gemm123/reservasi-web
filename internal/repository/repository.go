package repository

import "time"

type DatabaseRepo interface {
	InsertUser(name, email, password string) error
	Authenticated(email, password string) (int, string, error)
	SearchAvailabilityByDatesAndRoomID(start, end time.Time, roomID int) (bool, error)
}
