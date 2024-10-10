package types

type RegisterUserPayload struct {
	Email       string `json:"email"`
	FirstName   string `json:"firstName"`
	Lastname    string `json:"lastName"`
	Password    string `json:"password"`
	Country     string `json:"country"`
	State       string `json:"state"`
	PhoneNumber string `json:"phoneNumber"`
}

type User struct {
	ID          int    `json:"id"`
	Email       string `json:"email"`
	FirstName   string `json:"firstName"`
	Lastname    string `json:"lastName"`
	Password    string `json:"password"`
	Country     string `json:"country"`
	State       string `json:"state"`
	PhoneNumber string `json:"phoneNumber"`
	CreatedAt   string `json:"createdAt"`
}

// make this an interface for easy testing
type UserStore interface {
	GetUserByEmail(email string) (*User, error)
	GetUserByID(id int) (*User, error)
	CreateUser(User) error
}
