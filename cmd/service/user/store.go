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
		&user.Country,
		&user.State,
		&user.Role,
		&user.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (s *Store) GetUserByID(id int) (*types.User, error) {
	rows, err := s.db.Query(`
	SELECT 
		u.id, 
		u.email, 
		u.firstName, 
		u.lastName, 
		u.country, 
		u.state, 
		r.name AS role_name,
		u.createdAt 
	FROM users AS u 
	RIGHT JOIN user_roles AS ur ON ur.user_id = u.id
	LEFT JOIN roles AS r ON r.id = ur.role_id
	WHERE u.id = ?`, id)

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
	defer rows.Close()
	return user, nil
}

func (s *Store) CreateUser(user types.RegisterUserPayload) error {
	var defaultRole = "user"

	// Insert user into the `users` table
	res, err := s.db.Exec(
		"INSERT INTO users (firstName, lastName, email, state, country, password) VALUES (?,?,?,?,?,?)", 
		user.FirstName, user.Lastname, user.Email, user.State, user.Country, user.Password)
	if err != nil {
		return err
	}

	// Get the last inserted user ID
	lastInsertedId, err := res.LastInsertId()
	if err != nil {
		return err
	}

	// Insert the default role if it does not exist using INSERT IGNORE
	_, err = s.db.Exec("INSERT IGNORE INTO roles (name) VALUES (?)", defaultRole)
	if err != nil {
		return err
	}

	// Fetch the role ID of the default role (either existing or just inserted)
	role := new(types.Role)
	err = s.db.QueryRow("SELECT id FROM roles WHERE name = ?", defaultRole).Scan(&role.ID)
	if err != nil {
		return err
	}

	// Assign the default role to the user (assuming user_roles table exists)
	_, err = s.db.Exec("INSERT INTO user_roles (user_id, role_id) VALUES (?, ?)", lastInsertedId, role.ID)
	if err != nil {
		return err
	}

	return nil
}
