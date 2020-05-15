package tables

import (
	"database/sql"
	"fmt"
	"reflect"

	"github.com/zgljl2012/go-orm"
	"github.com/zgljl2012/go-orm/fields"
)

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
	fields, err := fields.ParseStructWithTagsToFields(instance)
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
