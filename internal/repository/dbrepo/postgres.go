package dbrepo

import (
	"context"
	"errors"
	"log"
	"time"

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
