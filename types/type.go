package types

// Interface types interface
type Interface interface {
	Name() string
	Build(Builder)
	MergeClause(*Type)
}

type TypesBuilder func(Type, Builder)

type Writer interface {
	WriteByte(byte) error
	WriteString(string) (int, error)
}

// Builder builder interface
type Builder interface {
	Writer
	WriteQuoted(field interface{})
	AddVar(Writer, ...interface{})
}

// Type
type Type struct {
	Name                string // WHERE
	BeforeExpression    Expression
	AfterNameExpression Expression
	AfterExpression     Expression
	Expression          Expression
	Builder             TypesBuilder
}

// Build build types
func (t Type) Build(builder Builder) {
	if t.Builder != nil {
		t.Builder(t, builder)
	} else if t.Expression != nil {
		if t.BeforeExpression != nil {
			t.BeforeExpression.Build(builder)
			builder.WriteByte(' ')
		}

		if t.Name != "" {
			builder.WriteString(t.Name)
			builder.WriteByte(' ')
		}

		if t.AfterNameExpression != nil {
			t.AfterNameExpression.Build(builder)
			builder.WriteByte(' ')
		}

		t.Expression.Build(builder)

		if t.AfterExpression != nil {
			builder.WriteByte(' ')
			t.AfterExpression.Build(builder)
		}
	}
}

const (
	PrimaryKey   string = "~~~py~~~" // primary key
	CurrentTable string = "~~~ct~~~" // current table
	Associations string = "~~~as~~~" // associations
)

var (
	currentTable  = Table{Name: CurrentTable}
	PrimaryColumn = Column{Table: CurrentTable, Name: PrimaryKey}
)

// Column quote with name
type Column struct {
	Table string
	Name  string
	Alias string
	Raw   bool
}

// Table quote with name
type Table struct {
	Name  string
	Alias string
	Raw   bool
}