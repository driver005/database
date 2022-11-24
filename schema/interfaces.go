package schema

import (
	"github.com/driver005/database/clause"
)

// DBDataTypeInterface database data type interface
type DBDataTypeInterface interface {
	DBDataType() string
}

// FieldNewValuePool field new scan value pool
type FieldNewValuePool interface {
	Get() interface{}
	Put(interface{})
}

// CreateClausesInterface create clauses interface
type CreateClausesInterface interface {
	CreateClauses(*Field) []clause.Interface
}

// QueryClausesInterface query clauses interface
type QueryClausesInterface interface {
	QueryClauses(*Field) []clause.Interface
}

// UpdateClausesInterface update clauses interface
type UpdateClausesInterface interface {
	UpdateClauses(*Field) []clause.Interface
}

// DeleteClausesInterface delete clauses interface
type DeleteClausesInterface interface {
	DeleteClauses(*Field) []clause.Interface
}
