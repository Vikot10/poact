package articleservice

import "context"

type ArticleService struct {
	storage storage.Storage
}

type ArticleCRUD interface {
	AddArticle(ctx context.Context, article *dto.Article) error
	Create(ctx context.Context, article *dto.Article) error
	Update(ctx context.Context, article *dto.Article) error
	Delete(ctx context.Context, id int64) error
}
