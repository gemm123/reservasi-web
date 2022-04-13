package helpers

import (
	"net/http"

	"github.com/gemm123/reservasi-web/internal/config"
)

var app *config.AppConfig

//setup app config untuk helpers
func NewHelpers(a *config.AppConfig) {
	app = a
}

func IsAuthenticated(request *http.Request) bool {
	exist := app.Session.Exists(request.Context(), "user_id")
	return exist
}
