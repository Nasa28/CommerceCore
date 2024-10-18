package role

import (
	"database/sql"

	"github.com/Nasa28/CommerceCore/types"
)

type RoleStore struct {
	db *sql.DB
}

func NewRolesStore(db *sql.DB) *RoleStore {
	return &RoleStore{db: db}
}

func (r *RoleStore) CreateRole(role types.Role) error {

	_, err := r.db.Exec("INSERT INTO roles (name) VALUES (?)", role.Name)

	if err != nil {
		return err
	}
	return nil
}

func (r *RoleStore) GetAllRoles() ([]types.Role, error) {

	rows, err := r.db.Query("SELECT * FROM roles")

	if err != nil {
		return nil, err
	}
	roles := []types.Role{}

	for rows.Next() {
		var role types.Role
		err := rows.Scan(&role.ID, &role.Name)
		if err != nil {
			return nil, err
		}
		roles = append(roles, role)

	}
	defer rows.Close()

	return roles, nil
}
