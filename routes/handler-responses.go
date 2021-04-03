package routes

import (
	"github.com/rs/zerolog"
	"net/http"
)

func RespondAndLog(w http.ResponseWriter, msg, contentType string, statusCode int, logger zerolog.Logger) {
	Respond(w, msg, contentType, statusCode, true, logger)
}

func RespondAndDontLog(w http.ResponseWriter, msg, contentType string, statusCode int, logger zerolog.Logger) {
	Respond(w, msg, contentType, statusCode, false, logger)
}

func Respond(w http.ResponseWriter, msg, contentType string, statusCode int, writeLog bool, logger zerolog.Logger) {
	if writeLog {
		logger.Debug().Msg(msg)
	}
	w.Header().Set(ContentType, contentType)
	w.WriteHeader(statusCode)
	_, _ = w.Write([]byte(msg))
	return
}
