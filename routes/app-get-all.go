package routes

import (
	"encoding/json"
	"github.com/leonj1/compass/exceptions"
	"github.com/rs/zerolog"
	"net/http"
	"os"
)

func (a *App) FetchAll(w http.ResponseWriter, r *http.Request) {
	logger := zerolog.New(os.Stdout).With().Logger()
	logger.Info().Msg("Get All Apps invoked")
	apps, err := a.Compass.FetchAll()
	if err != nil {
		if perr, ok := err.(*exceptions.NotFound); ok {
			RespondAndLog(w, perr.Error(), TEXT, http.StatusNotFound, logger)
			return
		}
		RespondAndLog(w, err.Error(), TEXT, http.StatusInternalServerError, logger)
		return
	}
	asJson, err := json.Marshal(apps)
	if err != nil {
		RespondAndLog(w, err.Error(), TEXT, http.StatusInternalServerError, logger)
		return
	}
	RespondAndDontLog(w, string(asJson), JSON, http.StatusOK, logger)
}
