package callbacks

import "github.com/driver005/database"

var (
	createClauses = []string{"INSERT", "VALUES", "ON CONFLICT"}
	queryClauses  = []string{"SELECT", "FROM", "WHERE", "GROUP BY", "ORDER BY", "LIMIT", "FOR"}
	updateClauses = []string{"UPDATE", "SET", "WHERE"}
	deleteClauses = []string{"DELETE", "FROM", "WHERE"}
)

type Config struct {
	LastInsertIDReversed bool
	CreateClauses        []string
	QueryClauses         []string
	UpdateClauses        []string
	DeleteClauses        []string
}

func RegisterDefaultCallbacks(db *database.DB, config *Config) {
	enableTransaction := func(db *database.DB) bool {
		return !db.SkipDefaultTransaction
	}

	if len(config.CreateClauses) == 0 {
		config.CreateClauses = createClauses
	}
	if len(config.QueryClauses) == 0 {
		config.QueryClauses = queryClauses
	}
	if len(config.DeleteClauses) == 0 {
		config.DeleteClauses = deleteClauses
	}
	if len(config.UpdateClauses) == 0 {
		config.UpdateClauses = updateClauses
	}

	createCallback := db.Callback().Create()
	createCallback.Match(enableTransaction).Register("database:begin_transaction", BeginTransaction)
	createCallback.Register("database:before_create", BeforeCreate)
	createCallback.Register("database:save_before_associations", SaveBeforeAssociations(true))
	createCallback.Register("database:create", Create(config))
	createCallback.Register("database:save_after_associations", SaveAfterAssociations(true))
	createCallback.Register("database:after_create", AfterCreate)
	createCallback.Match(enableTransaction).Register("database:commit_or_rollback_transaction", CommitOrRollbackTransaction)
	createCallback.Clauses = config.CreateClauses

	queryCallback := db.Callback().Query()
	queryCallback.Register("database:query", Query)
	queryCallback.Register("database:preload", Preload)
	queryCallback.Register("database:after_query", AfterQuery)
	queryCallback.Clauses = config.QueryClauses

	deleteCallback := db.Callback().Delete()
	deleteCallback.Match(enableTransaction).Register("database:begin_transaction", BeginTransaction)
	deleteCallback.Register("database:before_delete", BeforeDelete)
	deleteCallback.Register("database:delete_before_associations", DeleteBeforeAssociations)
	deleteCallback.Register("database:delete", Delete(config))
	deleteCallback.Register("database:after_delete", AfterDelete)
	deleteCallback.Match(enableTransaction).Register("database:commit_or_rollback_transaction", CommitOrRollbackTransaction)
	deleteCallback.Clauses = config.DeleteClauses

	updateCallback := db.Callback().Update()
	updateCallback.Match(enableTransaction).Register("database:begin_transaction", BeginTransaction)
	updateCallback.Register("database:setup_reflect_value", SetupUpdateReflectValue)
	updateCallback.Register("database:before_update", BeforeUpdate)
	updateCallback.Register("database:save_before_associations", SaveBeforeAssociations(false))
	updateCallback.Register("database:update", Update(config))
	updateCallback.Register("database:save_after_associations", SaveAfterAssociations(false))
	updateCallback.Register("database:after_update", AfterUpdate)
	updateCallback.Match(enableTransaction).Register("database:commit_or_rollback_transaction", CommitOrRollbackTransaction)
	updateCallback.Clauses = config.UpdateClauses

	rowCallback := db.Callback().Row()
	rowCallback.Register("database:row", RowQuery)
	rowCallback.Clauses = config.QueryClauses

	rawCallback := db.Callback().Raw()
	rawCallback.Register("database:raw", RawExec)
	rawCallback.Clauses = config.QueryClauses
}
