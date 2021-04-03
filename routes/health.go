package routes

import (
	"encoding/json"
	"github.com/leonj1/compass/models"
	"github.com/rs/zerolog"
	"net/http"
	"os"
)

func (a *App) Health(w http.ResponseWriter, r *http.Request) {
	logger := zerolog.New(os.Stdout).With().Logger()
	asJson, _ := json.Marshal(models.HealthResponse{Version: a.Version})
	RespondAndDontLog(w, string(asJson), JSON, http.StatusOK, logger)
}
