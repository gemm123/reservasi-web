package helpers

import "github.com/gemm123/reservasi-web/internal/config"

var app *config.AppConfig

func NewHelpers(a *config.AppConfig) {
	app = a
}
