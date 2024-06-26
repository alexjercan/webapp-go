package services

import (
	"context"
	"webapp-go/webapp/models"
	"webapp-go/webapp/repositories"

	"github.com/google/uuid"
)

type UsersService interface {
	CreateOrUpdateUser(c context.Context, user models.GitHubUser) (models.User, error)
	CreateAnonymousUser(c context.Context) (models.User, error)
	GetUser(c context.Context, id uuid.UUID) (models.User, error)
}

type usersService struct {
	usersRepository repositories.UsersRepository
}

func NewUsersService(usersRepository repositories.UsersRepository) UsersService {
	return usersService{usersRepository}
}

func (this usersService) CreateOrUpdateUser(c context.Context, user models.GitHubUser) (u models.User, err error) {
	u, err = this.usersRepository.GetUserByLogin(c, user.Login)

	if err != nil {
		u = models.NewUser(user)
		u, err = this.usersRepository.CreateUser(c, u)
	} else {
		u.Name = user.Name
		u, err = this.usersRepository.UpdateUser(c, u.ID, u)
	}

	return
}

func (this usersService) CreateAnonymousUser(c context.Context) (u models.User, err error) {
	u = models.NewAnonymousUser()
	u, err = this.usersRepository.CreateUser(c, u)

	return
}

func (this usersService) GetUser(c context.Context, id uuid.UUID) (u models.User, err error) {
	u, err = this.usersRepository.GetUser(c, id)

	return
}
