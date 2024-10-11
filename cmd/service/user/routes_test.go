package user

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Nasa28/CommerceCore/types"
	"github.com/gorilla/mux"
)

func TestUserServicehandlers(t *testing.T) {
	userStore := &MockUserStore{}
	handler := NewHandler(userStore)

	t.Run("It should fail if the payload is invalid", func(t *testing.T) {
		payload := types.RegisterUserPayload{
			Email:       "johndoeemail.com",
			FirstName:   "firstName",
			Lastname:    "lastName",
			Password:    "password",
			Country:     "country",
			State:       "state",
			PhoneNumber: "phoneNumber",
		}
		marshalled, _ := json.Marshal(payload)
		req, err := http.NewRequest(http.MethodPost, "/register", bytes.NewBuffer(marshalled))

		if err != nil {
			log.Fatal(err)
		}

		rr := httptest.NewRecorder()

		router := mux.NewRouter()
		router.HandleFunc("/register", handler.handleRegister)
		router.ServeHTTP(rr, req)

		if rr.Code != http.StatusBadRequest {
			t.Errorf("Expected status code %d , got %d", http.StatusBadRequest, rr.Code)
		}
	})

	t.Run("Should successfully register the user", func(t *testing.T) {
		payload := types.RegisterUserPayload{
			Email:       "johndoe@email.com",
			FirstName:   "firstName",
			Lastname:    "lastName",
			Password:    "password",
			Country:     "country",
			State:       "state",
			PhoneNumber: "phoneNumber",
		}
		marshalled, _ := json.Marshal(payload)
		req, err := http.NewRequest(http.MethodPost, "/register", bytes.NewBuffer(marshalled))

		if err != nil {
			log.Fatal(err)
		}

		rr := httptest.NewRecorder()

		router := mux.NewRouter()
		router.HandleFunc("/register", handler.handleRegister)
		router.ServeHTTP(rr, req)

		if rr.Code != http.StatusCreated {
			t.Errorf("Expected status code %d , got %d", http.StatusCreated, rr.Code)
		}
	})

}

type MockUserStore struct{}

func (m *MockUserStore) GetUserByEmail(email string) (*types.User, error) {
	return nil, fmt.Errorf("User not found")
}

func (m *MockUserStore) GetUserByID(id int) (*types.User, error) {
	return nil, nil
}

func (m *MockUserStore) CreateUser(types.User) error {
	return nil
}
