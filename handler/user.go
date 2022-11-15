package handler

import (
	"encoding/json"
	"github.com/mg52/go-api/domain"
	"github.com/mg52/go-api/helper"
	"github.com/sirupsen/logrus"
	"net/http"
)

type userHandler struct {
	logger *logrus.Entry
}

func NewUserHandler(logger *logrus.Entry) Handler {
	return &userHandler{logger: logger}
}

func (h *userHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.logger.WithFields(logrus.Fields{
		"field_1": "abc",
	}).Info("test")

	w.Header().Set("content-type", "application/json")
	switch {
	case r.Method == http.MethodGet:
		h.List(w, r)
		return
	case r.Method == http.MethodPost:
		h.Create(w, r)
		return
	default:
		helper.NotFound(w, r)
		return
	}
}
func (h *userHandler) List(w http.ResponseWriter, r *http.Request) {
	var users []domain.User

	user1 := domain.User{
		ID:   "1",
		Name: "aaaaa",
	}
	user2 := domain.User{
		ID:   "2",
		Name: "bbbbb",
	}
	users = append(users, user1)
	users = append(users, user2)

	jsonBytes, err := json.Marshal(users)
	if err != nil {
		helper.InternalServerError(w, r)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(jsonBytes)
}
func (h *userHandler) Create(w http.ResponseWriter, r *http.Request) {
	var u domain.User
	if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
		helper.InternalServerError(w, r)
		return
	}

	jsonBytes, err := json.Marshal(u)
	if err != nil {
		helper.InternalServerError(w, r)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(jsonBytes)
}
