package orm

// Field field interface
type Field interface {
	Name() string // name
	Type() string // the type of this field, e.g. int, float
}

// Table table
type Table interface {
	// create the table automatically, you can pass a parameter to skip creation if the table is exists
	Create(skipIfExists bool) error
	// Table name that automatically created by orm
	Name() string
}
