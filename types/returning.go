package types

type Returning struct {
	Columns []Column
}

// Name where types name
func (returning Returning) Name() string {
	return "RETURNING"
}

// Build build where types
func (returning Returning) Build(builder Builder) {
	if len(returning.Columns) > 0 {
		for idx, column := range returning.Columns {
			if idx > 0 {
				builder.WriteByte(',')
			}

			builder.WriteQuoted(column)
		}
	} else {
		builder.WriteByte('*')
	}
}

// MergeClause merge order by clauses
func (returning Returning) MergeClause(types *Type) {
	if v, ok := types.Expression.(Returning); ok {
		returning.Columns = append(v.Columns, returning.Columns...)
	}

	types.Expression = returning
}