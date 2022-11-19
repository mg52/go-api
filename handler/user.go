package handler

import (
	"encoding/json"
	_ "github.com/mg52/go-api/docs"
	"github.com/mg52/go-api/domain"
	"github.com/mg52/go-api/helper"
	"github.com/mg52/go-api/repository"
	"github.com/sirupsen/logrus"
	"net/http"
)

type userHandler struct {
	logger         *logrus.Entry
	userRepository repository.IUser
}

func NewUserHandler(logger *logrus.Entry, userRepository repository.IUser) Handler {
	return &userHandler{logger: logger, userRepository: userRepository}
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

// user godoc
// @Summary     User List
// @Description User List
// @Tags        user
// @ID          auth-login
// @Accept      json
// @Produce     json
// @Param Authorization header string true "Token with the bearer started"
// @Success     200 {array} domain.User
// @Router      /user [get]
func (h *userHandler) List(w http.ResponseWriter, r *http.Request) {
	allUsers, err := h.userRepository.GetAll()
	if err != nil {
		helper.InternalServerError(w, r, err)
		return
	}

	jsonBytes, err := json.Marshal(allUsers)
	if err != nil {
		helper.InternalServerError(w, r, err)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(jsonBytes)
}

// user godoc
// @Summary     User Create
// @Description User Create
// @Tags        user
// @Accept      json
// @Produce     json
// @Param Authorization header string true "Token with the bearer started"
// @Param       authLogin body     domain.User true "User Input"
// @Success     200       {object} domain.User
// @Router      /user [post]
func (h *userHandler) Create(w http.ResponseWriter, r *http.Request) {
	var u domain.User
	if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
		helper.InternalServerError(w, r, err)
		return
	}

	jsonBytes, err := json.Marshal(u)
	if err != nil {
		helper.InternalServerError(w, r, err)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(jsonBytes)
}
