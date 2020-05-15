package tables

import (
	"database/sql"
	"fmt"
	"reflect"
	"strconv"

	"github.com/zgljl2012/go-orm"
	"github.com/zgljl2012/go-orm/fields"
	log "github.com/zgljl2012/slog"
)

func parseWithParameters(field reflect.StructField) ([]fields.FieldOption, error) {
	tags := []struct {
		tag   string       // tag
		_type reflect.Kind // type
	}{
		{
			tag:   "primaryKey",
			_type: reflect.Bool,
		},
		{
			tag:   "primaryKey",
			_type: reflect.Bool,
		},
		{
			tag:   "length",
			_type: reflect.Int,
		},
	}
	options := []fields.FieldOption{}
	for _, tag := range tags {
		value := field.Tag.Get(tag.tag)
		if value != "" {
			if tag._type == reflect.Bool && value == "true" && tag.tag == "primaryKey" {
				options = append(options, fields.WithPrimaryKey(true))
			} else if tag._type == reflect.Int && tag.tag == "length" {
				length, err := strconv.Atoi(value)
				if err != nil {
					return nil, fmt.Errorf(`parse length tag error, field: "%s", length: "%s", err: "%s"`,
						field.Name, value, err)
				}
				options = append(options, fields.WithLength(length))
			} else if tag._type == reflect.Bool && tag.tag == "null" {
				if value == "true" {
					options = append(options, fields.WithNull(true))
				} else if value == "false" {
					options = append(options, fields.WithNull(true))
				}
			}
		}
	}
	return options, nil
}

func parseStructTags(instance interface{}) ([]orm.Field, error) {
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
			log.Info("iterate field", "kind", kind, "name", name)
			options, err := parseWithParameters(field)
			if err != nil {
				return nil, err
			}
			// TODO: 如果 Field 的类型是 String，但不包含 length tag 就报错
			if kind == reflect.Int {
				f := fields.NewIntField(name, options...)
				results = append(results, f)
			} else if kind == reflect.String {
				f := fields.NewCharField(name, options...)
				results = append(results, f)
			} else {
				return nil, fmt.Errorf(`Unsupport type "%s" of field "%s"`, kind, field.Name)
			}
		}
	}
	return results, nil
}

// NewStructTagsTable new a table with tags
func NewStructTagsTable(db *sql.DB, instance interface{}) (orm.Table, error) {
	t := reflect.TypeOf(instance)
	kind := t.Kind()
	if kind != reflect.Ptr {
		return nil, fmt.Errorf(ErrTableShouldBePointer)
	}
	if reflect.Indirect(reflect.ValueOf(instance)).Kind() != reflect.Struct {
		return nil, fmt.Errorf(ErrTableShouldBePointer)
	}
	fields, err := parseStructTags(instance)
	if err != nil {
		return nil, err
	}

	if len(fields) == 0 {
		return nil, fmt.Errorf("There are no fields in your instance")
	}

	return &simpleTable{
		fields: fields,
		db:     db,
		table:  instance,
		name:   reflect.TypeOf(reflect.Indirect(reflect.ValueOf(instance)).Interface()).Name(),
	}, nil
}
