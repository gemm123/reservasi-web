package handlers

import (
	"net/http"

	"github.com/gemm123/reservasi-web/internal/config"
	"github.com/gemm123/reservasi-web/internal/driver"
	"github.com/gemm123/reservasi-web/internal/models"
	"github.com/gemm123/reservasi-web/internal/render"
	"github.com/gemm123/reservasi-web/internal/repository"
	"github.com/gemm123/reservasi-web/internal/repository/dbrepo"
)

//repo yang digunakan oleh handlers
var Repo *Repository

type Repository struct {
	App *config.AppConfig
	DB  repository.DatabaseRepo
}

//membuat repo baru
func NewRepo(app *config.AppConfig, db *driver.DB) *Repository {
	return &Repository{
		App: app,
		DB:  dbrepo.NewPostgresRepo(db.SQL, app),
	}
}

//set repo untuk handler
func NewHandlers(repo *Repository) {
	Repo = repo
}

func (repo *Repository) Home(writer http.ResponseWriter, request *http.Request) {
	render.Template(writer, request, "home.page.html", &models.TemplateData{})
}
