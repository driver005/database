package repository

import "github.com/driver005/database"

// Repository is a generic DB handler that cares about default error handling
type Repository interface {
	GetAll(target interface{}, preloads ...string) error
	GetBatch(target interface{}, limit, offset int, preloads ...string) error

	GetWhere(target interface{}, condition string, preloads ...string) error
	GetWhereBatch(target interface{}, condition string, limit, offset int, preloads ...string) error

	GetByField(target interface{}, field string, value interface{}, preloads ...string) error
	GetByFields(target interface{}, filters map[string]interface{}, preloads ...string) error
	GetByFieldBatch(target interface{}, field string, value interface{}, limit, offset int, preloads ...string) error
	GetByFieldsBatch(target interface{}, filters map[string]interface{}, limit, offset int, preloads ...string) error

	GetOneByField(target interface{}, field string, value interface{}, preloads ...string) error
	GetOneByFields(target interface{}, filters map[string]interface{}, preloads ...string) error

	// GetOneByID assumes you have a PK column "id" which is a UUID. If this is not the case just ignore the method
	// and add a custom struct with this Repository embedded.
	GetOneByID(target interface{}, id string, preloads ...string) error

	Create(target interface{}) error
	Save(target interface{}) error
	Delete(target interface{}) error

	DB() *database.DB
	DBWithPreloads(preloads []string) *database.DB
	HandleError(res *database.DB) error
	HandleOneError(res *database.DB) error
}

// TransactionRepository extends Repository with modifier functions that accept a transaction
type TransactionRepository interface {
	Repository
	CreateTx(target interface{}, tx *database.DB) error
	SaveTx(target interface{}, tx *database.DB) error
	DeleteTx(target interface{}, tx *database.DB) error
}
