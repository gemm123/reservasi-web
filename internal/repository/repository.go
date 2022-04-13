package repository

type DatabaseRepo interface {
	InsertUser(name, email, password string) error
	Authenticated(email, password string) (int, string, error)
}
