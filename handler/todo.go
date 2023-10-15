package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-playground/validator/v10"
	_ "github.com/mg52/go-api/docs"
	"github.com/mg52/go-api/domain"
	"github.com/mg52/go-api/helper"
	"github.com/mg52/go-api/repository"
)

type todoHandler struct {
	todoRepository repository.ITodo
}

func NewTodoHandler(todoRepository repository.ITodo) Handler {
	return &todoHandler{todoRepository: todoRepository}
}

func (h *todoHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "application/json")
	switch {
	case r.Method == http.MethodPut:
		h.Create(w, r)
		return
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
// @Tags        Todo
// @ID          auth-login
// @Accept      json
// @Produce     json
// @Param Authorization header string true "Token with the Bearer started"
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

// Auth godoc
// @Summary     Todo Create
// @Description Todo Create
// @Tags        Todo
// @ID          auth-login
// @Accept      json
// @Produce     json
// @Param Authorization header string true "Token with the Bearer started"
// @Param       auth body     domain.TodoRequest true "Todo Input"
// @Success     200 {array} domain.Todo
// @Router      /todo [put]
func (h *todoHandler) Create(w http.ResponseWriter, r *http.Request) {
	var u domain.Todo
	if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
		helper.InternalServerError(w, r, err)
		return
	}

	uId, err := strconv.Atoi(r.Header.Get("UID"))
	if err != nil {
		helper.InternalServerError(w, r, err)
		return
	}

	u.UserID = uId

	validate := validator.New()
	err = validate.Struct(u)
	if err != nil {
		helper.Resp(w, r, http.StatusBadRequest, err.Error())
		return
	}

	_, err = h.todoRepository.CreateTodo(&u)
	if err != nil {
		helper.InternalServerError(w, r, err)
		return
	}

	returnTodo := domain.Todo{
		UserID:  uId,
		Content: u.Content,
	}

	jsonBytes, err := json.Marshal(returnTodo)
	if err != nil {
		helper.InternalServerError(w, r, err)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(jsonBytes)
}
