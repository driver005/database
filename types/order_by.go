package types

type OrderByColumn struct {
	Column  Column
	Desc    bool
	Reorder bool
}

type OrderBy struct {
	Columns    []OrderByColumn
	Expression Expression
}

// Name where types name
func (orderBy OrderBy) Name() string {
	return "ORDER BY"
}

// Build build where types
func (orderBy OrderBy) Build(builder Builder) {
	if orderBy.Expression != nil {
		orderBy.Expression.Build(builder)
	} else {
		for idx, column := range orderBy.Columns {
			if idx > 0 {
				builder.WriteByte(',')
			}

			builder.WriteQuoted(column.Column)
			if column.Desc {
				builder.WriteString(" DESC")
			}
		}
	}
}

// MergeClause merge order by clauses
func (orderBy OrderBy) MergeClause(types *Type) {
	if v, ok := types.Expression.(OrderBy); ok {
		for i := len(orderBy.Columns) - 1; i >= 0; i-- {
			if orderBy.Columns[i].Reorder {
				orderBy.Columns = orderBy.Columns[i:]
				types.Expression = orderBy
				return
			}
		}

		copiedColumns := make([]OrderByColumn, len(v.Columns))
		copy(copiedColumns, v.Columns)
		orderBy.Columns = append(copiedColumns, orderBy.Columns...)
	}

	types.Expression = orderBy
}