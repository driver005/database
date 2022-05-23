package database

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"reflect"

	"github.com/driver005/database/schema"
	"github.com/driver005/database/types"
)

type DeletedAt sql.NullTime

// Scan implements the Scanner interface.
func (n *DeletedAt) Scan(value interface{}) error {
	return (*sql.NullTime)(n).Scan(value)
}

// Value implements the driver Valuer interface.
func (n DeletedAt) Value() (driver.Value, error) {
	if !n.Valid {
		return nil, nil
	}
	return n.Time, nil
}

func (n DeletedAt) MarshalJSON() ([]byte, error) {
	if n.Valid {
		return json.Marshal(n.Time)
	}
	return json.Marshal(nil)
}

func (n *DeletedAt) UnmarshalJSON(b []byte) error {
	if string(b) == "null" {
		n.Valid = false
		return nil
	}
	err := json.Unmarshal(b, &n.Time)
	if err == nil {
		n.Valid = true
	}
	return err
}

func (DeletedAt) QueryClauses(f *schema.Field) []types.Interface {
	return []types.Interface{SoftDeleteQueryClause{Field: f}}
}

type SoftDeleteQueryClause struct {
	Field *schema.Field
}

func (sd SoftDeleteQueryClause) Name() string {
	return ""
}

func (sd SoftDeleteQueryClause) Build(types.Builder) {
}

func (sd SoftDeleteQueryClause) MergeClause(*types.Type) {
}

func (sd SoftDeleteQueryClause) ModifyStatement(stmt *Statement) {
	if _, ok := stmt.Types["soft_delete_enabled"]; !ok && !stmt.Statement.Unscoped {
		if c, ok := stmt.Types["WHERE"]; ok {
			if where, ok := c.Expression.(types.Where); ok && len(where.Exprs) >= 1 {
				for _, expr := range where.Exprs {
					if orCond, ok := expr.(types.OrConditions); ok && len(orCond.Exprs) == 1 {
						where.Exprs = []types.Expression{types.And(where.Exprs...)}
						c.Expression = where
						stmt.Types["WHERE"] = c
						break
					}
				}
			}
		}

		stmt.AddClause(types.Where{Exprs: []types.Expression{
			types.Eq{Column: types.Column{Table: types.CurrentTable, Name: sd.Field.DBName}, Value: nil},
		}})
		stmt.Types["soft_delete_enabled"] = types.Type{}
	}
}

func (DeletedAt) UpdateClauses(f *schema.Field) []types.Interface {
	return []types.Interface{SoftDeleteUpdateClause{Field: f}}
}

type SoftDeleteUpdateClause struct {
	Field *schema.Field
}

func (sd SoftDeleteUpdateClause) Name() string {
	return ""
}

func (sd SoftDeleteUpdateClause) Build(types.Builder) {
}

func (sd SoftDeleteUpdateClause) MergeClause(*types.Type) {
}

func (sd SoftDeleteUpdateClause) ModifyStatement(stmt *Statement) {
	if stmt.SQL.Len() == 0 && !stmt.Statement.Unscoped {
		SoftDeleteQueryClause(sd).ModifyStatement(stmt)
	}
}

func (DeletedAt) DeleteClauses(f *schema.Field) []types.Interface {
	return []types.Interface{SoftDeleteDeleteClause{Field: f}}
}

type SoftDeleteDeleteClause struct {
	Field *schema.Field
}

func (sd SoftDeleteDeleteClause) Name() string {
	return ""
}

func (sd SoftDeleteDeleteClause) Build(types.Builder) {
}

func (sd SoftDeleteDeleteClause) MergeClause(*types.Type) {
}

func (sd SoftDeleteDeleteClause) ModifyStatement(stmt *Statement) {
	if stmt.SQL.Len() == 0 && !stmt.Statement.Unscoped {
		curTime := stmt.DB.NowFunc()
		stmt.AddClause(types.Set{{Column: types.Column{Name: sd.Field.DBName}, Value: curTime}})
		stmt.SetColumn(sd.Field.DBName, curTime, true)

		if stmt.Schema != nil {
			_, queryValues := schema.GetIdentityFieldValuesMap(stmt.Context, stmt.ReflectValue, stmt.Schema.PrimaryFields)
			column, values := schema.ToQueryValues(stmt.Table, stmt.Schema.PrimaryFieldDBNames, queryValues)

			if len(values) > 0 {
				stmt.AddClause(types.Where{Exprs: []types.Expression{types.IN{Column: column, Values: values}}})
			}

			if stmt.ReflectValue.CanAddr() && stmt.Dest != stmt.Model && stmt.Model != nil {
				_, queryValues = schema.GetIdentityFieldValuesMap(stmt.Context, reflect.ValueOf(stmt.Model), stmt.Schema.PrimaryFields)
				column, values = schema.ToQueryValues(stmt.Table, stmt.Schema.PrimaryFieldDBNames, queryValues)

				if len(values) > 0 {
					stmt.AddClause(types.Where{Exprs: []types.Expression{types.IN{Column: column, Values: values}}})
				}
			}
		}

		SoftDeleteQueryClause(sd).ModifyStatement(stmt)
		stmt.AddClauseIfNotExists(types.Update{})
		stmt.Build(stmt.DB.Callback().Update().Clauses...)
	}
}