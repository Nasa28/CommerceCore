package types

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
	Password  string `json:"password"`
	Country   string `json:"country"`
	State     string `json:"state"`
	CreatedAt string `json:"createdAt"`
}

// make this an interface for easy testing
type UserStore interface {
	GetUserByEmail(email string) (*User, error)
	GetUserByID(id int) (*User, error)
	CreateUser(User) error
}
