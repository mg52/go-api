package handler

import (
	"bytes"
	"encoding/json"
	"github.com/mg52/go-api/domain"
	"github.com/mg52/go-api/helper"
	"golang.org/x/crypto/bcrypt"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
)

type mockUserEntity struct {
}

func (entity *mockUserEntity) GetOneByUsername(username string) (*domain.User, error) {
	passBytes, _ := bcrypt.GenerateFromPassword([]byte("tafa"), 14)
	thePassword := string(passBytes)
	return &domain.User{
		ID:       10,
		Username: "mus",
		Password: thePassword,
	}, nil
}

func (entity *mockUserEntity) GetOneByUsernameAndPassword(username string, password string) (*domain.User, error) {
	passBytes, _ := bcrypt.GenerateFromPassword([]byte("tafa"), 14)
	thePassword := string(passBytes)
	return &domain.User{
		ID:       10,
		Username: "mus",
		Password: thePassword,
	}, nil
}

func (entity *mockUserEntity) CreateUser(user *domain.User) (int, error) {
	return 10, nil
}

func TestHealthCheckHandler(t *testing.T) {
	// Create a request to pass to our handler. We don't have any query parameters for now, so we'll
	// pass 'nil' as the third parameter.
	f := domain.Login{
		Username: "mus",
		Password: "tafa",
	}
	var buf bytes.Buffer
	err := json.NewEncoder(&buf).Encode(f)
	if err != nil {
		log.Fatal(err)
	}

	req, err := http.NewRequest("POST", "/login", &buf)
	if err != nil {
		t.Fatal(err)
	}

	userRepository := &mockUserEntity{}
	logrusEntry := helper.NewLogger("dev")
	authHandler := NewAuthHandler(logrusEntry, userRepository)

	// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(authHandler.ServeHTTP)

	// Our handlers satisfy http.Handler, so we can call their ServeHTTP method
	// directly and pass in our Request and ResponseRecorder.
	handler.ServeHTTP(rr, req)

	// Check the status code is what we expect.
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	actualResp := rr.Body.String()
	// Check the response body is what we expect.
	expected := `{"token":"ey`
	if actualResp[0:12] != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}
