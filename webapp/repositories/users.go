package repositories

import (
	"context"
	"webapp-go/webapp/models"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type UsersRepository interface {
	GetUser(c context.Context, id uuid.UUID) (models.User, error)
	GetUserByLogin(c context.Context, githubUsername string) (models.User, error)
	CreateUser(c context.Context, user models.User) (models.User, error)
	UpdateUser(c context.Context, id uuid.UUID, user models.User) (models.User, error)
}

type usersRepository struct {
	db *bun.DB
}

func NewUserRepository(db *bun.DB) UsersRepository {
	return usersRepository{db}
}

func (this usersRepository) GetUser(c context.Context, id uuid.UUID) (user models.User, err error) {
	err = this.db.NewSelect().Model(&user).Where("id = ?", id).Scan(c)

	return
}

func (this usersRepository) GetUserByLogin(c context.Context, githubUsername string) (user models.User, err error) {
	err = this.db.NewSelect().Model(&user).Where("github_username = ?", githubUsername).Scan(c)

	return
}

func (this usersRepository) CreateUser(c context.Context, user models.User) (models.User, error) {
	_, err := this.db.NewInsert().Model(&user).Exec(c)

	return user, err
}

func (this usersRepository) UpdateUser(c context.Context, id uuid.UUID, user models.User) (models.User, error) {
	user.ID = id

	_, err := this.db.NewUpdate().Model(&user).OmitZero().WherePK().Exec(c)

	return user, err
}
