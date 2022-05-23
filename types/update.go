package types

type Update struct {
	Modifier string
	Table    Table
}

// Name update types name
func (update Update) Name() string {
	return "UPDATE"
}

// Build build update types
func (update Update) Build(builder Builder) {
	if update.Modifier != "" {
		builder.WriteString(update.Modifier)
		builder.WriteByte(' ')
	}

	if update.Table.Name == "" {
		builder.WriteQuoted(currentTable)
	} else {
		builder.WriteQuoted(update.Table)
	}
}

// MergeClause merge update types
func (update Update) MergeClause(types *Type) {
	if v, ok := types.Expression.(Update); ok {
		if update.Modifier == "" {
			update.Modifier = v.Modifier
		}
		if update.Table.Name == "" {
			update.Table = v.Table
		}
	}
	types.Expression = update
}