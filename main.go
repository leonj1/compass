package main

import (
	"database/sql"
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/sqlite3"
	"github.com/golang-migrate/migrate/v4/source/file"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/leonj1/compass/routes"
	"github.com/leonj1/compass/services"
	_ "github.com/mattn/go-sqlite3"
	"github.com/rs/zerolog/log"
	"net/http"
	"os"
	"time"
)

func main() {
	log.Info().Msgf("Starting compass")
	sqliteDatabase := os.Getenv("DB_PATH")
	serverPort := os.Getenv("HTTP_PORT")

	if _, err := os.Stat(sqliteDatabase); err == nil {
		log.Info().Msgf("Exists: %s", sqliteDatabase)
	} else if os.IsNotExist(err) {
		log.Info().Msgf("Does not exist: %s", sqliteDatabase)
		os.Exit(1)
	} else {
		// Schrodinger: file may or may not exist. See err for details.
		// Therefore, do *NOT* use !os.IsNotExist(err) to test for file existence
		log.Info().Msgf("Problem checking if file exists: %s", sqliteDatabase)
		os.Exit(1)
	}

	log.Info().Msgf("Attempting to connect to db %s", sqliteDatabase)
	db, err := sql.Open("sqlite3", sqliteDatabase)
	if err != nil {
		log.Error().Msgf("Problem connecting to db %s", err.Error())
		os.Exit(1)
	}
	defer func() {
		if err := db.Close(); err != nil {
			log.Error().Msgf("Problem closing db connection %s", err.Error())
			return
		}
	}()

	config := &sqlite3.Config{
		MigrationsTable: "my_migration_table",
	}
	driver, err := sqlite3.WithInstance(db, config)
	if err != nil {
		log.Printf(err.Error())
		os.Exit(1)
	}

	fsrc, err := (&file.File{}).Open("file://.//migrations")
	if err != nil {
		log.Printf("Problem %s", err.Error())
		os.Exit(1)
	}

	m, err := migrate.NewWithInstance(
		"file",
		fsrc,
		"sqlite3",
		driver,
	)

	log.Printf("Performing database migrations\n")
	err = m.Up()
	if err != nil {
		if err.Error() == "no change" {
			log.Printf("No migrations performed\n")
		} else {
			log.Printf("Migrate UP failed: %s", err.Error())
			os.Exit(1)
		}
	}

	app := routes.App{
		Compass: *services.NewCompass(db),
	}

	s := mux.NewRouter()
	s.Handle("/health", http.HandlerFunc(app.Health)).Methods(routes.GET)
	s.Handle("/applications", http.HandlerFunc(app.FetchAll)).Methods(routes.GET)
	s.Handle("/applications/{name}/{environment}", http.HandlerFunc(app.AppGetVersion)).Methods(routes.GET)
	s.Handle("/applications/{name}/{environment}/{version}", http.HandlerFunc(app.SetAppVersion)).Methods(routes.PUT)

	log.Info().Msgf("Starting Web Server on %s", serverPort)
	port := fmt.Sprintf(":%s", serverPort)
	ss := &http.Server{
		Addr: port,
		Handler: handlers.CORS(
			handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization"}),
			handlers.AllowedMethods([]string{routes.GET, routes.POST, routes.PUT, routes.DELETE, routes.HEAD, routes.OPTIONS}),
			handlers.AllowedOrigins([]string{"*"}))(s),
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 2 << 20,
	}
	err = ss.ListenAndServe()
	if err != nil {
		log.Error().Msgf("Problem with web server: %s", err.Error())
		os.Exit(1)
	}
	os.Exit(0)
}
