package api

import (
	"net/http/httptest"
	"testing"

	"github.com/jmoiron/sqlx"
)

func TestInit(t *testing.T) {
	t.Run("Initialize API Handler", func(t *testing.T) {
		// First, init DB handle
		dbHandle := &sqlx.DB{}

		apiHandler := Init(dbHandle)
		if apiHandler == nil {
			t.Error("API Handler should be initialized.")
		}
	})
}

func TestSendJSONResponse(t *testing.T) {
	t.Run("Send JSON response", func(t *testing.T) {
		// Init response writer
		w := httptest.NewRecorder()

		SendJSONResponse(w, 200, "")
		if w.Header().Get("Content-Type") != "application/json" ||
			w.Code != 200 {
			t.Error("Invalid response.")
		}

	})
}

func TestSendError(t *testing.T) {
	t.Run("Send Error", func(t *testing.T) {
		// Init response writer
		w := httptest.NewRecorder()

		SendError(w, 500, "")
		if w.Code != 500 {
			t.Error("Invalid response.")
		}

	})
}
