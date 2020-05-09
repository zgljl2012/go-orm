package fields

import "orm"

// IntField int field
type intField struct {
	name    string
	options *FieldOptions
}

// NewIntField new an int field
func NewIntField(name string, opts ...FieldOption) orm.Field {
	options := defaultOptions

	for _, o := range opts {
		o(&options)
	}

	return &intField{
		name:    name,
		options: &options,
	}
}

// Type return type
func (f *intField) Type() string {
	return "INT"
}

func (f *intField) Name() string {
	return f.name
}
