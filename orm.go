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
	// Upsert add or update
	Upsert(instance interface{}) error
	// Delete operate will delete via primary keys
	Delete(instance interface{}) error
	// Update operate will select those row via primary keys, then update other fields.
	// So your should be sure of your primary keys won't be updated.
	Update(instance interface{}) error
	// Filter rows
	Filter(...*QueryParameter) FilterSet
	// Count get the counts
	Count(instance interface{}) (int, error)
}

// QueryParameter for filter
type QueryParameter struct {
	Name     string
	Value    interface{}
	Operator string // 操作符
}

// WithParameter create paramter pair
func WithParameter(name string, value interface{}) *QueryParameter {
	return &QueryParameter{
		Name:     name,
		Value:    value,
		Operator: "=",
	}
}

// FilterSet for select
// you can iterate FilterSet via range
type FilterSet interface {
	// Filter with paramters
	Filter(parameters ...*QueryParameter) FilterSet
	// OrderBy specify ordering fields, plus means ASC, minus(-) means DESC
	OrderBy(...string) FilterSet
	// Limit rows
	Limit(int) FilterSet
	// Offset set offset
	Offset(int) FilterSet
	// All return all rows, returned data just an array of objects, not pointer.
	All() []interface{}
}
