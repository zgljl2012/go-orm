package tables

import (
	"database/sql"
	"fmt"
	"orm"
	"reflect"
	"strings"

	log "github.com/zgljl2012/slog"
)

const (
	// ErrTableShouldBePointer table should be a pointer of a struct
	ErrTableShouldBePointer = "table should be a pointer of a struct"
	// ErrTableShouldBeStruct table should be struct
	ErrTableShouldBeStruct = "table should be struct"
	// ErrTableNotImplementModelFields implement modelFields
	ErrTableNotImplementModelFields = "table is not implement ModelFields, there are not found function: Fields() []orm.Field "
)

type simpleTable struct {
	db    *sql.DB
	table interface{}
}

// NewTable create a table instance, you can input every struct.
// All pub fields will be checked if their type is orm.Field.
func NewTable(db *sql.DB, table interface{}) (orm.Table, error) {
	t := reflect.TypeOf(table)
	kind := t.Kind()
	log.Debug("table", "type", t, "kind", kind, "ptrTo", reflect.Indirect(reflect.ValueOf(table)).Kind())
	if kind != reflect.Ptr {
		return nil, fmt.Errorf(ErrTableShouldBePointer)
	}
	if reflect.Indirect(reflect.ValueOf(table)).Kind() != reflect.Struct {
		return nil, fmt.Errorf(ErrTableShouldBePointer)
	}
	// check if implement ModelField
	if !t.Implements(reflect.TypeOf((*orm.ModelFields)(nil)).Elem()) {
		return nil, fmt.Errorf(ErrTableNotImplementModelFields)
	}
	return &simpleTable{
		db:    db,
		table: table,
	}, nil
}

func (t *simpleTable) Create(skipIfExists bool) error {
	var primaryKeys []string
	sql := "CREATE TABLE "
	if skipIfExists {
		sql += " IF NOT EXISTS "
	}
	sql += t.Name()
	sql += `(`
	// iterate fields
	fields := t.table.(orm.ModelFields).Fields()
	for i, field := range fields {
		log.Debug("iterare field", "table", t.Name(), "field", field.Name(), "type", field.Type())
		sql += fmt.Sprintf(`%s %s`, field.Name(), field.Type())
		if i < len(fields)-1 {
			sql += ","
		}
		if field.PrimaryKey() {
			primaryKeys = append(primaryKeys, field.Name())
		}
	}
	// primary keys
	if len(primaryKeys) > 0 {
		sql += ", PRIMARY KEY("
		sql += strings.Join(primaryKeys, ",")
		sql += ")"
	}
	sql += `)`
	log.Debug(sql)
	if _, err := t.db.Exec(sql); err != nil {
		return err
	}
	return nil
}

func (t *simpleTable) Name() string {
	name := reflect.TypeOf(reflect.Indirect(reflect.ValueOf(t.table)).Interface()).Name()
	return name
}
