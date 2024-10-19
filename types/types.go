package types

import (
	"time"
)

type RegisterUserPayload struct {
	Email     string `json:"email" validate:"required,email"`
	FirstName string `json:"firstName" validate:"required"`
	Lastname  string `json:"lastName" validate:"required"`
	Password  string `json:"password" validate:"required,min=3,max=120"`
	Country   string `json:"country"`
	State     string `json:"state"`
}

type LoginUserPayload struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}
type User struct {
	ID        int    `json:"id"`
	Email     string `json:"email"`
	FirstName string `json:"firstName"`
	Lastname  string `json:"lastName"`
	Password  string `json:"password,omitempty"`
	Country   string `json:"country"`
	State     string `json:"state"`
	Role      string `json:"role,omitempty"`
	CreatedAt string `json:"createdAt"`
}

type Product struct {
	ID          int       `json:"id" validate:"required"`
	Name        string    `json:"name" validate:"required"`
	Description string    `json:"description"`
	Price       float64   `json:"price" validate:"required,gte=0"`
	Image       string    `json:"image_url"`
	IsActive    bool      `json:"isActive" validate:"required"`
	CreatedAt   time.Time `json:"createdAt"`
}

type ProductInventory struct {
	ID        int       `json:"id" validate:"required"`
	ProductID int       `json:"product_id" validate:"required"`
	Quantity  float64   `json:"quantity_available" validate:"required,gte=0"`
	Stock     int64     `json:"stock" validate:"required,gte=0"`
	CreatedAt time.Time `json:"createdAt"`
}

type ProductAndInventory struct {
	ID                int       `json:"id"`
	Name              string    `json:"name"`
	Description       string    `json:"description"`
	Price             float64   `json:"price"`
	ImageURL          string    `json:"image_url"`
	QuantityAvailable float64   `json:"quantity_available"`
	Stock             int64     `json:"stock"`
	CreatedAt         time.Time `json:"createdAt"`
}

type CreateProductPayload struct {
	Name        string    `json:"name" validate:"required"`
	Description string    `json:"description"`
	Price       float64   `json:"price" validate:"required,gte=0"`
	Image       string    `json:"image_url"`
	Quantity    float64   `json:"quantity_available" validate:"required,gte=0"`
	Stock       int64     `json:"stock" validate:"required,gte=0"`
	CreatedAt   time.Time `json:"createdAt"`
}
type ProductAndInventoryUpdate struct {
	ProductID   int      `json:"id" validate:"required"`
	Name        *string  `json:"name,omitempty"`
	Description *string  `json:"description,omitempty"`
	Price       *float64 `json:"price,omitempty"`
	Image       *string  `json:"image_url,omitempty"`
	IsActive    *bool    `json:"isActive,omitempty"`
	Quantity    *float64 `json:"quantity_available,omitempty"`
	Stock       *int64   `json:"stock,omitempty"`
}

type Role struct {
	ID   int64  `json:"id"`
	Name string `json:"name" validate:"required"`
}

type ProductRepository interface {
	CreateProduct(product CreateProductPayload) error
	GetProductByID(id int) (*ProductAndInventory, error)
	UpdateProduct(product ProductAndInventoryUpdate) error
	// DeleteProduct(id int) error
	ListProducts(offset, limit int) ([]ProductAndInventory, error)
}

// make this an interface for easy testing
type UserStore interface {
	GetUserByEmail(email string) (*User, error)
	GetUserByID(id int) (*User, error)
	CreateUser(user RegisterUserPayload) error
	// UpdateUser(user User) error
	// DeleteUser(id int) error
	ListUsers(offset, limit int) ([]User, error)
}

type RoleStore interface {
	CreateRole(role Role) error
	GetAllRoles() ([]Role, error)
}
