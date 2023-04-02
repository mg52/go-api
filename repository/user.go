package repository

import (
	"database/sql"
	"errors"
	_ "github.com/lib/pq"
	"github.com/mg52/go-api/domain"
)

type IUser interface {
	//GetAll() ([]domain.User, error)
	GetOneByUsername(username string) (*domain.User, error)
	GetOneByUsernameAndPassword(username string, password string) (*domain.User, error)
	CreateUser(*domain.User) (int, error)
}

var UserEntity IUser

type userEntity struct {
	db *sql.DB
}

func NewUserEntity(db *sql.DB) IUser {
	UserEntity = &userEntity{db: db}
	return UserEntity
}

func (entity *userEntity) GetOneByUsername(username string) (*domain.User, error) {
	sqlStatementSelect := `SELECT * FROM users WHERE username=$1;`
	var user domain.User
	row := entity.db.QueryRow(sqlStatementSelect, username)
	errSelect := row.Scan(&user.ID, &user.Username, &user.Password)
	if errSelect != nil && errSelect == sql.ErrNoRows {
		return nil, errors.New("user not found")
	} else {
		return &user, nil
	}
}

func (entity *userEntity) GetOneByUsernameAndPassword(username string, password string) (*domain.User, error) {
	sqlStatementSelect := `SELECT * FROM users WHERE username=$1 and password=$2;`
	var user domain.User
	row := entity.db.QueryRow(sqlStatementSelect, username, password)
	errSelect := row.Scan(&user.ID, &user.Username, &user.Password)
	if errSelect != nil && errSelect == sql.ErrNoRows {
		return nil, errors.New("user not found")
	} else {
		return &user, nil
	}
}

func (entity *userEntity) CreateUser(user *domain.User) (int, error) {
	insertStatement := `
	INSERT INTO users (username, password)
	VALUES ($1, $2)
	RETURNING id`
	id := 0
	err := entity.db.QueryRow(insertStatement, user.Username, user.Password).Scan(&id)
	if err != nil {
		return -1, errors.New("user cannot be created")
	}
	return id, nil
}
