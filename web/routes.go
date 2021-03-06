package main

import (
	"net/http"

	"github.com/gemm123/reservasi-web/internal/config"
	"github.com/gemm123/reservasi-web/internal/handlers"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func routes(app *config.AppConfig) http.Handler {
	mux := chi.NewRouter()

	mux.Use(middleware.Recoverer)
	mux.Use(NoSurf)
	mux.Use(SessionLoad)

	fileServer := http.FileServer(http.Dir("./static/"))
	mux.Handle("/static/*", http.StripPrefix("/static", fileServer))

	mux.Get("/", handlers.Repo.Home)
	mux.Get("/about", handlers.Repo.About)

	mux.Route("/room", func(mux chi.Router) {
		mux.Get("/president", handlers.Repo.President)
		mux.Get("/royal", handlers.Repo.Royal)
		mux.Get("/tower-club", handlers.Repo.TowerClub)
		mux.Get("/grand-deluxe", handlers.Repo.GrandDeluxe)
		mux.Get("/deluxe", handlers.Repo.Deluxe)
	})

	mux.Post("/search-availability", handlers.Repo.CheckAvailability)
	mux.Get("/book-room", handlers.Repo.BookRoom)
	mux.Get("/make-reservation", handlers.Repo.Reservation)
	mux.Post("/make-reservation", handlers.Repo.PostReservation)
	mux.Get("/reservation-summary", handlers.Repo.ReservationSummary)

	mux.Get("/login", handlers.Repo.ShowLogin)
	mux.Post("/login", handlers.Repo.PostShowLogin)
	mux.Get("/register", handlers.Repo.Register)
	mux.Post("/register", handlers.Repo.PostRegister)
	mux.Get("/logout", handlers.Repo.Logout)

	mux.Route("/admin", func(mux chi.Router) {
		mux.Use(Auth)
		mux.Get("/all-reservation", handlers.Repo.AdminShowAllReservation)
		mux.Get("/all-reservation/{id}", handlers.Repo.AdminShowAllReservationByID)
		mux.Post("/all-reservation/{id}", handlers.Repo.AdminPostAllReservation)

		mux.Get("/new-reservation", handlers.Repo.AdminShowNewReservation)
		mux.Get("/new-reservation/{id}", handlers.Repo.AdminShowNewReservationByID)
		mux.Post("/new-reservation/{id}", handlers.Repo.AdminPostNewReservation)

		mux.Get("/process-reservation/{id}", handlers.Repo.AdminProcessReservation)

		mux.Get("/delete-new-reservation/{id}", handlers.Repo.AdminDeleteNewReservation)
		mux.Get("/delete-all-reservation/{id}", handlers.Repo.AdminDeleteAllReservation)

		mux.Get("/account", handlers.Repo.AdminAccount)
	})

	return mux
}
