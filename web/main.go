package main

import (
	"encoding/gob"
	"log"
	"os"
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/gemm123/reservasi-web/internal/config"
	"github.com/gemm123/reservasi-web/internal/driver"
	"github.com/gemm123/reservasi-web/internal/models"
)

var app config.AppConfig
var session *scs.SessionManager

func main() {

}

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
}
