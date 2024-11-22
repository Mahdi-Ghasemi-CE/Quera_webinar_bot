package dto

import (
	"database/sql"
	"time"
)

type UsecaseResponse[T any] struct {
	TotalRows int64 `json:"totalRows"`
	Items     *[]T  `json:"items"`
}

type BaseModel struct {
	Id int64 `mapper:"_id"`

	CreatedAt time.Time
	CreatedBy int

	ModifiedAt sql.NullTime
	ModifiedBy *sql.NullInt64

	IsActive bool
}
