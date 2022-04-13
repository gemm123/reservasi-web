package handlers

import (
	"log"
	"net/http"

	"github.com/gemm123/reservasi-web/internal/config"
	"github.com/gemm123/reservasi-web/internal/driver"
	"github.com/gemm123/reservasi-web/internal/forms"
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

func (repo *Repository) About(writer http.ResponseWriter, request *http.Request) {
	render.Template(writer, request, "about.page.html", &models.TemplateData{})
}

func (repo *Repository) President(writer http.ResponseWriter, request *http.Request) {
	render.Template(writer, request, "president.page.html", &models.TemplateData{})
}

func (repo *Repository) Royal(writer http.ResponseWriter, request *http.Request) {
	render.Template(writer, request, "royal.page.html", &models.TemplateData{})
}

func (repo *Repository) TowerClub(writer http.ResponseWriter, request *http.Request) {
	render.Template(writer, request, "tower-club.page.html", &models.TemplateData{})
}

func (repo *Repository) GrandDeluxe(writer http.ResponseWriter, request *http.Request) {
	render.Template(writer, request, "grand-deluxe.page.html", &models.TemplateData{})
}

func (repo *Repository) Deluxe(writer http.ResponseWriter, request *http.Request) {
	render.Template(writer, request, "deluxe.page.html", &models.TemplateData{})
}

func (repo *Repository) ShowLogin(writer http.ResponseWriter, request *http.Request) {
	render.Template(writer, request, "login.page.html", &models.TemplateData{
		Form: forms.New(nil),
	})
}

func (repo *Repository) Register(writer http.ResponseWriter, request *http.Request) {
	render.Template(writer, request, "register.page.html", &models.TemplateData{
		Form: forms.New(nil),
	})
}

func (repo *Repository) PostRegister(writer http.ResponseWriter, request *http.Request) {
	err := request.ParseForm()
	if err != nil {
		repo.App.Session.Put(request.Context(), "error", "Gagal parse form!")
		http.Redirect(writer, request, "/register", http.StatusSeeOther)
		return
	}

	name := request.Form.Get("name")
	email := request.Form.Get("email")
	password := request.Form.Get("password")

	form := forms.New(request.PostForm)
	form.Required("name", "email", "password")
	form.MinLength("name", 3)
	form.IsEmail("email")

	if !form.Valid() {
		render.Template(writer, request, "register.page.html", &models.TemplateData{
			Form: form,
		})
		return
	}

	err = repo.DB.InsertUser(name, email, password)
	if err != nil {
		repo.App.Session.Put(request.Context(), "error", "Gagal register user!")
		http.Redirect(writer, request, "/register", http.StatusSeeOther)
		return
	}

	repo.App.Session.Put(request.Context(), "success", "Register berhasil")
	http.Redirect(writer, request, "/login", http.StatusSeeOther)
}

func (repo *Repository) PostShowLogin(writer http.ResponseWriter, request *http.Request) {
	_ = repo.App.Session.RenewToken(request.Context())

	err := request.ParseForm()
	if err != nil {
		log.Println(err)
	}

	email := request.Form.Get("email")
	password := request.Form.Get("password")

	form := forms.New(request.PostForm)
	form.Required("email", "password")
	form.IsEmail("email")

	if !form.Valid() {
		render.Template(writer, request, "login.page.html", &models.TemplateData{
			Form: form,
		})
		return
	}

	id, _, err := repo.DB.Authenticated(email, password)
	if err != nil {
		repo.App.Session.Put(request.Context(), "error", "Email atau password salah")
		http.Redirect(writer, request, "/login", http.StatusSeeOther)
		return
	}

	repo.App.Session.Put(request.Context(), "user_id", id)
	repo.App.Session.Put(request.Context(), "success", "Berhasil login")
	http.Redirect(writer, request, "/", http.StatusSeeOther)
}

func (repo *Repository) Logout(writer http.ResponseWriter, request *http.Request) {
	_ = repo.App.Session.Destroy(request.Context())
	_ = repo.App.Session.RenewToken(request.Context())
	http.Redirect(writer, request, "/login", http.StatusSeeOther)
}

func (repo *Repository) AdminDashboard(writer http.ResponseWriter, request *http.Request) {
	render.Template(writer, request, "admin-dashboard.page.html", &models.TemplateData{})
}
