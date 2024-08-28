package articleservice

import (
	"context"

	"github.com/Vikot10/viarticles/internal/dto"
)

type ArticleService struct {
	ArticleStore ArticleStore
}

type ArticleStore interface {
	AddArticle(ctx context.Context, article *dto.Article) error
	Create(ctx context.Context, article *dto.Article) error
	Update(ctx context.Context, article *dto.Article) error
	Delete(ctx context.Context, id int64) error
}

type FaveStore interface {
	AddFave(ctx context.Context, fave *dto.Fave) error
	GetFaves(ctx context.Context) ([]*dto.Fave, error)
}

func New(store ArticleStore) *ArticleService {
	return &ArticleService{
		ArticleStore: store,
	}
}
