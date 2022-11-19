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

type loginHandler struct {
	logger         *logrus.Entry
	userRepository repository.IUser
}

func NewLoginHandler(logger *logrus.Entry, userRepository repository.IUser) Handler {
	return &loginHandler{logger: logger, userRepository: userRepository}
}

func (h *loginHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "application/json")
	switch {
	case r.Method == http.MethodPost:
		h.Login(w, r)
		return
	default:
		helper.NotFound(w, r)
		return
	}
}

// login godoc
// @Summary     Login
// @Description Login
// @Tags        login
// @Accept      json
// @Produce     json
// @Param       authLogin body     domain.Login true "Login Input"
// @Success     200       {object} domain.Authentication
// @Router      /login [post]
func (h *loginHandler) Login(w http.ResponseWriter, r *http.Request) {
	var u domain.Login
	if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
		helper.InternalServerError(w, r, err)
		return
	}

	theUser, err := h.userRepository.GetOneByUsernameAndPassword(u.Name, u.Password)
	if err != nil || theUser == nil {
		helper.Resp(w, r, http.StatusUnauthorized, "username or password is incorrect")
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
