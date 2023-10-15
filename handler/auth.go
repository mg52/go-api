package handler

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/go-playground/validator/v10"
	_ "github.com/mg52/go-api/docs"
	"github.com/mg52/go-api/domain"
	"github.com/mg52/go-api/helper"
	"github.com/mg52/go-api/repository"
	"golang.org/x/crypto/bcrypt"
)

type authHandler struct {
	userRepository repository.IUser
}

func NewAuthHandler(userRepository repository.IUser) Handler {
	return &authHandler{userRepository: userRepository}
}

func (h *authHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "application/json")
	switch {
	case r.Method == http.MethodPost:
		h.Login(w, r)
		return
	case r.Method == http.MethodPut:
		h.Signup(w, r)
		return
	default:
		helper.NotFound(w, r)
		return
	}
}

// Auth godoc
// @Summary     Login
// @Description Login
// @Tags        Auth
// @Accept      json
// @Produce     json
// @Param       auth body     domain.Login true "Auth Input"
// @Success     200       {object} domain.Authentication
// @Router      /auth [post]
func (h *authHandler) Login(w http.ResponseWriter, r *http.Request) {
	var u domain.Login
	if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
		helper.InternalServerError(w, r, err)
		return
	}

	validate := validator.New()
	err := validate.Struct(u)
	if err != nil {
		helper.Resp(w, r, http.StatusBadRequest, err.Error())
		return
	}

	theUser, err := h.userRepository.GetOneByUsername(u.Username)
	if err != nil || theUser == nil {
		helper.Unauthorized(w, r, errors.New("cannot find username"))
		return
	}

	match := helper.CheckPasswordHash(u.Password, theUser.Password)
	if !match {
		helper.Unauthorized(w, r, errors.New("password is incorrect"))
		return
	}

	token, err := helper.GenerateJWTToken(*theUser)
	if err != nil {
		helper.InternalServerError(w, r, err)
		return
	}
	auth := domain.Authentication{Token: token}

	jsonBytes, err := json.Marshal(auth)
	if err != nil {
		helper.InternalServerError(w, r, err)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(jsonBytes)
}

// Auth godoc
// @Summary     Signup
// @Description Signup
// @Tags        Auth
// @Accept      json
// @Produce     json
// @Param       signup body     domain.Login true "Sign up Input"
// @Success     200       {object} domain.Authentication
// @Router      /auth [put]
func (h *authHandler) Signup(w http.ResponseWriter, r *http.Request) {
	var u domain.Login
	if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
		helper.InternalServerError(w, r, err)
		return
	}

	validate := validator.New()
	err := validate.Struct(u)
	if err != nil {
		helper.Resp(w, r, http.StatusBadRequest, err.Error())
		return
	}

	theUser, err := h.userRepository.GetOneByUsername(u.Username)
	if theUser != nil {
		helper.Resp(w, r, http.StatusBadRequest, "username is already taken.")
		return
	}
	aUser := domain.User{}
	aUser.Username = u.Username
	passBytes, _ := bcrypt.GenerateFromPassword([]byte(u.Password), 14)
	aUser.Password = string(passBytes)

	id, err := h.userRepository.CreateUser(&aUser)
	if err != nil {
		helper.InternalServerError(w, r, err)
		return
	}

	aUser.ID = id

	token, err := helper.GenerateJWTToken(aUser)
	if err != nil {
		helper.InternalServerError(w, r, err)
		return
	}
	auth := domain.Authentication{Token: token}

	jsonBytes, err := json.Marshal(auth)
	if err != nil {
		helper.InternalServerError(w, r, err)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(jsonBytes)
}
