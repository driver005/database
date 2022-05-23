package types

import "strconv"

// Limit limit types
type Limit struct {
	Limit  int
	Offset int
}

// Name where types name
func (limit Limit) Name() string {
	return "LIMIT"
}

// Build build where types
func (limit Limit) Build(builder Builder) {
	if limit.Limit > 0 {
		builder.WriteString("LIMIT ")
		builder.WriteString(strconv.Itoa(limit.Limit))
	}
	if limit.Offset > 0 {
		if limit.Limit > 0 {
			builder.WriteByte(' ')
		}
		builder.WriteString("OFFSET ")
		builder.WriteString(strconv.Itoa(limit.Offset))
	}
}

// MergeClause merge order by clauses
func (limit Limit) MergeClause(types *Type) {
	types.Name = ""

	if v, ok := types.Expression.(Limit); ok {
		if limit.Limit == 0 && v.Limit != 0 {
			limit.Limit = v.Limit
		}

		if limit.Offset == 0 && v.Offset > 0 {
			limit.Offset = v.Offset
		} else if limit.Offset < 0 {
			limit.Offset = 0
		}
	}

	types.Expression = limit
}