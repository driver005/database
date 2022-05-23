package tests

import (
	"github.com/driver005/database"
	"github.com/driver005/database/logger"
	"github.com/driver005/database/schema"
	"github.com/driver005/database/types"
)

type DummyDialector struct{}

func (DummyDialector) Name() string {
	return "dummy"
}

func (DummyDialector) Initialize(*database.DB) error {
	return nil
}

func (DummyDialector) DefaultValueOf(field *schema.Field) types.Expression {
	return types.Expr{SQL: "DEFAULT"}
}

func (DummyDialector) Migrator(*database.DB) database.Migrator {
	return nil
}

func (DummyDialector) BindVarTo(writer types.Writer, stmt *database.Statement, v interface{}) {
	writer.WriteByte('?')
}

func (DummyDialector) QuoteTo(writer types.Writer, str string) {
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

func (DummyDialector) Explain(sql string, vars ...interface{}) string {
	return logger.ExplainSQL(sql, nil, `"`, vars...)
}

func (DummyDialector) DataTypeOf(*schema.Field) string {
	return ""
}
