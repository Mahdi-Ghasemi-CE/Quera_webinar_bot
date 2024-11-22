package usecase

import (
	"Quera_webinar_bot/config"
	"Quera_webinar_bot/internal/filter"
	"Quera_webinar_bot/internal/persistence/repository"
	"Quera_webinar_bot/tools"
	"Quera_webinar_bot/usecase/dto"
	"context"
	"github.com/jinzhu/copier"
)

type BaseUsecase[TEntity any, TCreate any, TUpdate any, TResponse any] struct {
	repository repository.BaseRepository[TEntity]
}

func NewBaseUsecase[TEntity any, TCreate any, TUpdate any, TResponse any](cfg *config.Config, repository repository.BaseRepository[TEntity]) *BaseUsecase[TEntity, TCreate, TUpdate, TResponse] {
	return &BaseUsecase[TEntity, TCreate, TUpdate, TResponse]{
		repository: repository,
	}
}

func (u *BaseUsecase[TEntity, TCreate, TUpdate, TResponse]) Create(ctx context.Context, req TCreate) (TResponse, error) {
	var response TResponse
	entity, _ := tools.TypeConverter[TEntity](req)

	entity, err := u.repository.Create(ctx, entity)
	if err != nil {
		return response, err
	}

	response, _ = tools.TypeConverter[TResponse](entity)
	return response, nil
}

func (u *BaseUsecase[TEntity, TCreate, TUpdate, TResponse]) Delete(ctx context.Context, id int) error {

	return u.repository.Delete(ctx, id)
}

func (u *BaseUsecase[TEntity, TCreate, TUpdate, TResponse]) GetById(ctx context.Context, id int) (TResponse, error) {
	var response TResponse
	entity, err := u.repository.GetById(ctx, id)
	if err != nil {
		return response, err
	}
	return tools.TypeConverter[TResponse](entity)
}

func (u *BaseUsecase[TEntity, TCreate, TUpdate, TResponse]) GetByFilter(ctx context.Context, req filter.QueryOptions) (*dto.UsecaseResponse[TResponse], error) {
	var response *dto.UsecaseResponse[TResponse]
	count, entities, err := u.repository.GetByFilter(req)
	if err != nil {
		return response, err
	}

	var items *[]TResponse

	err = copier.Copy(entities, items)
	if err != nil {
		return response, err
	}

	response.Items = items
	response.TotalRows = count

	return response, nil
}
