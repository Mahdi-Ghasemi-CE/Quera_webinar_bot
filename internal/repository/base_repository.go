package repository

import (
	"Quera_webinar_bot/internal/filter"
	"context"
)

type BaseRepository[TEntity any] interface {
	Create(ctx context.Context, entity TEntity) (TEntity, error)
	//Update(ctx context.Context, id int, entity map[string]interface{}) (TEntity, error)
	Delete(ctx context.Context, id int) error
	GetById(ctx context.Context, id int) (TEntity, error)
	GetByFilter(req filter.QueryOptions) (int64, *[]TEntity, error)
}
