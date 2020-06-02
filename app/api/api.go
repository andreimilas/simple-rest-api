package api

import (
	"encoding/json"
	"net/http"

	"github.com/jmoiron/sqlx"
)

// Handler - Holds API specific dependencies
type Handler struct {
	DB *sqlx.DB
}

// Init - Initialize API
func Init(db *sqlx.DB) *Handler {
	return &Handler{
		DB: db,
	}
}

// SendJSONResponse -
func SendJSONResponse(w http.ResponseWriter, statusCode int, content interface{}) {
	// Try to marshal the content
	jsonContent, err := json.Marshal(content)
	if err != nil {
		// Marshalling error, send 500
		SendError(w, http.StatusInternalServerError, err.Error())
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	w.Write(jsonContent)
}

// SendError -
func SendError(w http.ResponseWriter, statusCode int, errorMessage string) {
	http.Error(w, errorMessage, statusCode)
}
