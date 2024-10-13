package user

import (
	"database/sql"
	"fmt"

	"github.com/Nasa28/CommerceCore/types"
)

type Store struct {
	db *sql.DB
}

func NewStore(db *sql.DB) *Store {
	return &Store{db: db}
}

func (s *Store) GetUserByEmail(email string) (*types.User, error) {
	rows, err := s.db.Query("SELECT * FROM users WHERE email = ?", email)
	if err != nil {
		return nil, err
	}
	user := new(types.User)
	for rows.Next() {
		user, err = scanUsersIntoRows(rows)
		if err != nil {
			return nil, err
		}
	}
	if user.ID == 0 {
		return nil, fmt.Errorf("user not found")
	}
	return user, nil
}

func scanUsersIntoRows(rows *sql.Rows) (*types.User, error) {
	user := new(types.User)
	err := rows.Scan(
		&user.ID,
		&user.FirstName,
		&user.Lastname,
		&user.Email,
		&user.Password,
		&user.Country,
		&user.State,
		&user.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (s *Store) GetUserByID(id int) (*types.User, error) {
	rows, err := s.db.Query("SELECT * FROM users WHERE id = ?", id)
	if err != nil {
		return nil, err
	}
	user := new(types.User)
	for rows.Next() {
		user, err = scanUsersIntoRows(rows)
		if err != nil {
			return nil, err
		}
	}
	if user.ID == 0 {
		return nil, fmt.Errorf("user not found")
	}
	return user, nil
}

func (s *Store) CreateUser(user types.User) error {
	res, err := s.db.Exec("INSERT INTO users (firstName, lastName, email, state, country, password) VALUES (?,?,?,?,?,?)", user.FirstName, user.Lastname, user.Email, user.State, user.Country, user.Password)
	if err != nil {
		return err
	}
	rowsAffected, _ := res.RowsAffected()
	fmt.Printf("Rows affected: %d\n", rowsAffected)
	return nil
}
