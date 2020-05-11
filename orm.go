package orm

// Field field interface
type Field interface {
	Name() string     // name
	Type() string     // the type of this field, e.g. int, float
	PrimaryKey() bool // primary key
}

// Table table
type Table interface {
	// create the table automatically, you can pass a parameter to skip creation if the table is exists
	Create(skipIfExists bool) error
	// Table name that automatically created by orm
	Name() string
	// Add
	Add(instance interface{}) error
	// Delete
	Delete(instance interface{}) error
	// Update operate will select those row via primary keys, then update other fields.
	// So your should be sure of your primary keys won't be updated.
	Update(instance interface{}) error
	// Query
	// Query() error
}
