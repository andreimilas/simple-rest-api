package user

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"

	"sample-rest-api/app/api"
)

// userAPI container - holds dependencies for the user API
type userAPI struct {
	handler *api.Handler
	store   *userStore
}

// AddRoutes - defines routes for the user resource
func AddRoutes(router *mux.Router, apiHandler *api.Handler) {
	// Initialize userAPI handler
	uAPI := &userAPI{
		apiHandler,
		&userStore{apiHandler.DB},
	}
	router.HandleFunc("/users", uAPI.listUsers).Methods("GET")
	router.HandleFunc("/users", uAPI.createUser).Methods("POST")
	router.HandleFunc("/users/{id}", uAPI.getUser).Methods("GET")
	router.HandleFunc("/users/{id}", uAPI.deleteUser).Methods("DELETE")
}

func (uAPI *userAPI) listUsers(w http.ResponseWriter, r *http.Request) {
	// Pagination parameters
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	offset, _ := strconv.Atoi(r.URL.Query().Get("offset"))
	if limit < 1 || limit > 25 {
		limit = 10
	}
	if offset < 0 {
		offset = 0
	}
	// List users
	users, err := uAPI.store.List(limit, offset)
	if err != nil {
		api.SendError(w, http.StatusInternalServerError, err.Error())
		return
	}

	// Send the JSON response
	api.SendJSONResponse(w, http.StatusOK, users)
}

func (uAPI *userAPI) createUser(w http.ResponseWriter, r *http.Request) {
	user := &User{}
	// Try to decode the request body into the user instance
	err := json.NewDecoder(r.Body).Decode(user)
	if err != nil {
		api.SendError(w, http.StatusBadRequest, "")
		return
	}

	// Create user
	err = uAPI.store.Create(user)
	if err != nil {
		api.SendError(w, http.StatusInternalServerError, err.Error())
		return
	}

	// Send the JSON response
	api.SendJSONResponse(w, http.StatusCreated, nil)
}

func (uAPI *userAPI) getUser(w http.ResponseWriter, r *http.Request) {
	// Get path parameters
	params := mux.Vars(r)
	userID := params["id"]

	// Get user
	user, err := uAPI.store.Get(userID)
	if err != nil {
		// If the entry does not exist, return 404
		if err == sql.ErrNoRows {
			api.SendError(w, http.StatusNotFound, "")
			return
		}

		api.SendError(w, http.StatusInternalServerError, err.Error())
		return
	}

	// Send the JSON response
	api.SendJSONResponse(w, http.StatusOK, user)
}

func (uAPI *userAPI) deleteUser(w http.ResponseWriter, r *http.Request) {
	// Get path parameters
	params := mux.Vars(r)
	userID := params["id"]

	// Delete user
	err := uAPI.store.Delete(userID)
	if err != nil {
		api.SendError(w, http.StatusInternalServerError, err.Error())
		return
	}

	// Send the JSON response
	api.SendJSONResponse(w, http.StatusNoContent, nil)
}
