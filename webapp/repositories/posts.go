package repositories

import (
	"context"
	"webapp-go/webapp/models"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type PostsRepository interface {
	GetPost(c context.Context, slug uuid.UUID) (models.Post, error)
	GetPosts(c context.Context) ([]models.Post, error)
	CreatePost(c context.Context, post models.Post) (models.Post, error)
	UpdatePost(c context.Context, slug uuid.UUID, post models.Post) (models.Post, error)
	DeletePost(c context.Context, slug uuid.UUID) (uuid.UUID, error)
}

type postsRepository struct {
	db *bun.DB
}

func NewPostsRepository(db *bun.DB) PostsRepository {
	return postsRepository{db}
}

func (p postsRepository) GetPost(c context.Context, slug uuid.UUID) (post models.Post, err error) {
	err = p.db.NewSelect().Model(&post).Where("slug = ?", slug).Scan(c);

    return
}

func (p postsRepository) GetPosts(c context.Context) (posts []models.Post, err error) {
	err = p.db.NewSelect().Model(&posts).Scan(c)

    return
}

func (p postsRepository) CreatePost(c context.Context, post models.Post) (models.Post, error) {
	_, err := p.db.NewInsert().Model(&post).Exec(c)

    return post, err
}

func (p postsRepository) UpdatePost(c context.Context, slug uuid.UUID, post models.Post) (models.Post, error) {
	post.Slug = slug

    _, err := p.db.NewUpdate().Model(&post).OmitZero().WherePK().Exec(c)

    return post, err
}

func (p postsRepository) DeletePost(c context.Context, slug uuid.UUID) (uuid.UUID, error) {
    _, err := p.db.NewDelete().Model((*models.Post)(nil)).Where("slug = ?", slug).Exec(c)

    return slug, err
}
