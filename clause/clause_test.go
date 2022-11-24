package clause_test

import (
	"reflect"
	"strings"
	"sync"
	"testing"

	"github.com/driver005/database"
	"github.com/driver005/database/clause"
	"github.com/driver005/database/schema"
	"github.com/driver005/database/utils/tests"
)

var db, _ = database.Open(tests.DummyDialector{}, nil)

func checkBuildClauses(t *testing.T, clauses []clause.Interface, result string, vars []interface{}) {
	var (
		buildNames    []string
		buildNamesMap = map[string]bool{}
		user, _       = schema.Parse(&tests.User{}, &sync.Map{}, db.NamingStrategy)
		stmt          = database.Statement{DB: db, Table: user.Table, Schema: user, Clauses: map[string]clause.Clause{}}
	)

	for _, c := range clauses {
		if _, ok := buildNamesMap[c.Name()]; !ok {
			buildNames = append(buildNames, c.Name())
			buildNamesMap[c.Name()] = true
		}

		stmt.AddClause(c)
	}

	stmt.Build(buildNames...)

	if strings.TrimSpace(stmt.SQL.String()) != result {
		t.Errorf("SQL expects %v got %v", result, stmt.SQL.String())
	}

	if !reflect.DeepEqual(stmt.Vars, vars) {
		t.Errorf("Vars expects %+v got %v", stmt.Vars, vars)
	}
}
