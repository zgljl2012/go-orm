package fields

import "orm"

// Field int field
type myField struct {
	name    string
	_type   string
	options *FieldOptions
}

// NewIntField new an int field
func NewIntField(name string, opts ...FieldOption) orm.Field {
	options := defaultOptions

	for _, o := range opts {
		o(&options)
	}

	return &myField{
		name:    name,
		_type:   "INT",
		options: &options,
	}
}

// Type return type
func (f *myField) Type() string {
	return f._type
}

func (f *myField) Name() string {
	return f.name
}
