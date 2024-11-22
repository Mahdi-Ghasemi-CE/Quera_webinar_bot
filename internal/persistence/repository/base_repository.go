package repository

import (
	"Quera_webinar_bot/config"
	"Quera_webinar_bot/internal/enum"
	"Quera_webinar_bot/internal/filter"
	"Quera_webinar_bot/internal/persistence/database"
	"Quera_webinar_bot/internal/service_errors"
	"context"
	"gorm.io/gorm"
)

const softDeleteExp string = "id = ?"

type BaseRepository[TEntity any] struct {
	database *gorm.DB
	preloads []string
}

func NewBaseRepository[TEntity any](cfg *config.Config, preloads []string) *BaseRepository[TEntity] {
	return &BaseRepository[TEntity]{
		database: database.GetDb(),
		preloads: preloads,
	}
}

func (r BaseRepository[TEntity]) Create(ctx context.Context, entity TEntity) (TEntity, error) {
	tx := r.database.WithContext(ctx).Begin()
	err := tx.
		Create(&entity).
		Error
	if err != nil {
		tx.Rollback()
		return entity, err
	}
	tx.Commit()

	return entity, nil
}

func (r BaseRepository[TEntity]) Delete(ctx context.Context, id int) error {
	tx := r.database.WithContext(ctx).Begin()

	model := new(TEntity)

	if cnt := tx.
		Model(model).
		Where(softDeleteExp, id).
		/* Updates(deleteMap). */
		RowsAffected; cnt == 0 {
		tx.Rollback()
		return &service_errors.ServiceError{EndUserMessage: string(enum.RecordNotFound)}
	}
	tx.Commit()
	return nil
}

func (r BaseRepository[TEntity]) GetById(ctx context.Context, id int) (TEntity, error) {
	model := new(TEntity)
	db := filter.ApplyPreloads(r.database, r.preloads)
	err := db.
		Where(softDeleteExp, id).
		First(model).
		Error
	if err != nil {
		return *model, err
	}
	return *model, nil
}

func (r *BaseRepository[TEntity]) GetByFilter(options filter.QueryOptions) (int64, *[]TEntity, error) {
	var entities *[]TEntity
	var totalRows int64

	query := filter.BuildQuery(r.database, options)

	result := query.Find(&entities).Count(&totalRows)

	return totalRows, entities, result.Error
}
