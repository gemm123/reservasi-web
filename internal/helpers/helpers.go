package helpers

import (
	"fmt"
	"net/http"
	"runtime/debug"

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

func ServerError(writer http.ResponseWriter, err error) {
	trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())
	app.ErrorLog.Println(trace)
	http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}
