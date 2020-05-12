package fields

import (
	"fmt"

	"github.com/zgljl2012/go-orm"
)

// Field int field
type myField struct {
	name    string
	_type   Type
	options *FieldOptions
}

func newFiled(name string, _type Type, opts ...FieldOption) orm.Field {
	options := defaultOptions

	for _, o := range opts {
		o(&options)
	}

	return &myField{
		name:    name,
		_type:   _type,
		options: &options,
	}
}

// NewIntField new an int field
func NewIntField(name string, opts ...FieldOption) orm.Field {
	return newFiled(name, INT, opts...)
}

// NewUInt64Field new an uint64 field
func NewUInt64Field(name string, opts ...FieldOption) orm.Field {
	return newFiled(name, UINT64, opts...)
}

// NewBoolField new a bool field
func NewBoolField(name string, opts ...FieldOption) orm.Field {
	return newFiled(name, BOOL, opts...)
}

// NewFloatField new a bool field
func NewFloatField(name string, opts ...FieldOption) orm.Field {
	return newFiled(name, FLOAT, opts...)
}

// NewDatetimeField new a datetime field
func NewDatetimeField(name string, opts ...FieldOption) orm.Field {
	return newFiled(name, DATETIME, opts...)
}

// NewCharField new a char field
// you can set the length with WithLength, the default length is 100.
func NewCharField(name string, opts ...FieldOption) orm.Field {
	return newFiled(name, CHAR, opts...)
}

// Type return type
func (f *myField) Type() string {
	var t string

	if f._type == CHAR {
		t = f._type.String() + fmt.Sprintf("(%d)", f.options.Length)
	} else {
		t = f._type.String()
	}
	if !f.options.Null {
		t += " NOT NULL"
	} else {
		t += " NULL"
	}
	return t
}

func (f *myField) Name() string {
	return f.name
}

func (f *myField) PrimaryKey() bool {
	return f.options.PrimaryKey
}
