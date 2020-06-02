package user

import (
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
)

func TestStoreList(t *testing.T) {
	t.Run("List users - Limit 3; Offset 1", func(t *testing.T) {

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
			WithArgs(3, 1).
			WillReturnRows(rows)

		dbHandle := sqlx.NewDb(db, "mysql")
		// Initialize user store
		userStore := &userStore{dbHandle}
		userList, err := userStore.List(3, 1)
		if err != nil {
			t.Error("Unexpected error.")
		}

		// Check user count
		if len(userList) != 3 {
			t.Error("User count should be 3")
		}
	})
}

func TestStoreCreate(t *testing.T) {
	t.Run("Create user", func(t *testing.T) {

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
		// Initialize user store
		userStore := &userStore{dbHandle}
		// Build user instance
		user := &User{FirstName: "User1FirstName", LastName: "User1LastName", Email: "u1fn.u1ln@mail.test", IsActive: true}
		err = userStore.Create(user)
		if err != nil {
			t.Error("Unexpected error.")
		}
	})
}

func TestStoreGet(t *testing.T) {
	t.Run("Get user", func(t *testing.T) {

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
		// Initialize user store
		userStore := &userStore{dbHandle}
		user, err := userStore.Get("1e7aceca-9da3-11ea-bd4c-0242ac140002")
		if err != nil {
			t.Error("Unexpected error.")
		}

		// Check user count
		if user.UUID.String() != "1e7aceca-9da3-11ea-bd4c-0242ac140002" {
			t.Error("Invalid user")
		}
	})
}

func TestStoreDelete(t *testing.T) {
	t.Run("Delete user", func(t *testing.T) {

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
		// Initialize user store
		userStore := &userStore{dbHandle}
		err = userStore.Delete("1e7aceca-9da3-11ea-bd4c-0242ac140002")
		if err != nil {
			t.Error("Unexpected error.")
		}
	})
}
