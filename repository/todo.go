package repository

import (
	"database/sql"
	_ "github.com/lib/pq"
	"github.com/mg52/go-api/domain"
)

type ITodo interface {
	GetAll(uId int) ([]domain.Todo, error)
	CreateTodo(uId int, todo *domain.Todo) (int, error)
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
		var td domain.Todo
		if err := rows.Scan(&td.ID, &td.UserID, &td.Content); err != nil {
			return todos, err
		}
		todos = append(todos, td)
	}
	if err = rows.Err(); err != nil {
		return todos, err
	}
	return todos, nil
}

func (entity *todoEntity) CreateTodo(uId int, todo *domain.Todo) (int, error) {
	return 0, nil
}

//
//func (entity *userEntity) GetOneByUsername(username string) (*domain.User, error) {
//	sqlStatementSelect := `SELECT * FROM users WHERE username=$1;`
//	var user domain.User
//	row := entity.db.QueryRow(sqlStatementSelect, username)
//	errSelect := row.Scan(&user.ID, &user.Username, &user.Password)
//	if errSelect != nil && errSelect == sql.ErrNoRows {
//		return nil, errors.New("user not found")
//	} else {
//		return &user, nil
//	}
//}
//
//func (entity *userEntity) GetOneByUsernameAndPassword(username string, password string) (*domain.User, error) {
//	sqlStatementSelect := `SELECT * FROM users WHERE username=$1 and password=$2;`
//	var user domain.User
//	row := entity.db.QueryRow(sqlStatementSelect, username, password)
//	errSelect := row.Scan(&user.ID, &user.Username, &user.Password)
//	if errSelect != nil && errSelect == sql.ErrNoRows {
//		return nil, errors.New("user not found")
//	} else {
//		return &user, nil
//	}
//}
//
//func (entity *userEntity) CreateUser(user *domain.User) (int, error) {
//	insertStatement := `
//	INSERT INTO users (username, password)
//	VALUES ($1, $2)
//	RETURNING id`
//	id := 0
//	err := entity.db.QueryRow(insertStatement, user.Username, user.Password).Scan(&id)
//	if err != nil {
//		return -1, errors.New("user cannot be created")
//	}
//	return id, nil
//}
