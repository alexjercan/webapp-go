package repositories

import (
	"context"
	"webapp-go/webapp/models"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type UsersRepository interface {
	GetUser(c context.Context, id uuid.UUID) (models.User, error)
	GetUsers(c context.Context) ([]models.User, error)
	CreateUser(c context.Context, user models.User) (models.User, error)
	UpdateUser(c context.Context, id uuid.UUID, user models.User) (models.User, error)
	DeleteUser(c context.Context, id uuid.UUID) (uuid.UUID, error)
}

type usersRepository struct {
	db *bun.DB
}

func NewUserRepository(db *bun.DB) UsersRepository {
	return usersRepository{db}
}

func (p usersRepository) GetUser(c context.Context, id uuid.UUID) (user models.User, err error) {
	err = p.db.NewSelect().Model(&user).Where("id = ?", id).Scan(c);

    return
}

func (p usersRepository) GetUsers(c context.Context) (users []models.User, err error) {
	err = p.db.NewSelect().Model(&users).Scan(c)

    return
}

func (p usersRepository) CreateUser(c context.Context, user models.User) (models.User, error) {
	_, err := p.db.NewInsert().Model(&user).Exec(c)

    return user, err
}

func (p usersRepository) UpdateUser(c context.Context, id uuid.UUID, user models.User) (models.User, error) {
	user.ID = id

    _, err := p.db.NewUpdate().Model(&user).OmitZero().WherePK().Exec(c)

    return user, err
}

func (p usersRepository) DeleteUser(c context.Context, id uuid.UUID) (uuid.UUID, error) {
    _, err := p.db.NewDelete().Model((*models.User)(nil)).Where("id = ?", id).Exec(c)

    return id, err
}
