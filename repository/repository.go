package repository

import (
	"context"
	"fmt"

	// "github.com/aklinkert/go-logging"
	"github.com/driver005/database"
	"github.com/driver005/database/logger"
)

type Repositories struct {
	logger       logger.Interface
	db           *database.DB
	defaultJoins []string
	ctx 		context.Context
}

// NewRepositories returns a new base repositories that implements TransactionRepositories
func NewRepositories(ctx context.Context, db *database.DB, logger logger.Interface, defaultJoins ...string) TransactionRepository {
	return &Repositories{
		defaultJoins: defaultJoins,
		logger:       logger,
		db:           db,
		ctx: ctx,
	}
}

func (r *Repositories) DB() *database.DB {
	return r.DBWithPreloads(nil)
}

func (r *Repositories) GetAll(target interface{}, preloads ...string) error {
	r.logger.Info(r.ctx, "Executing GetAll on %T", target)

	res := r.DBWithPreloads(preloads).
		Unscoped().
		Find(target)

	return r.HandleError(res)
}

func (r *Repositories) GetBatch(target interface{}, limit, offset int, preloads ...string) error {
	r.logger.Info(r.ctx,"Executing GetBatch on %T", target)

	res := r.DBWithPreloads(preloads).
		Unscoped().
		Limit(limit).
		Offset(offset).
		Find(target)

	return r.HandleError(res)
}

func (r *Repositories) GetWhere(target interface{}, condition string, preloads ...string) error {
	r.logger.Info(r.ctx,"Executing GetWhere on %T with %v ", target, condition)

	res := r.DBWithPreloads(preloads).
		Where(condition).
		Find(target)

	return r.HandleError(res)
}

func (r *Repositories) GetWhereBatch(target interface{}, condition string, limit, offset int, preloads ...string) error {
	r.logger.Info(r.ctx,"Executing GetWhere on %T with %v ", target, condition)

	res := r.DBWithPreloads(preloads).
		Where(condition).
		Limit(limit).
		Offset(offset).
		Find(target)

	return r.HandleError(res)
}

func (r *Repositories) GetByField(target interface{}, field string, value interface{}, preloads ...string) error {
	r.logger.Info(r.ctx,"Executing GetByField on %T with %v = %v", target, field, value)

	res := r.DBWithPreloads(preloads).
		Where(fmt.Sprintf("%v = ?", field), value).
		Find(target)

	return r.HandleError(res)
}

func (r *Repositories) GetByFields(target interface{}, filters map[string]interface{}, preloads ...string) error {
	r.logger.Info(r.ctx,"Executing GetByField on %T with filters = %+v", target, filters)

	db := r.DBWithPreloads(preloads)
	for field, value := range filters {
		db = db.Where(fmt.Sprintf("%v = ?", field), value)
	}

	res := db.Find(target)

	return r.HandleError(res)
}

func (r *Repositories) GetByFieldBatch(target interface{}, field string, value interface{}, limit, offset int, preloads ...string) error {
	r.logger.Info(r.ctx,"Executing GetByField on %T with %v = %v", target, field, value)

	res := r.DBWithPreloads(preloads).
		Where(fmt.Sprintf("%v = ?", field), value).
		Limit(limit).
		Offset(offset).
		Find(target)

	return r.HandleError(res)
}

func (r *Repositories) GetByFieldsBatch(target interface{}, filters map[string]interface{}, limit, offset int, preloads ...string) error {
	r.logger.Info(r.ctx,"Executing GetByField on %T with filters = %+v", target, filters)

	db := r.DBWithPreloads(preloads)
	for field, value := range filters {
		db = db.Where(fmt.Sprintf("%v = ?", field), value)
	}

	res := db.
		Limit(limit).
		Offset(offset).
		Find(target)

	return r.HandleError(res)
}

func (r *Repositories) GetOneByField(target interface{}, field string, value interface{}, preloads ...string) error {
	r.logger.Info(r.ctx,"Executing GetOneByField on %T with %v = %v", target, field, value)

	res := r.DBWithPreloads(preloads).
		Where(fmt.Sprintf("%v = ?", field), value).
		First(target)

	return r.HandleOneError(res)
}

func (r *Repositories) GetOneByFields(target interface{}, filters map[string]interface{}, preloads ...string) error {
	r.logger.Info(r.ctx,"Executing FindOneByField on %T with filters = %+v", target, filters)

	db := r.DBWithPreloads(preloads)
	for field, value := range filters {
		db = db.Where(fmt.Sprintf("%v = ?", field), value)
	}

	res := db.First(target)
	return r.HandleOneError(res)
}

func (r *Repositories) GetOneByID(target interface{}, id string, preloads ...string) error {
	r.logger.Info(r.ctx,"Executing GetOneByID on %T with ID %v", target, id)

	res := r.DBWithPreloads(preloads).
		Where("id = ?", id).
		First(target)

	return r.HandleOneError(res)
}

func (r *Repositories) Create(target interface{}) error {
	r.logger.Info(r.ctx,"Executing Create on %T", target)

	res := r.db.Create(target)
	return r.HandleError(res)
}

func (r *Repositories) CreateTx(target interface{}, tx *database.DB) error {
	r.logger.Info(r.ctx,"Executing Create on %T", target)

	res := tx.Create(target)
	return r.HandleError(res)
}

func (r *Repositories) Save(target interface{}) error {
	r.logger.Info(r.ctx,"Executing Save on %T", target)

	res := r.db.Save(target)
	return r.HandleError(res)
}

func (r *Repositories) SaveTx(target interface{}, tx *database.DB) error {
	r.logger.Info(r.ctx,"Executing Save on %T", target)

	res := tx.Save(target)
	return r.HandleError(res)
}

func (r *Repositories) Delete(target interface{}) error {
	r.logger.Info(r.ctx,"Executing Delete on %T", target)

	res := r.db.Delete(target)
	return r.HandleError(res)
}

func (r *Repositories) DeleteTx(target interface{}, tx *database.DB) error {
	r.logger.Info(r.ctx,"Executing Delete on %T", target)

	res := tx.Delete(target)
	return r.HandleError(res)
}

func (r *Repositories) HandleError(res *database.DB) error {
	if res.Error != nil && res.Error != database.ErrRecordNotFound {
		err := fmt.Errorf("Error: %w", res.Error)
		r.logger.Error(r.ctx, "%v", err)
		return err
	}

	return nil
}

func (r *Repositories) HandleOneError(res *database.DB) error {
	if err := r.HandleError(res); err != nil {
		return err
	}

	if res.RowsAffected != 1 {
		return ErrNotFound
	}

	return nil
}

func (r *Repositories) DBWithPreloads(preloads []string) *database.DB {
	dbConn := r.db

	for _, join := range r.defaultJoins {
		dbConn = dbConn.Joins(join)
	}

	for _, preload := range preloads {
		dbConn = dbConn.Preload(preload)
	}

	return dbConn
}
