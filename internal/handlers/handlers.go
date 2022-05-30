package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gemm123/reservasi-web/internal/config"
	"github.com/gemm123/reservasi-web/internal/driver"
	"github.com/gemm123/reservasi-web/internal/forms"
	"github.com/gemm123/reservasi-web/internal/helpers"
	"github.com/gemm123/reservasi-web/internal/models"
	"github.com/gemm123/reservasi-web/internal/render"
	"github.com/gemm123/reservasi-web/internal/repository"
	"github.com/gemm123/reservasi-web/internal/repository/dbrepo"
	"github.com/go-chi/chi/v5"
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

//menampilkan page home
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

type jsonResponse struct {
	OK        bool   `json:"ok"`
	Message   string `json:"message"`
	RoomID    string `json:"room_id"`
	StartDate string `json:"start_date"`
	EndDate   string `json:"end_date"`
}

func (repo *Repository) CheckAvailability(writer http.ResponseWriter, request *http.Request) {
	err := request.ParseForm()
	if err != nil {
		response := jsonResponse{
			OK:      false,
			Message: "Internal server error",
		}

		out, _ := json.Marshal(response)
		writer.Header().Add("Content-Type", "application/json")
		writer.Write(out)
		return
	}

	sd := request.Form.Get("start")
	ed := request.Form.Get("end")

	layout := "2006-01-02"
	startDate, _ := time.Parse(layout, sd)
	endDate, _ := time.Parse(layout, ed)
	roomID, _ := strconv.Atoi(request.Form.Get("room_id"))

	available, err := repo.DB.SearchAvailabilityByDatesAndRoomID(startDate, endDate, roomID)
	if err != nil {
		response := jsonResponse{
			OK:      false,
			Message: "Error queying database",
		}

		out, _ := json.Marshal(response)
		writer.Header().Add("Content-Type", "application/json")
		writer.Write(out)
		return
	}

	response := jsonResponse{
		OK:        available,
		Message:   "",
		StartDate: sd,
		EndDate:   ed,
		RoomID:    strconv.Itoa(roomID),
	}

	out, _ := json.Marshal(response)
	writer.Header().Add("Content-Type", "application/json")
	writer.Write(out)
}

func (repo *Repository) BookRoom(writer http.ResponseWriter, request *http.Request) {
	roomID, _ := strconv.Atoi(request.URL.Query().Get("id"))
	sd := request.URL.Query().Get("s")
	ed := request.URL.Query().Get("e")

	layout := "2006-01-02"
	startDate, _ := time.Parse(layout, sd)
	endDate, _ := time.Parse(layout, ed)

	var reservation models.Reservation

	room, err := repo.DB.GetRoomByID(roomID)
	if err != nil {
		repo.App.Session.Put(request.Context(), "error", "Gagal mendapatkan room dari database!")
		http.Redirect(writer, request, "/", http.StatusSeeOther)
		return
	}

	reservation.StartDate = startDate
	reservation.EndDate = endDate
	reservation.RoomID = roomID
	reservation.Room.RoomName = room.RoomName

	repo.App.Session.Put(request.Context(), "reservation", reservation)

	http.Redirect(writer, request, "/make-reservation", http.StatusSeeOther)
}

func (repo *Repository) Reservation(writer http.ResponseWriter, request *http.Request) {
	reservation, ok := repo.App.Session.Get(request.Context(), "reservation").(models.Reservation)
	if !ok {
		repo.App.Session.Put(request.Context(), "error", "Gagal mendapatkan session reservation")
		http.Redirect(writer, request, "/", http.StatusSeeOther)
		return
	}

	sd := reservation.StartDate.Format("2006-01-02")
	ed := reservation.EndDate.Format("2006-01-02")

	stringMap := make(map[string]string)
	stringMap["start_date"] = sd
	stringMap["end_date"] = ed

	data := make(map[string]interface{})
	data["reservation"] = reservation

	render.Template(writer, request, "make-reservation.page.html", &models.TemplateData{
		Form:      forms.New(nil),
		StringMap: stringMap,
		Data:      data,
	})
}

func (repo *Repository) PostReservation(writer http.ResponseWriter, request *http.Request) {
	res, ok := repo.App.Session.Get(request.Context(), "reservation").(models.Reservation)
	if !ok {
		repo.App.Session.Put(request.Context(), "error", "Gagal mendapatkan session reservation")
		http.Redirect(writer, request, "/", http.StatusSeeOther)
		return
	}

	err := request.ParseForm()
	if err != nil {
		repo.App.Session.Put(request.Context(), "error", "gagal parse form")
		http.Redirect(writer, request, "/", http.StatusSeeOther)
		return
	}

	sd := request.Form.Get("start_date")
	ed := request.Form.Get("end_date")
	layout := "2006-01-02"
	startDate, _ := time.Parse(layout, sd)
	endDate, _ := time.Parse(layout, ed)

	roomID, err := strconv.Atoi(request.Form.Get("room_id"))
	if err != nil {
		repo.App.Session.Put(request.Context(), "error", "data salah")
		http.Redirect(writer, request, "/", http.StatusSeeOther)
		return
	}

	reservation := models.Reservation{
		Name:      request.Form.Get("name"),
		Email:     request.Form.Get("email"),
		Phone:     request.Form.Get("phone"),
		StartDate: startDate,
		EndDate:   endDate,
		RoomID:    roomID,
		Room:      res.Room,
	}

	form := forms.New(request.PostForm)
	form.Required("name", "email", "phone")
	form.MinLength("name", 3)
	form.IsEmail("email")

	if !form.Valid() {
		data := make(map[string]interface{})
		data["reservation"] = reservation
		stringMap := make(map[string]string)
		stringMap["start_date"] = sd
		stringMap["end_date"] = ed

		render.Template(writer, request, "make-reservation.page.html", &models.TemplateData{
			Form:      form,
			Data:      data,
			StringMap: stringMap,
		})

		return
	}

	reservationID, err := repo.DB.InsertReservation(reservation)
	if err != nil {
		repo.App.Session.Put(request.Context(), "error", "Gagal insert ke database")
		http.Redirect(writer, request, "/", http.StatusSeeOther)
		return
	}

	roomRestriction := models.RoomRestrictions{
		StartDate:     startDate,
		EndDate:       endDate,
		RoomID:        roomID,
		ReservationID: reservationID,
		RestrictionID: 1,
	}

	err = repo.DB.InsertRoomRestriction(roomRestriction)
	if err != nil {
		repo.App.Session.Put(request.Context(), "error", "Gagal insert ke database")
		http.Redirect(writer, request, "/", http.StatusSeeOther)
		return
	}

	repo.App.Session.Put(request.Context(), "reservation", reservation)
	http.Redirect(writer, request, "/reservation-summary", http.StatusSeeOther)
}

func (repo *Repository) ReservationSummary(writer http.ResponseWriter, request *http.Request) {
	reservation, ok := repo.App.Session.Get(request.Context(), "reservation").(models.Reservation)
	if !ok {
		repo.App.Session.Put(request.Context(), "error", "Gagal mendapatkan session reservation")
		http.Redirect(writer, request, "/", http.StatusSeeOther)
		return
	}

	repo.App.Session.Remove(request.Context(), "reservation")

	data := make(map[string]interface{})
	data["reservation"] = reservation

	startDate := reservation.StartDate.Format("2006-01-02")
	endDate := reservation.EndDate.Format("2006-01-02")
	stringMap := make(map[string]string)
	stringMap["start_date"] = startDate
	stringMap["end_date"] = endDate

	render.Template(writer, request, "reservation-summary.page.html", &models.TemplateData{
		Data:      data,
		StringMap: stringMap,
	})
}

func (repo *Repository) Register(writer http.ResponseWriter, request *http.Request) {
	render.Template(writer, request, "register.page.html", &models.TemplateData{
		Form: forms.New(nil),
	})
}

//proses register
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
		repo.App.Session.Put(request.Context(), "error", "Register failed!")
		http.Redirect(writer, request, "/register", http.StatusSeeOther)
		return
	}

	repo.App.Session.Put(request.Context(), "success", "Register success")
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
		repo.App.Session.Put(request.Context(), "error", "Email atau password wrong")
		http.Redirect(writer, request, "/login", http.StatusSeeOther)
		return
	}

	repo.App.Session.Put(request.Context(), "user_id", id)
	repo.App.Session.Put(request.Context(), "success", "Success login")
	http.Redirect(writer, request, "/", http.StatusSeeOther)
}

//proses logout
func (repo *Repository) Logout(writer http.ResponseWriter, request *http.Request) {
	_ = repo.App.Session.Destroy(request.Context())
	_ = repo.App.Session.RenewToken(request.Context())
	http.Redirect(writer, request, "/login", http.StatusSeeOther)
}

func (repo *Repository) AdminShowAllReservation(writer http.ResponseWriter, request *http.Request) {
	reservations, err := repo.DB.ShowAllReservation()
	if err != nil {
		helpers.ServerError(writer, err)
		return
	}

	data := make(map[string]interface{})
	data["reservations"] = reservations

	render.Template(writer, request, "admin-all-reservation.page.html", &models.TemplateData{
		Data: data,
	})
}

func (repo *Repository) AdminShowAllReservationByID(writer http.ResponseWriter, request *http.Request) {
	splitURL := strings.Split(request.RequestURI, "/")
	id, err := strconv.Atoi(splitURL[3])
	if err != nil {
		helpers.ServerError(writer, err)
		return
	}

	reservation, err := repo.DB.GetReservationByID(id)
	if err != nil {
		helpers.ServerError(writer, err)
		return
	}

	data := make(map[string]interface{})
	data["reservation"] = reservation

	render.Template(writer, request, "admin-show-all-reservation.page.html", &models.TemplateData{
		Data: data,
		Form: forms.New(nil),
	})
}

func (repo *Repository) AdminShowNewReservation(writer http.ResponseWriter, request *http.Request) {
	reservations, err := repo.DB.ShowNewReservation()
	if err != nil {
		helpers.ServerError(writer, err)
		return
	}

	data := make(map[string]interface{})
	data["reservations"] = reservations

	render.Template(writer, request, "admin-new-reservation.page.html", &models.TemplateData{
		Data: data,
	})
}

func (repo *Repository) AdminShowNewReservationByID(writer http.ResponseWriter, request *http.Request) {
	splitURL := strings.Split(request.RequestURI, "/")
	id, err := strconv.Atoi(splitURL[3])
	if err != nil {
		helpers.ServerError(writer, err)
		return
	}

	reservation, err := repo.DB.GetReservationByID(id)
	if err != nil {
		helpers.ServerError(writer, err)
		return
	}

	data := make(map[string]interface{})
	data["reservation"] = reservation

	render.Template(writer, request, "admin-show-new-reservation.page.html", &models.TemplateData{
		Data: data,
		Form: forms.New(nil),
	})
}

func (repo *Repository) AdminProcessReservation(writer http.ResponseWriter, request *http.Request) {
	id, _ := strconv.Atoi(chi.URLParam(request, "id"))
	err := repo.DB.UpdateProcessedReservation(id)
	if err != nil {
		helpers.ServerError(writer, err)
		return
	}

	http.Redirect(writer, request, "/admin/new-reservation", http.StatusSeeOther)
}

func (repo *Repository) AdminDeleteNewReservation(writer http.ResponseWriter, request *http.Request) {
	id, _ := strconv.Atoi(chi.URLParam(request, "id"))
	_ = repo.DB.DeleteReservationByID(id)

	repo.App.Session.Put(request.Context(), "success", "Reservation deleted")

	http.Redirect(writer, request, "/admin/new-reservation", http.StatusSeeOther)
}

func (repo *Repository) AdminDeleteAllReservation(writer http.ResponseWriter, request *http.Request) {
	id, _ := strconv.Atoi(chi.URLParam(request, "id"))
	_ = repo.DB.DeleteReservationByID(id)

	repo.App.Session.Put(request.Context(), "success", "Reservation deleted")

	http.Redirect(writer, request, "/admin/all-reservation", http.StatusSeeOther)
}

func (repo *Repository) AdminPostNewReservation(writer http.ResponseWriter, request *http.Request) {
	err := request.ParseForm()
	if err != nil {
		helpers.ServerError(writer, err)
		return
	}

	id, _ := strconv.Atoi(chi.URLParam(request, "id"))

	reservation, err := repo.DB.GetReservationByID(id)
	if err != nil {
		helpers.ServerError(writer, err)
		return
	}

	reservation.Name = request.Form.Get("name")
	reservation.Email = request.Form.Get("email")
	reservation.Phone = request.Form.Get("phone")

	err = repo.DB.UpdateReservation(reservation)
	if err != nil {
		helpers.ServerError(writer, err)
		return
	}

	repo.App.Session.Put(request.Context(), "success", "Changes saved")

	http.Redirect(writer, request, "/admin/new-reservation", http.StatusSeeOther)
}

func (repo *Repository) AdminPostAllReservation(writer http.ResponseWriter, request *http.Request) {
	err := request.ParseForm()
	if err != nil {
		helpers.ServerError(writer, err)
		return
	}

	id, _ := strconv.Atoi(chi.URLParam(request, "id"))

	reservation, err := repo.DB.GetReservationByID(id)
	if err != nil {
		helpers.ServerError(writer, err)
		return
	}

	reservation.Name = request.Form.Get("name")
	reservation.Email = request.Form.Get("email")
	reservation.Phone = request.Form.Get("phone")

	err = repo.DB.UpdateReservation(reservation)
	if err != nil {
		helpers.ServerError(writer, err)
		return
	}

	repo.App.Session.Put(request.Context(), "success", "Changes saved")

	http.Redirect(writer, request, "/admin/all-reservation", http.StatusSeeOther)
}

func (repo *Repository) AdminAccount(writer http.ResponseWriter, request *http.Request) {
	render.Template(writer, request, "admin-account.page.html", &models.TemplateData{})
}
