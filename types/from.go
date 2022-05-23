package types

// From from types
type From struct {
	Tables []Table
	Joins  []Join
}

// Name from types name
func (from From) Name() string {
	return "FROM"
}

// Build build from types
func (from From) Build(builder Builder) {
	if len(from.Tables) > 0 {
		for idx, table := range from.Tables {
			if idx > 0 {
				builder.WriteByte(',')
			}

			builder.WriteQuoted(table)
		}
	} else {
		builder.WriteQuoted(currentTable)
	}

	for _, join := range from.Joins {
		builder.WriteByte(' ')
		join.Build(builder)
	}
}

// MergeClause merge from types
func (from From) MergeClause(types *Type) {
	types.Expression = from
}