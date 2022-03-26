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

	mux.Get("/", handlers.Repo.Home)

	mux.Route("/room", func(mux chi.Router) {
		mux.Get("/kopi", handlers.Repo.Kopi)
		mux.Get("/teh", handlers.Repo.Teh)
		mux.Get("/susu", handlers.Repo.Susu)
		mux.Get("/jahe", handlers.Repo.Jahe)
		mux.Get("/jus", handlers.Repo.Jus)
	})

	return mux
}
