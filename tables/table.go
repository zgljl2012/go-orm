package tables

import (
	"database/sql"
	"fmt"
	"reflect"
	"strings"

	"github.com/zgljl2012/go-orm"
	log "github.com/zgljl2012/slog"
)

const (
	// ErrTableShouldBePointer table should be a pointer of a struct
	ErrTableShouldBePointer = "table should be a pointer of a struct"
	// ErrTableShouldBeStruct table should be struct
	ErrTableShouldBeStruct = "table should be struct"
	// ErrTableNotImplementModelFields implement modelFields
	ErrTableNotImplementModelFields = "table is not implement ModelFields, there are not found function: Fields() []orm.Field "
	// ErrRowIsNotExists exists error
	ErrRowIsNotExists = "The row not exists"
)

type simpleTable struct {
	db     *sql.DB
	fields []orm.Field
	table  interface{}
	name   string
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
		db:     db,
		table:  table,
		fields: table.(orm.ModelFields).Fields(),
		name:   reflect.TypeOf(reflect.Indirect(reflect.ValueOf(table)).Interface()).Name(),
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
	for i, field := range t.fields {
		log.Debug("iterare field", "table", t.Name(), "field", field.Name(), "type", field.Type())
		sql += fmt.Sprintf(`%s %s`, field.Name(), field.Type())
		if i < len(t.fields)-1 {
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
	log.Info(sql)
	if _, err := t.db.Exec(sql); err != nil {
		return err
	}
	return nil
}

func (t *simpleTable) Name() string {
	return t.name
}

func (t *simpleTable) exec(sql string, values []interface{}) error {
	tx, err := t.db.Begin()
	if err != nil {
		return err
	}
	stmt, err := tx.Prepare(sql)
	if err != nil {
		return err
	}
	if _, err := stmt.Exec(values...); err != nil {
		return err
	}
	if err := tx.Commit(); err != nil {
		return err
	}
	if err := stmt.Close(); err != nil {
		return err
	}
	return nil
}

// Add
func (t *simpleTable) Add(instance interface{}) error {
	sql := "INSERT INTO " + t.Name() + " ("
	// fields
	names, values := t.parseInstance(instance, false)
	params := []string{}
	for range names {
		params = append(params, "?")
	}
	sql += strings.Join(names, ",")
	sql += ") VALUES ("
	// values
	sql += strings.Join(params, ",")
	sql += ")"

	log.Debug(sql)

	if err := t.exec(sql, values); err != nil {
		log.Error("got an error when add data", "err", err)
		return err
	}
	return nil
}

// Delete
func (t *simpleTable) Delete(instance interface{}) error {
	// get primary keys
	primaryKeys, primaryValues := t.parseInstance(instance, true)
	for i, key := range primaryKeys {
		primaryKeys[i] = fmt.Sprintf("%s=?", key)
	}
	sql := "DELETE FROM " + t.Name() + " WHERE " + strings.Join(primaryKeys, ",")

	log.Debug(sql)

	if err := t.exec(sql, primaryValues); err != nil {
		log.Error("got an error when delete data", "err", err)
		return err
	}
	return nil
}

func (t *simpleTable) parseInstance(instance interface{}, justPrimaryKeys bool) ([]string, []interface{}) {
	// fields
	names := []string{}
	values := []interface{}{}
	params := []string{}
	_ = values
	for _, field := range t.fields {
		if !justPrimaryKeys || field.PrimaryKey() {
			names = append(names, field.Name())
			value := reflect.ValueOf(instance).Elem().FieldByName(field.Name()).Interface()
			values = append(values, value)
			params = append(params, "?")
		}
	}
	return names, values
}

func (t *simpleTable) Exists(instance interface{}) error {
	cnt, err := t.Count(instance)
	if err != nil {
		return err
	}
	if cnt == 0 {
		return fmt.Errorf(ErrRowIsNotExists)
	}
	return nil
}

func (t *simpleTable) Count(instance interface{}) (int, error) {
	names, values := t.parseInstance(instance, true)
	sql := "SELECT COUNT(*) FROM " + t.Name() + " WHERE "
	for i, name := range names {
		names[i] = fmt.Sprintf("%s=?", name)
	}
	sql += strings.Join(names, " AND ")
	log.Debug(sql)
	tx, err := t.db.Begin()
	if err != nil {
		return 0, err
	}
	stmt, err := tx.Prepare(sql)
	if err != nil {
		return 0, err
	}
	cnt := 0
	if err := stmt.QueryRow(values...).Scan(&cnt); err != nil {
		return 0, err
	}
	if err := tx.Commit(); err != nil {
		return 0, err
	}
	if err := stmt.Close(); err != nil {
		return 0, err
	}
	return cnt, nil
}

// Update
func (t *simpleTable) Update(instance interface{}) error {
	// check the row exists or not
	if err := t.Exists(instance); err != nil {
		return err
	}
	// get primary keys
	primaryKeys, primaryValues := t.parseInstance(instance, true)
	for i, key := range primaryKeys {
		primaryKeys[i] = fmt.Sprintf("%s=?", key)
	}

	// keys, values
	names, values := t.parseInstance(instance, false)

	for i, name := range names {
		names[i] = fmt.Sprintf("%s=?", name)
	}

	// sql
	sql := "UPDATE " + t.Name() + " SET "
	sql += strings.Join(names, ",")
	sql += " WHERE " + strings.Join(primaryKeys, " AND ")

	log.Debug(sql)

	if err := t.exec(sql, append(values, primaryValues...)); err != nil {
		log.Error("got an error when update data", "err", err)
		return err
	}
	return nil
}

func (t *simpleTable) Filter(filters ...*orm.QueryParameter) orm.FilterSet {
	// validate parameters
	return newFilterSet(t.db, t.Name(), t.table).Filter(filters...)
}

func (t *simpleTable) Upsert(instance interface{}) error {
	// check the row exists or not
	if err := t.Exists(instance); err == nil {
		return t.Update(instance)
	}
	return t.Add(instance)
}
