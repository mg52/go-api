package repository

import (
	"errors"
	"github.com/mg52/go-api/domain"
)

type IUser interface {
	GetAll() ([]domain.User, error)
	GetOneByUsername(username string) (*domain.User, error)
	GetOneByUsernameAndPassword(username string, password string) (*domain.User, error)
}

var UserEntity IUser

type userEntity struct {
	users []domain.User
}

func NewUserEntity() IUser {
	var users []domain.User
	users = append(users, domain.User{
		ID:       1,
		Name:     "aaaaa",
		Password: "passaaaaa",
	})
	users = append(users, domain.User{
		ID:       2,
		Name:     "bbbbb",
		Password: "passbbbbb",
	})
	users = append(users, domain.User{
		ID:       3,
		Name:     "ccccc",
		Password: "passccccc",
	})
	users = append(users, domain.User{
		ID:       4,
		Name:     "ddddd",
		Password: "passddddd",
	})
	UserEntity = &userEntity{users: users}
	return UserEntity
}

func (entity *userEntity) GetAll() ([]domain.User, error) {
	return entity.users, nil
}

func (entity *userEntity) GetOneByUsername(username string) (*domain.User, error) {
	for _, theUser := range entity.users {
		if theUser.Name == username {
			return &theUser, nil
			break
		}
	}
	return nil, errors.New("user not found")
}

func (entity *userEntity) GetOneByUsernameAndPassword(username string, password string) (*domain.User, error) {
	for _, theUser := range entity.users {
		if theUser.Name == username {
			if theUser.Password == password {
				return &theUser, nil
				break
			}
		}
	}
	return nil, errors.New("user not found")
}
