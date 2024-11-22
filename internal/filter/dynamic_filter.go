package filter

import (
	"Quera_webinar_bot/internal/enum"
	"gorm.io/gorm"
)

type Filter struct {
	Field    string
	Operator enum.QueryOperation
	Value    interface{}
}

type Sort struct {
	Field      string
	Descending bool
}

type Pagination struct {
	Page     int
	PageSize int
}

type QueryOptions struct {
	Filters    []Filter
	Sorts      []Sort
	Pagination *Pagination
	Preloads   []string
}

func ApplyFilters(db *gorm.DB, filters []Filter) *gorm.DB {
	for _, filter := range filters {
		switch filter.Operator {
		case enum.Equal:
			db = db.Where(filter.Field+" = ?", filter.Value)
		case enum.NoEqual:
			db = db.Where(filter.Field+" != ?", filter.Value)
		case enum.GreaterThan:
			db = db.Where(filter.Field+" > ?", filter.Value)
		case enum.GreaterEqual:
			db = db.Where(filter.Field+" => ?", filter.Value)
		case enum.SmallerThan:
			db = db.Where(filter.Field+" < ?", filter.Value)
		case enum.SmallerEqual:
			db = db.Where(filter.Field+" =< ?", filter.Value)
		case enum.LIKE:
			db = db.Where(filter.Field+" LIKE ?", filter.Value)
		case enum.NoLIKE:
			db = db.Where(filter.Field+" NOT LIKE ?", filter.Value)
		}
	}
	return db
}

func ApplySorting(db *gorm.DB, sorts []Sort) *gorm.DB {
	for _, sort := range sorts {
		if sort.Descending {
			db = db.Order(sort.Field + " DESC")
		} else {
			db = db.Order(sort.Field + " ASC")
		}
	}
	return db
}

func ApplyPagination(db *gorm.DB, pagination *Pagination) *gorm.DB {
	if pagination != nil {
		offset := (pagination.Page - 1) * pagination.PageSize
		db = db.Offset(offset).Limit(pagination.PageSize)
	}
	return db
}

func ApplyPreloads(db *gorm.DB, preloads []string) *gorm.DB {
	for _, preload := range preloads {
		db = db.Preload(preload)
	}
	return db
}

func BuildQuery(db *gorm.DB, options QueryOptions) *gorm.DB {
	db = ApplyFilters(db, options.Filters)
	db = ApplySorting(db, options.Sorts)
	db = ApplyPagination(db, options.Pagination)
	db = ApplyPreloads(db, options.Preloads)
	return db
}
