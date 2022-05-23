package types

type Insert struct {
	Table    Table
	Modifier string
}

// Name insert types name
func (insert Insert) Name() string {
	return "INSERT"
}

// Build build insert types
func (insert Insert) Build(builder Builder) {
	if insert.Modifier != "" {
		builder.WriteString(insert.Modifier)
		builder.WriteByte(' ')
	}

	builder.WriteString("INTO ")
	if insert.Table.Name == "" {
		builder.WriteQuoted(currentTable)
	} else {
		builder.WriteQuoted(insert.Table)
	}
}

// MergeClause merge insert types
func (insert Insert) MergeClause(types *Type) {
	if v, ok := types.Expression.(Insert); ok {
		if insert.Modifier == "" {
			insert.Modifier = v.Modifier
		}
		if insert.Table.Name == "" {
			insert.Table = v.Table
		}
	}
	types.Expression = insert
}