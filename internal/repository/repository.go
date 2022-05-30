package repository

import (
	"time"

	"github.com/gemm123/reservasi-web/internal/models"
)

type DatabaseRepo interface {
	InsertUser(name, email, password string) error
	Authenticated(email, password string) (int, string, error)
	SearchAvailabilityByDatesAndRoomID(start, end time.Time, roomID int) (bool, error)
	GetRoomByID(roomID int) (models.Room, error)
	InsertReservation(reservation models.Reservation) (int, error)
	InsertRoomRestriction(roomRestriction models.RoomRestrictions) error
	ShowAllReservation() ([]models.Reservation, error)
	ShowNewReservation() ([]models.Reservation, error)
	GetReservationByID(id int) (models.Reservation, error)
	UpdateProcessedReservation(id int) error
	DeleteReservationByID(id int) error
	UpdateReservation(reservation models.Reservation) error
}
