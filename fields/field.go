package fields

import (
	"fmt"
	"reflect"
	"strconv"

	"github.com/zgljl2012/go-orm"
	"github.com/zgljl2012/slog"
	log "github.com/zgljl2012/slog"
)

// Field int field
type myField struct {
	id      string
	name    string
	_type   Type
	options *FieldOptions
}

func newFiled(id string, name string, _type Type, opts ...FieldOption) orm.Field {
	options := defaultOptions

	for _, o := range opts {
		o(&options)
	}

	return &myField{
		id:      id,
		name:    name,
		_type:   _type,
		options: &options,
	}
}

type valueValidator func(string) error

// parseFieldOptions parse options
func parseFieldOptions(field reflect.StructField) ([]FieldOption, error) {
	boolValidator := func(value string) error {
		if value != "true" && value != "false" {
			return fmt.Errorf(`can't support "%v" for bool field`, value)
		}
		return nil
	}
	tags := []struct {
		tag        string                         // tag
		_type      reflect.Kind                   // type
		fun        func(value string) FieldOption // handler
		validators []valueValidator
	}{
		{
			tag:   "primaryKey",
			_type: reflect.Bool,
			validators: []valueValidator{
				boolValidator,
			},
			fun: func(value string) FieldOption {
				return func(options *FieldOptions) {
					options.PrimaryKey = value == "true"
				}
			},
		},
		{
			tag:   "null",
			_type: reflect.Bool,
			fun: func(value string) FieldOption {
				return func(options *FieldOptions) {
					options.Null = value == "true"
				}
			},
		},
		{
			tag:   "length",
			_type: reflect.Int,
			fun: func(value string) FieldOption {
				return func(options *FieldOptions) {
					val, _ := strconv.Atoi(value)
					options.Length = val
				}
			},
			validators: []valueValidator{
				func(value string) error {
					_, err := strconv.Atoi(value)
					if err != nil {
						return fmt.Errorf(`parse length tag error, field: "%s", length: "%s", err: "%s"`,
							field.Name, value, err)
					}
					return nil
				},
			},
		},
	}
	options := []FieldOption{}
	for _, tag := range tags {
		value := field.Tag.Get(tag.tag)
		if value != "" {
			if tag.validators != nil {
				for _, validator := range tag.validators {
					if err := validator(value); err != nil {
						return nil, err
					}
				}
			}
			options = append(options, tag.fun(value))
		}
	}
	return options, nil
}

// ParseStructWithTagsToFields parse the struct's fields with tags to orm.field
func ParseStructWithTagsToFields(instance interface{}) ([]orm.Field, error) {
	// TODO: 校验 name 命名的合法性；校验 name 是否重复
	results := []orm.Field{}
	// iterate fields of instance
	value := reflect.Indirect(reflect.ValueOf(instance))
	t := value.Type()
	log.Info(t)
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		// type
		kind := field.Type.Kind()
		// name
		name := field.Tag.Get("name")
		if name != "" {
			slog.Info("iterate field", "kind", kind, "name", name)
			options, err := parseFieldOptions(field)
			if err != nil {
				return nil, err
			}
			// TODO: 如果 Field 的类型是 String，但不包含 length tag 就报错
			var f orm.Field
			if kind == reflect.Int {
				f = newFiled(field.Name, name, INT, options...)
			} else if kind == reflect.Uint64 {
				f = newFiled(field.Name, name, UINT64, options...)
			} else if kind == reflect.String {
				f = newFiled(field.Name, name, CHAR, options...)
			} else if kind == reflect.Bool {
				f = newFiled(field.Name, name, BOOL, options...)
			} else if kind == reflect.Float32 {
				f = newFiled(field.Name, name, FLOAT, options...)
			} else if kind == reflect.Struct && field.Type.String() == "time.Time" {
				f = newFiled(field.Name, name, DATETIME, options...)
			} else {
				return nil, fmt.Errorf(`Unsupport type "%s - %s" of field "%s"`, field.Type, kind, field.Name)
			}
			results = append(results, f)
		}
	}
	return results, nil
}

// NewIntField new an int field
func NewIntField(name string, opts ...FieldOption) orm.Field {
	return newFiled(name, name, INT, opts...)
}

// NewUInt64Field new an uint64 field
func NewUInt64Field(name string, opts ...FieldOption) orm.Field {
	return newFiled(name, name, UINT64, opts...)
}

// NewBoolField new a bool field
func NewBoolField(name string, opts ...FieldOption) orm.Field {
	return newFiled(name, name, BOOL, opts...)
}

// NewFloatField new a bool field
func NewFloatField(name string, opts ...FieldOption) orm.Field {
	return newFiled(name, name, FLOAT, opts...)
}

// NewDatetimeField new a datetime field
func NewDatetimeField(name string, opts ...FieldOption) orm.Field {
	return newFiled(name, name, DATETIME, opts...)
}

// NewCharField new a char field
// you can set the length with WithLength, the default length is 100.
func NewCharField(name string, opts ...FieldOption) orm.Field {
	return newFiled(name, name, CHAR, opts...)
}

// Type return type
func (f *myField) Type() string {
	var t string

	if f._type == CHAR {
		t = f._type.String() + fmt.Sprintf("(%d)", f.options.Length)
	} else {
		t = f._type.String()
	}
	if !f.options.Null || f.PrimaryKey() {
		t += " NOT NULL"
	} else if !f.PrimaryKey() {
		t += " NULL"
	}
	return t
}

func (f *myField) Name() string {
	return f.name
}

func (f *myField) ID() string {
	return f.id
}

func (f *myField) PrimaryKey() bool {
	return f.options.PrimaryKey
}
