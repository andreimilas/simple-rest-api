package user

import (
	"log"

	"github.com/jmoiron/sqlx"
)

type userStore struct {
	DB *sqlx.DB
}

// List - store method for listing users
func (ss *userStore) List(limit int, offset int) ([]User, error) {
	users := make([]User, 0)
	userQuery := `SELECT id, uuid, first_name, last_name, email, is_active, created, modified FROM user LIMIT ? OFFSET ?`
	// Execute the query while preventing SQL injection
	err := ss.DB.Select(&users, userQuery, limit, offset)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	return users, nil
}

// Create - store method for creating a user
func (ss *userStore) Create(user *User) error {
	userQuery := `INSERT INTO user (uuid, first_name, last_name, email, is_active, created, modified) 
				VALUES (UUID(), ?, ?, ?, ?, NOW(), NOW())`
	// Execute the query while preventing SQL injection
	_, err := ss.DB.Exec(userQuery, user.FirstName, user.LastName, user.Email, user.IsActive)
	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}

// Get - store method for fetching a user
func (ss *userStore) Get(userID string) (*User, error) {
	user := &User{}
	userQuery := `SELECT id, uuid, first_name, last_name, email, is_active, created, modified FROM user WHERE uuid = ? LIMIT 1`
	// Execute the query while preventing SQL injection
	err := ss.DB.Get(user, userQuery, userID)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	return user, nil
}

// Delete - store method for deleting a user
func (ss *userStore) Delete(userID string) error {
	userQuery := `DELETE FROM user WHERE uuid = ?`
	// Execute the query while preventing SQL injection
	_, err := ss.DB.Exec(userQuery, userID)
	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}
