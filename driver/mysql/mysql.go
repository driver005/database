package mysql

import (
	"context"
	"database/sql"
	"fmt"
	"math"
	"strings"
	"time"

	"github.com/driver005/database"
	"github.com/driver005/database/callbacks"
	"github.com/driver005/database/logger"
	"github.com/driver005/database/migrator"
	"github.com/driver005/database/schema"
	"github.com/driver005/database/types"
	_ "github.com/go-sql-driver/mysql"
)

type Config struct {
	DriverName                string
	ServerVersion             string
	DSN                       string
	Conn                      database.ConnPool
	SkipInitializeWithVersion bool
	DefaultStringSize         uint
	DefaultDatetimePrecision  *int
	DisableDatetimePrecision  bool
	DontSupportRenameIndex    bool
	DontSupportRenameColumn   bool
	DontSupportForShareClause bool
}

type Dialector struct {
	*Config
}

var (
	// CreateClauses create clauses
	CreateClauses = []string{"INSERT", "VALUES", "ON CONFLICT"}
	// QueryClauses query clauses
	QueryClauses = []string{}
	// UpdateClauses update clauses
	UpdateClauses = []string{"UPDATE", "SET", "WHERE", "ORDER BY", "LIMIT"}
	// DeleteClauses delete clauses
	DeleteClauses = []string{"DELETE", "FROM", "WHERE", "ORDER BY", "LIMIT"}

	defaultDatetimePrecision = 3
)

func Open(dsn string) database.Dialector {
	return &Dialector{Config: &Config{DSN: dsn}}
}

func New(config Config) database.Dialector {
	return &Dialector{Config: &config}
}

func (dialector Dialector) Name() string {
	return "mysql"
}

// NowFunc return now func
func (dialector Dialector) NowFunc(n int) func() time.Time {
	return func() time.Time {
		round := time.Second / time.Duration(math.Pow10(n))
		return time.Now().Round(round)
	}
}

func (dialector Dialector) Apply(config *database.Config) error {
	if config.NowFunc == nil {
		if dialector.DefaultDatetimePrecision == nil {
			dialector.DefaultDatetimePrecision = &defaultDatetimePrecision
		}

		// while maintaining the readability of the code, separate the business logic from
		// the general part and leave it to the function to do it here.
		config.NowFunc = dialector.NowFunc(*dialector.DefaultDatetimePrecision)
	}

	return nil
}

func (dialector Dialector) Initialize(db *database.DB) (err error) {
	ctx := context.Background()

	// register callbacks
	callbacks.RegisterDefaultCallbacks(db, &callbacks.Config{
		CreateClauses: CreateClauses,
		QueryClauses:  QueryClauses,
		UpdateClauses: UpdateClauses,
		DeleteClauses: DeleteClauses,
	})

	if dialector.DriverName == "" {
		dialector.DriverName = "mysql"
	}

	if dialector.DefaultDatetimePrecision == nil {
		dialector.DefaultDatetimePrecision = &defaultDatetimePrecision
	}

	if dialector.Conn != nil {
		db.ConnPool = dialector.Conn
	} else {
		db.ConnPool, err = sql.Open(dialector.DriverName, dialector.DSN)
		if err != nil {
			return err
		}
	}

	if !dialector.Config.SkipInitializeWithVersion {
		err = db.ConnPool.QueryRowContext(ctx, "SELECT VERSION()").Scan(&dialector.ServerVersion)
		if err != nil {
			return err
		}

		if strings.Contains(dialector.ServerVersion, "MariaDB") {
			dialector.Config.DontSupportRenameIndex = true
			dialector.Config.DontSupportRenameColumn = true
			dialector.Config.DontSupportForShareClause = true
		} else if strings.HasPrefix(dialector.ServerVersion, "5.6.") {
			dialector.Config.DontSupportRenameIndex = true
			dialector.Config.DontSupportRenameColumn = true
			dialector.Config.DontSupportForShareClause = true
		} else if strings.HasPrefix(dialector.ServerVersion, "5.7.") {
			dialector.Config.DontSupportRenameColumn = true
			dialector.Config.DontSupportForShareClause = true
		} else if strings.HasPrefix(dialector.ServerVersion, "5.") {
			dialector.Config.DisableDatetimePrecision = true
			dialector.Config.DontSupportRenameIndex = true
			dialector.Config.DontSupportRenameColumn = true
			dialector.Config.DontSupportForShareClause = true
		}
	}

	for k, v := range dialector.ClauseBuilders() {
		db.TypesBuilders[k] = v
	}
	return
}

const (
	// ClauseOnConflict for types.TypesBuilder ON CONFLICT key
	ClauseOnConflict = "ON CONFLICT"
	// ClauseValues for types.TypesBuilder VALUES key
	ClauseValues = "VALUES"
	// ClauseValues for types.TypesBuilder FOR key
	ClauseFor = "FOR"
)

func (dialector Dialector) ClauseBuilders() map[string]types.TypesBuilder {
	clauseBuilders := map[string]types.TypesBuilder{
		ClauseOnConflict: func(c types.Type, builder types.Builder) {
			onConflict, ok := c.Expression.(types.OnConflict)
			if !ok {
				c.Build(builder)
				return
			}

			builder.WriteString("ON DUPLICATE KEY UPDATE ")
			if len(onConflict.DoUpdates) == 0 {
				if s := builder.(*database.Statement).Schema; s != nil {
					var column types.Column
					onConflict.DoNothing = false

					if s.PrioritizedPrimaryField != nil {
						column = types.Column{Name: s.PrioritizedPrimaryField.DBName}
					} else if len(s.DBNames) > 0 {
						column = types.Column{Name: s.DBNames[0]}
					}

					if column.Name != "" {
						onConflict.DoUpdates = []types.Assignment{{Column: column, Value: column}}
					}
				}
			}

			for idx, assignment := range onConflict.DoUpdates {
				if idx > 0 {
					builder.WriteByte(',')
				}

				builder.WriteQuoted(assignment.Column)
				builder.WriteByte('=')
				if column, ok := assignment.Value.(types.Column); ok && column.Table == "excluded" {
					column.Table = ""
					builder.WriteString("VALUES(")
					builder.WriteQuoted(column)
					builder.WriteByte(')')
				} else {
					builder.AddVar(builder, assignment.Value)
				}
			}
		},
		ClauseValues: func(c types.Type, builder types.Builder) {
			if values, ok := c.Expression.(types.Values); ok && len(values.Columns) == 0 {
				builder.WriteString("VALUES()")
				return
			}
			c.Build(builder)
		},
	}

	if dialector.Config.DontSupportForShareClause {
		clauseBuilders[ClauseFor] = func(c types.Type, builder types.Builder) {
			if values, ok := c.Expression.(types.Locking); ok && strings.EqualFold(values.Strength, "SHARE") {
				builder.WriteString("LOCK IN SHARE MODE")
				return
			}
			c.Build(builder)
		}
	}

	return clauseBuilders
}

func (dialector Dialector) DefaultValueOf(field *schema.Field) types.Expression {
	return types.Expr{SQL: "DEFAULT"}
}

func (dialector Dialector) Migrator(db *database.DB) database.Migrator {
	return Migrator{
		Migrator: migrator.Migrator{
			Config: migrator.Config{
				DB:        db,
				Dialector: dialector,
			},
		},
		Dialector: dialector,
	}
}

func (dialector Dialector) BindVarTo(writer types.Writer, stmt *database.Statement, v interface{}) {
	writer.WriteByte('?')
}

func (dialector Dialector) QuoteTo(writer types.Writer, str string) {
	var (
		underQuoted, selfQuoted bool
		continuousBacktick      int8
		shiftDelimiter          int8
	)

	for _, v := range []byte(str) {
		switch v {
		case '`':
			continuousBacktick++
			if continuousBacktick == 2 {
				writer.WriteString("``")
				continuousBacktick = 0
			}
		case '.':
			if continuousBacktick > 0 || !selfQuoted {
				shiftDelimiter = 0
				underQuoted = false
				continuousBacktick = 0
				writer.WriteByte('`')
			}
			writer.WriteByte(v)
			continue
		default:
			if shiftDelimiter-continuousBacktick <= 0 && !underQuoted {
				writer.WriteByte('`')
				underQuoted = true
				if selfQuoted = continuousBacktick > 0; selfQuoted {
					continuousBacktick -= 1
				}
			}

			for ; continuousBacktick > 0; continuousBacktick -= 1 {
				writer.WriteString("``")
			}

			writer.WriteByte(v)
		}
		shiftDelimiter++
	}

	if continuousBacktick > 0 && !selfQuoted {
		writer.WriteString("``")
	}
	writer.WriteByte('`')
}

func (dialector Dialector) Explain(sql string, vars ...interface{}) string {
	return logger.ExplainSQL(sql, nil, `'`, vars...)
}

func (dialector Dialector) DataTypeOf(field *schema.Field) string {
	switch field.DataType {
	case schema.Bool:
		return "boolean"
	case schema.Int, schema.Uint:
		return dialector.getSchemaIntAndUnitType(field)
	case schema.Float:
		return dialector.getSchemaFloatType(field)
	case schema.String:
		return dialector.getSchemaStringType(field)
	case schema.Time:
		return dialector.getSchemaTimeType(field)
	case schema.Bytes:
		return dialector.getSchemaBytesType(field)
	}

	return string(field.DataType)
}

func (dialector Dialector) getSchemaFloatType(field *schema.Field) string {
	if field.Precision > 0 {
		return fmt.Sprintf("decimal(%d, %d)", field.Precision, field.Scale)
	}

	if field.Size <= 32 {
		return "float"
	}

	return "double"
}

func (dialector Dialector) getSchemaStringType(field *schema.Field) string {
	size := field.Size
	if size == 0 {
		if dialector.DefaultStringSize > 0 {
			size = int(dialector.DefaultStringSize)
		} else {
			hasIndex := field.TagSettings["INDEX"] != "" || field.TagSettings["UNIQUE"] != ""
			// TEXT, GEOMETRY or JSON column can't have a default value
			if field.PrimaryKey || field.HasDefaultValue || hasIndex {
				size = 191 // utf8mb4
			}
		}
	}

	if size >= 65536 && size <= int(math.Pow(2, 24)) {
		return "mediumtext"
	}

	if size > int(math.Pow(2, 24)) || size <= 0 {
		return "longtext"
	}

	return fmt.Sprintf("varchar(%d)", size)
}

func (dialector Dialector) getSchemaTimeType(field *schema.Field) string {
	precision := ""
	if !dialector.DisableDatetimePrecision && field.Precision == 0 {
		field.Precision = *dialector.DefaultDatetimePrecision
	}

	if field.Precision > 0 {
		precision = fmt.Sprintf("(%d)", field.Precision)
	}

	if field.NotNull || field.PrimaryKey {
		return "datetime" + precision
	}
	return "datetime" + precision + " NULL"
}

func (dialector Dialector) getSchemaBytesType(field *schema.Field) string {
	if field.Size > 0 && field.Size < 65536 {
		return fmt.Sprintf("varbinary(%d)", field.Size)
	}

	if field.Size >= 65536 && field.Size <= int(math.Pow(2, 24)) {
		return "mediumblob"
	}

	return "longblob"
}

func (dialector Dialector) getSchemaIntAndUnitType(field *schema.Field) string {
	sqlType := "bigint"
	switch {
	case field.Size <= 8:
		sqlType = "tinyint"
	case field.Size <= 16:
		sqlType = "smallint"
	case field.Size <= 24:
		sqlType = "mediumint"
	case field.Size <= 32:
		sqlType = "int"
	}

	if field.DataType == schema.Uint {
		sqlType += " unsigned"
	}

	if field.AutoIncrement {
		sqlType += " AUTO_INCREMENT"
	}

	return sqlType
}

func (dialector Dialector) SavePoint(tx *database.DB, name string) error {
	return tx.Exec("SAVEPOINT " + name).Error
}

func (dialector Dialector) RollbackTo(tx *database.DB, name string) error {
	return tx.Exec("ROLLBACK TO SAVEPOINT " + name).Error
}