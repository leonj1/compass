package routes

import (
	"github.com/gorilla/mux"
	"github.com/leonj1/compass/exceptions"
	"github.com/rs/zerolog"
	"net/http"
	"os"
)

func (a *App) SetAppVersion(w http.ResponseWriter, r *http.Request) {
	logger := zerolog.New(os.Stdout).With().Logger()
	logger.Info().Msg("Set App Version invoked")
	vars := mux.Vars(r)
	name := vars["name"]
	environment := vars["environment"]
	version := vars["version"]
	_, err := a.Compass.SetApplicationEnv(name, environment, version)
	if err != nil {
		if perr, ok := err.(*exceptions.NotFound); ok {
			RespondAndLog(w, perr.Error(), TEXT, http.StatusNotFound, logger)
			return
		}
		RespondAndLog(w, err.Error(), TEXT, http.StatusInternalServerError, logger)
		return
	}
	RespondAndDontLog(w, "OK", TEXT, http.StatusOK, logger)
}
