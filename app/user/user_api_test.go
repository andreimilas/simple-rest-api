package user

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"

	"sample-rest-api/app/api"
)

func TestAPIAddRoutes(t *testing.T) {
	t.Run("Add User routes", func(t *testing.T) {
		router := mux.NewRouter().StrictSlash(true)
		apiHandler := api.Init(&sqlx.DB{})

		AddRoutes(router, apiHandler)
		// Iterate over the registered routes
		exists := false
		router.Walk(func(route *mux.Route, router *mux.Router, ancestors []*mux.Route) error {
			path, _ := route.GetPathTemplate()
			if path == "/users" {
				exists = true
			}
			return nil
		})

		if !exists {
			t.Error("User routes not registered.")
		}
	})
}

func TestAPIListUsers(t *testing.T) {
	t.Run("API List users", func(t *testing.T) {
		// Create a mock sql db connection
		db, mock, err := sqlmock.New()
		if err != nil {
			t.Error("Error while opening mock SQL connection.")
		}
		defer db.Close()

		// Add rows to the database
		rows := sqlmock.NewRows([]string{"id", "uuid", "first_name", "last_name", "email", "is_active", "created", "modified"}).
			AddRow(1, "1e7aceca-9da3-11ea-bd4c-0242ac140002", "User1FirstName", "User1LastName", "u1fn.u1ln@mail.test", true, time.Now(), time.Now()).
			AddRow(2, "1e7ad3d8-9da3-11ea-bd4c-0242ac140002", "User2FirstName", "User2LastName", "u2fn.u2ln@mail.test", false, time.Now(), time.Now()).
			AddRow(3, "1e7ad456-9da3-11ea-bd4c-0242ac140002", "User3FirstName", "User3LastName", "u3fn.u3ln@mail.test", true, time.Now(), time.Now())
		mock.ExpectQuery("^SELECT id, uuid, first_name, last_name, email, is_active, created, modified FROM user LIMIT \\? OFFSET \\?").
			WithArgs(10, 0).
			WillReturnRows(rows)

		dbHandle := sqlx.NewDb(db, "mysql")
		// Initialize API and router
		apiHandler := api.Init(dbHandle)
		router := mux.NewRouter().StrictSlash(true)
		AddRoutes(router, apiHandler)

		// Send request
		req, _ := http.NewRequest("GET", "/users", nil)
		response := httptest.NewRecorder()
		router.ServeHTTP(response, req)

		// Check response code
		if response.Code != 200 {
			t.Error("Incorrect response code.")
			return
		}

		// Check response body
		userList := make([]User, 0)
		err = json.Unmarshal(response.Body.Bytes(), &userList)
		if err != nil {
			t.Error("Invalid JSON in response body.")
			return
		}

		// Check expected response
		if len(userList) != 3 {
			t.Error("Incorrect response body.")
			return
		}
	})
}

func TestAPICreateUser(t *testing.T) {
	t.Run("API Create user", func(t *testing.T) {
		// Create a mock sql db connection
		db, mock, err := sqlmock.New()
		if err != nil {
			t.Error("Error while opening mock SQL connection.")
		}
		defer db.Close()

		mock.ExpectExec("^INSERT INTO user \\(uuid, first_name, last_name, email, is_active, created, modified\\) VALUES \\(UUID\\(\\), \\?, \\?, \\?, \\?, NOW\\(\\), NOW\\(\\)\\)").
			WithArgs("User1FirstName", "User1LastName", "u1fn.u1ln@mail.test", true).
			WillReturnResult(sqlmock.NewResult(1, 1))

		dbHandle := sqlx.NewDb(db, "mysql")
		// Initialize API and router
		apiHandler := api.Init(dbHandle)
		router := mux.NewRouter().StrictSlash(true)
		AddRoutes(router, apiHandler)

		jsonUser, _ := json.Marshal(User{FirstName: "User1FirstName", LastName: "User1LastName", Email: "u1fn.u1ln@mail.test", IsActive: true})
		// Send request
		req, _ := http.NewRequest("POST", "/users", bytes.NewBuffer(jsonUser))
		response := httptest.NewRecorder()
		router.ServeHTTP(response, req)

		// Check response code
		if response.Code != 201 {
			t.Error("Incorrect response code.")
			return
		}
	})
}

func TestAPIGetUser(t *testing.T) {
	t.Run("API Get user", func(t *testing.T) {
		// Create a mock sql db connection
		db, mock, err := sqlmock.New()
		if err != nil {
			t.Error("Error while opening mock SQL connection.")
		}
		defer db.Close()

		// Add rows to the database
		rows := sqlmock.NewRows([]string{"id", "uuid", "first_name", "last_name", "email", "is_active", "created", "modified"}).
			AddRow(1, "1e7aceca-9da3-11ea-bd4c-0242ac140002", "User1FirstName", "User1LastName", "u1fn.u1ln@mail.test", true, time.Now(), time.Now())
		mock.ExpectQuery("^SELECT id, uuid, first_name, last_name, email, is_active, created, modified FROM user WHERE uuid = \\?").
			WithArgs("1e7aceca-9da3-11ea-bd4c-0242ac140002").
			WillReturnRows(rows)

		dbHandle := sqlx.NewDb(db, "mysql")
		// Initialize API and router
		apiHandler := api.Init(dbHandle)
		router := mux.NewRouter().StrictSlash(true)
		AddRoutes(router, apiHandler)

		// Send request
		req, _ := http.NewRequest("GET", "/users/1e7aceca-9da3-11ea-bd4c-0242ac140002", nil)
		response := httptest.NewRecorder()
		router.ServeHTTP(response, req)

		// Check response code
		if response.Code != 200 {
			t.Error("Incorrect response code.")
			return
		}

		// Check response body
		user := &User{}
		err = json.Unmarshal(response.Body.Bytes(), user)
		if err != nil {
			t.Error("Invalid JSON in response body.")
			return
		}

		// Check expected response
		if user.UUID.String() != "1e7aceca-9da3-11ea-bd4c-0242ac140002" {
			t.Error("Incorrect response body.")
			return
		}
	})
}

func TestAPIGetNonExistingUser(t *testing.T) {
	t.Run("API Get non-existing user", func(t *testing.T) {
		// Create a mock sql db connection
		db, mock, err := sqlmock.New()
		if err != nil {
			t.Error("Error while opening mock SQL connection.")
		}
		defer db.Close()

		// Add rows to the database
		rows := sqlmock.NewRows([]string{"id", "uuid", "first_name", "last_name", "email", "is_active", "created", "modified"})
		mock.ExpectQuery("^SELECT id, uuid, first_name, last_name, email, is_active, created, modified FROM user WHERE uuid = \\?").
			WithArgs("1e7aceca-9da3-11ea-bd4c-0242ac140002").
			WillReturnRows(rows)

		dbHandle := sqlx.NewDb(db, "mysql")
		// Initialize API and router
		apiHandler := api.Init(dbHandle)
		router := mux.NewRouter().StrictSlash(true)
		AddRoutes(router, apiHandler)

		// Send request
		req, _ := http.NewRequest("GET", "/users/1e7aceca-9da3-11ea-bd4c-0242ac140002", nil)
		response := httptest.NewRecorder()
		router.ServeHTTP(response, req)

		// Check response code
		if response.Code != 404 {
			t.Error("Incorrect response code.")
			return
		}
	})
}

func TestAPIDeleteUser(t *testing.T) {
	t.Run("API Delete user", func(t *testing.T) {
		// Create a mock sql db connection
		db, mock, err := sqlmock.New()
		if err != nil {
			t.Error("Error while opening mock SQL connection.")
		}
		defer db.Close()

		mock.ExpectExec("^DELETE FROM user WHERE uuid = \\?").
			WithArgs("1e7aceca-9da3-11ea-bd4c-0242ac140002").
			WillReturnResult(sqlmock.NewResult(1, 1))

		dbHandle := sqlx.NewDb(db, "mysql")
		// Initialize API and router
		apiHandler := api.Init(dbHandle)
		router := mux.NewRouter().StrictSlash(true)
		AddRoutes(router, apiHandler)

		// Send request
		req, _ := http.NewRequest("DELETE", "/users/1e7aceca-9da3-11ea-bd4c-0242ac140002", nil)
		response := httptest.NewRecorder()
		router.ServeHTTP(response, req)

		// Check response code
		if response.Code != 204 {
			t.Error("Incorrect response code.")
			return
		}
	})
}
