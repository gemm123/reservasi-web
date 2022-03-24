package main

import (
	"encoding/gob"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/gemm123/reservasi-web/internal/config"
	"github.com/gemm123/reservasi-web/internal/driver"
	"github.com/gemm123/reservasi-web/internal/handlers"
	"github.com/gemm123/reservasi-web/internal/helpers"
	"github.com/gemm123/reservasi-web/internal/models"
	"github.com/gemm123/reservasi-web/internal/render"
)

var app config.AppConfig
var session *scs.SessionManager

func run() (*driver.DB, error) {
	//model data yang ditaruh di session
	gob.Register(models.Reservation{})
	gob.Register(models.User{})
	gob.Register(models.Room{})
	gob.Register(models.Restriction{})
	gob.Register(map[string]int{})

	//ganti true kalau di production
	app.InProduction = false
	app.UseCache = false

	app.InfoLog = log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	app.ErrorLog = log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	session = scs.New()
	session.Lifetime = 24 * time.Hour
	session.Cookie.Persist = true
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = app.InProduction

	app.Session = session

	//koneksi ke database
	connectionString := "host=localhost port=5432 dbname=reservasi_web user=postgres password=gemmq123456"
	db, err := driver.ConnectSQL(connectionString)
	if err != nil {
		log.Fatal("cant connect to database")
	}
	log.Println("connected to database")

	templateCache, err := render.CreateTemplateCache()
	if err != nil {
		log.Fatal("cannot create template cache")
		return nil, err
	}

	app.TemplateCache = templateCache

	repo := handlers.NewRepo(&app, db)
	handlers.NewHandlers(repo)
	render.NewRenderer(&app)
	helpers.NewHelpers(&app)

	return db, nil
}

func main() {
	db, err := run()
	if err != nil {
		log.Fatal(err)
	}
	defer db.SQL.Close()

	fmt.Println("starting application on port 8080")

	server := &http.Server{
		Addr:    ":8080",
		Handler: routes(&app),
	}

	err = server.ListenAndServe()
	log.Fatal(err)
}
