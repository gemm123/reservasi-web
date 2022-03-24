package main

import (
	"net/http"

	"github.com/justinas/nosurf"
)

//menambahkan CSRF protection di semua post request
func NoSurf(next http.Handler) http.Handler {
	csrfHandler := nosurf.New(next)

	csrfHandler.SetBaseCookie(http.Cookie{
		HttpOnly: true,
		Path:     "/",
		Secure:   app.InProduction,
		SameSite: http.SameSiteLaxMode,
	})

	return csrfHandler
}

//load dan save session di setiap request
func SessionLoad(next http.Handler) http.Handler {
	return session.LoadAndSave(next)
}
