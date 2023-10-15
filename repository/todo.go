package repository

import (
	"database/sql"
	"errors"

	_ "github.com/lib/pq"
	"github.com/mg52/go-api/domain"
)

type ITodo interface {
	GetAll(uId int) ([]domain.Todo, error)
	CreateTodo(todo *domain.Todo) (int, error)
}

var TodoEntity ITodo

type todoEntity struct {
	db *sql.DB
}

func NewTodoEntity(db *sql.DB) ITodo {
	TodoEntity = &todoEntity{db: db}
	return TodoEntity
}

func (entity *todoEntity) GetAll(uId int) ([]domain.Todo, error) {
	sqlStatementSelect := `SELECT * FROM todos WHERE user_id=$1;`
	rows, err := entity.db.Query(sqlStatementSelect, uId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var todos []domain.Todo
	for rows.Next() {
		var rowId int
		var td domain.Todo
		if err := rows.Scan(&rowId, &td.UserID, &td.Content); err != nil {
			return todos, err
		}
		todos = append(todos, td)
	}
	if err = rows.Err(); err != nil {
		return todos, err
	}
	return todos, nil
}

func (entity *todoEntity) CreateTodo(todo *domain.Todo) (int, error) {
	insertStatement := `INSERT INTO todos (user_id, content) VALUES ($1, $2) RETURNING id`
	id := 0
	err := entity.db.QueryRow(insertStatement, todo.UserID, todo.Content).Scan(&id)
	if err != nil {
		return -1, errors.New("todo cannot be created")
	}
	return id, nil
}
