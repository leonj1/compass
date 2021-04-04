package main

import (
	"database/sql"
	"embed"
	"fmt"
	"github.com/Boostport/migration"
	"github.com/Boostport/migration/driver/sqlite"
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

// Create migration source
//go:embed migrations
var embedFS embed.FS

func main() {
	log.Info().Msgf("Starting compass")
	sqliteDatabase := os.Getenv("DB_PATH")
	serverPort := os.Getenv("HTTP_PORT")
	version := os.Getenv("VERSION")

	if _, err := os.Stat(sqliteDatabase); err == nil {
		log.Info().Msgf("Exists: %s", sqliteDatabase)
	} else if os.IsNotExist(err) {
		log.Info().Msgf("DB file does not exist: %s", sqliteDatabase)
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

	embedSource := &migration.EmbedMigrationSource{
		EmbedFS: embedFS,
		Dir:     "migrations",
	}

	driver, err := sqlite.New(fmt.Sprintf("file:%s", sqliteDatabase), true)
	if err != nil {
		log.Printf("Problem connecting to db: %s", err.Error())
		os.Exit(1)
	}

	// Run all up migrations
	applied, err := migration.Migrate(driver, embedSource, migration.Up, 0)
	if err != nil {
		log.Printf("Problem applying migrations: %s", err.Error())
		os.Exit(1)
	}
	log.Info().Msgf("Applied %d db migrations", applied)

	app := routes.App{
		Compass: *services.NewCompass(db),
		Version: version,
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
