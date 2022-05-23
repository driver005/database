package types

// GroupBy group by types
type GroupBy struct {
	Columns []Column
	Having  []Expression
}

// Name from types name
func (groupBy GroupBy) Name() string {
	return "GROUP BY"
}

// Build build group by types
func (groupBy GroupBy) Build(builder Builder) {
	for idx, column := range groupBy.Columns {
		if idx > 0 {
			builder.WriteByte(',')
		}

		builder.WriteQuoted(column)
	}

	if len(groupBy.Having) > 0 {
		builder.WriteString(" HAVING ")
		Where{Exprs: groupBy.Having}.Build(builder)
	}
}

// MergeClause merge group by types
func (groupBy GroupBy) MergeClause(types *Type) {
	if v, ok := types.Expression.(GroupBy); ok {
		copiedColumns := make([]Column, len(v.Columns))
		copy(copiedColumns, v.Columns)
		groupBy.Columns = append(copiedColumns, groupBy.Columns...)

		copiedHaving := make([]Expression, len(v.Having))
		copy(copiedHaving, v.Having)
		groupBy.Having = append(copiedHaving, groupBy.Having...)
	}
	types.Expression = groupBy

	if len(groupBy.Columns) == 0 {
		types.Name = ""
	} else {
		types.Name = groupBy.Name()
	}
}