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
	row := s.db.QueryRow("SELECT id, email, firstName, lastName, password, country, state FROM users WHERE email = ?", email)

	// Create a user object to hold the result
	user := new(types.User)

	// Scan the result into the user struct
	err := row.Scan(
		&user.ID,
		&user.Email,
		&user.FirstName,
		&user.Lastname,
		&user.Password,
		&user.Country,
		&user.State,
	)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("user not found")
	}

	if err != nil {
		return nil, err
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

func (s *Store) ListUsers(limit, offset int) ([]types.User, error) {
	query := `
		SELECT 
			u.id,
			u.email, 
			u.firstName, 
			u.lastName, 
			u.state,
			u.country,
			COALESCE(r.name, 'user') AS role, 
			u.createdAt 
		FROM users AS u 
		LEFT JOIN user_roles AS ur ON u.id = ur.user_id
		LEFT JOIN roles AS r ON ur.role_id = r.id
		LIMIT ? OFFSET ?
		`
	rows, err := s.db.Query(query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	users := []types.User{}

	for rows.Next() {
		var user types.User

		err := rows.Scan(
			&user.ID,
			&user.Email,
			&user.FirstName,
			&user.Lastname,
			&user.State,
			&user.Country,
			&user.Role,
			&user.CreatedAt,
		)
		if err != nil {
			return nil, err
		}

		users = append(users, user)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}
	return users, nil
}
