package handler

import (
	"encoding/json"
	_ "github.com/mg52/go-api/docs"
	"github.com/mg52/go-api/helper"
	"github.com/mg52/go-api/repository"
	"github.com/sirupsen/logrus"
	"net/http"
	"strconv"
)

type todoHandler struct {
	logger         *logrus.Entry
	todoRepository repository.ITodo
}

func NewTodoHandler(logger *logrus.Entry, todoRepository repository.ITodo) Handler {
	return &todoHandler{logger: logger, todoRepository: todoRepository}
}

func (h *todoHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.logger.WithFields(logrus.Fields{
		"field_1": "abc",
	}).Info("test")

	w.Header().Set("content-type", "application/json")
	switch {
	//case r.Method == http.MethodGet:
	//	h.List(w, r)
	//	return
	case r.Method == http.MethodGet:
		h.List(w, r)
		return
	default:
		helper.NotFound(w, r)
		return
	}
}

// todo godoc
// @Summary     Todo List
// @Description Todo List
// @Tags        todo
// @ID          auth-login
// @Accept      json
// @Produce     json
// @Param Authorization header string true "Token with the bearer started"
// @Success     200 {array} domain.Todo
// @Router      /todo [get]
func (h *todoHandler) List(w http.ResponseWriter, r *http.Request) {
	uId, err := strconv.Atoi(r.Header.Get("UID"))
	if err != nil {
		helper.InternalServerError(w, r, err)
		return
	}

	allTodos, err := h.todoRepository.GetAll(uId)
	if err != nil {
		helper.InternalServerError(w, r, err)
		return
	}

	jsonBytes, err := json.Marshal(allTodos)
	if err != nil {
		helper.InternalServerError(w, r, err)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(jsonBytes)
}
