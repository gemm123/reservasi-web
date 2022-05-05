package dbrepo

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/gemm123/reservasi-web/internal/models"
	"golang.org/x/crypto/bcrypt"
)

//insert ke tabel user
func (repo *postgresDBRepo) InsertUser(name, email, password string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), 12)

	query := `insert into users(name, email, password, created_at, updated_at) 
				values($1, $2, $3, $4, $5)`

	_, err := repo.DB.ExecContext(ctx, query, name, email, hashedPassword, time.Now(), time.Now())
	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}

//authenticated user
func (repo *postgresDBRepo) Authenticated(email, password string) (int, string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var id int
	var hashedPassword string

	row := repo.DB.QueryRowContext(ctx, "select id, password from users where email = $1", email)
	err := row.Scan(&id, &hashedPassword)
	if err != nil {
		return id, "", err
	}

	err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	if err == bcrypt.ErrMismatchedHashAndPassword {
		return 0, "", errors.New("password salah")
	} else if err != nil {
		return 0, "", err
	}

	return id, hashedPassword, nil
}

//return true jika room tersedia
func (repo *postgresDBRepo) SearchAvailabilityByDatesAndRoomID(start, end time.Time, roomID int) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `select count(id)
				from room_restrictions
				where room_id = $1 and $2 < end_date and $3 > start_date;`

	var numRows int
	row := repo.DB.QueryRowContext(ctx, query, roomID, start, end)
	err := row.Scan(&numRows)
	if err != nil {
		return false, err
	}

	if numRows == 0 {
		return true, nil
	}

	return false, err
}

//get room by id
func (repo *postgresDBRepo) GetRoomByID(roomID int) (models.Room, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `select id, room_name, created_at, updated_at from rooms where id = $1`

	var room models.Room
	row := repo.DB.QueryRowContext(ctx, query, roomID)
	err := row.Scan(
		&room.ID,
		&room.RoomName,
		&room.CreatedAt,
		&room.UpdatedAt,
	)
	if err != nil {
		return room, err
	}

	return room, nil
}

//insert reservasi ke database
func (repo *postgresDBRepo) InsertReservation(reservation models.Reservation) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var reservationID int

	query := `insert into reservations (name, email, phone, start_date, end_date, room_id, created_at, updated_at)
				values($1, $2, $3, $4, $5, $6, $7, $8) returning id`

	err := repo.DB.QueryRowContext(ctx, query,
		reservation.Name,
		reservation.Email,
		reservation.Phone,
		reservation.StartDate,
		reservation.EndDate,
		reservation.RoomID,
		reservation.CreatedAt,
		reservation.UpdatedAt,
	).Scan(&reservationID)
	if err != nil {
		return 0, err
	}

	return reservationID, nil
}

//insert room restriction ke database
func (repo *postgresDBRepo) InsertRoomRestriction(roomRestriction models.RoomRestrictions) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `insert into room_restrictions
				(start_date, end_date, room_id, reservation_id, restriction_id, created_at, updated_at)
				values($1, $2, $3, $4, $5, $6, $7)`

	_, err := repo.DB.ExecContext(ctx, query,
		roomRestriction.StartDate,
		roomRestriction.EndDate,
		roomRestriction.RoomID,
		roomRestriction.ReservationID,
		roomRestriction.RestrictionID,
		time.Now(),
		time.Now(),
	)
	if err != nil {
		return err
	}

	return nil
}

//return slice resersevation
func (repo *postgresDBRepo) ShowAllReservation() ([]models.Reservation, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var resersevations []models.Reservation

	query := `select res.id, res.name, res.start_date, res.end_date, room.room_name
				from reservations res
				left join rooms room on (res.room_id = room.id)
				order by res.start_date asc`

	rows, err := repo.DB.QueryContext(ctx, query)
	if err != nil {
		return resersevations, err
	}

	defer rows.Close()

	for rows.Next() {
		var resersevation models.Reservation
		err := rows.Scan(
			&resersevation.ID,
			&resersevation.Name,
			&resersevation.StartDate,
			&resersevation.EndDate,
			&resersevation.Room.RoomName,
		)
		if err != nil {
			return resersevations, err
		}

		resersevations = append(resersevations, resersevation)
	}

	if err = rows.Err(); err != nil {
		return resersevations, err
	}

	return resersevations, nil
}

func (repo *postgresDBRepo) GetReservationByID(id int) (models.Reservation, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `select res.id, res.name, res.email, res.phone, res.start_date, res.end_date, res.room_id, 
				res.processed, res.created_at, res.updated_at, room.id, room.room_name 
				from reservations res
				left join rooms room on (res.room_id = room.id)
				where res.id = $1`

	var resersevation models.Reservation

	row := repo.DB.QueryRowContext(ctx, query, id)
	err := row.Scan(
		&resersevation.ID,
		&resersevation.Name,
		&resersevation.Email,
		&resersevation.Phone,
		&resersevation.StartDate,
		&resersevation.EndDate,
		&resersevation.RoomID,
		&resersevation.Processed,
		&resersevation.CreatedAt,
		&resersevation.UpdatedAt,
		&resersevation.Room.ID,
		&resersevation.Room.RoomName,
	)
	if err != nil {
		return resersevation, err
	}

	return resersevation, nil
}
