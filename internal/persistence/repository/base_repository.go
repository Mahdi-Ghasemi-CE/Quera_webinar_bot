package repository

import (
	"context"
	"database/sql"
	"gorm.io/gorm"
	"menu_yar/config"
	"menu_yar/internal/enum"
	"menu_yar/internal/service_errors"
	"menu_yar/pkg/filter"
	"menu_yar/pkg/logging"
	"menu_yar/pkg/persistence/database"
	common "menu_yar/tools"
	"time"
)

const softDeleteExp string = "id = ?"

type BaseRepository[TEntity any] struct {
	database *gorm.DB
	logger   logging.Logger
	preloads []string
}

func NewBaseRepository[TEntity any](cfg *config.Config, preloads []string) *BaseRepository[TEntity] {
	return &BaseRepository[TEntity]{
		database: database.GetDb(),
		logger:   logging.NewLogger(cfg),
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
		r.logger.Error(logging.Postgres, logging.Insert, err.Error(), nil)
		return entity, err
	}
	tx.Commit()

	return entity, nil
}

func (r BaseRepository[TEntity]) Update(ctx context.Context, id int, entity map[string]interface{}) (TEntity, error) {
	snakeMap := map[string]interface{}{}
	for k, v := range entity {
		snakeMap[common.ToSnakeCase(k)] = v
	}
	snakeMap["modified_by"] = enum.AdminId
	snakeMap["modified_at"] = sql.NullTime{Valid: true, Time: time.Now().UTC()}
	model := new(TEntity)
	tx := r.database.WithContext(ctx).Begin()
	if err := tx.Model(model).
		Where(softDeleteExp, id).
		Updates(snakeMap).
		Error; err != nil {
		tx.Rollback()
		r.logger.Error(logging.Postgres, logging.Update, err.Error(), nil)
		return *model, err
	}
	tx.Commit()
	return *model, nil
}

func (r BaseRepository[TEntity]) Delete(ctx context.Context, id int) error {
	tx := r.database.WithContext(ctx).Begin()

	model := new(TEntity)

	/*	deleteMap := map[string]interface{}{
		"deleted_by": &sql.NullInt64{Int64: int64(ctx.Value(constant.UserIdKey).(float64)), Valid: true},
		"deleted_at": sql.NullTime{Valid: true, Time: time.Now().UTC()},
	}*/

	if ctx.Value(enum.UserIdKey) == nil {
		return &service_errors.ServiceError{EndUserMessage: string(enum.PermissionDenied)}
	}
	if cnt := tx.
		Model(model).
		Where(softDeleteExp, id).
		/* Updates(deleteMap). */
		RowsAffected; cnt == 0 {
		tx.Rollback()
		r.logger.Error(logging.Postgres, logging.Update, string(enum.RecordNotFound), nil)
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

/*func (r BaseRepository[TEntity]) GetByFilter(ctx context.Context, req filter.QueryOptions) (int64, *[]TEntity, error) {
	model := new(TEntity)
	var items *[]TEntity

	db := database.Preload(r.database, r.preloads)
	query := database.GenerateDynamicQuery[TEntity](&req.DynamicFilter)
	sort := database.GenerateDynamicSort[TEntity](&req.DynamicFilter)
	var totalRows int64 = 0

	db.
		Model(model).
		Where(query).
		Count(&totalRows)

	err := db.
		Where(query).
		Offset(req.GetOffset()).
		Limit(req.GetPageSize()).
		Order(sort).
		Find(&items).
		Error

	if err != nil {
		return 0, &[]TEntity{}, err
	}
	return totalRows, items, err

}*/
