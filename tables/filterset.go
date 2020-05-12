package tables

import (
	"database/sql"
	"fmt"
	"reflect"
	"strings"

	"github.com/zgljl2012/go-orm"
	log "github.com/zgljl2012/slog"
)

type filterSet struct {
	instance   interface{}
	table      string
	offset     int
	limit      int
	parameters []*orm.QueryParameter
	order      []string
	db         *sql.DB
}

func newFilterSet(db *sql.DB, table string, instance interface{}) orm.FilterSet {
	return &filterSet{
		instance:   instance,
		db:         db,
		table:      table,
		limit:      0,
		offset:     0,
		parameters: []*orm.QueryParameter{},
	}
}

// Filter with paramters
func (f *filterSet) Filter(parameters ...*orm.QueryParameter) orm.FilterSet {
	f.parameters = append(f.parameters, parameters...)
	return f
}

// OrderBy specify ordering fields, plus means ASC, minus(-) means DESC
func (f *filterSet) OrderBy(orders ...string) orm.FilterSet {
	f.order = append(f.order, orders...)
	return f
}

// Limit rows
func (f *filterSet) Limit(limit int) orm.FilterSet {
	if limit > 0 {
		f.limit = limit
	}
	return f
}

// All return all rows
func (f *filterSet) All() []interface{} {
	var (
		sql    string
		names  []string
		values []interface{}
	)
	// filter
	sql = "SELECT * FROM " + f.table
	if len(f.parameters) > 0 {
		sql += " WHERE "
		for _, parameter := range f.parameters {
			names = append(names, parameter.Name+" "+parameter.Operator+" ?")
			values = append(values, parameter.Value)
		}
		sql += strings.Join(names, ",")
	}
	// order
	var orders []string
	for _, order := range f.order {
		order = strings.Trim(order, " ")
		if order[0] == '-' {
			orders = append(orders, order[1:]+" DESC")
		} else {
			orders = append(orders, order)
		}
	}
	if len(orders) > 0 {
		sql += " ORDER BY " + strings.Join(orders, ",")
	}
	// limit
	if f.limit > 0 {
		sql += fmt.Sprintf(" LIMIT %v", f.limit)
	}
	// offset
	if f.offset > 0 {
		if f.limit <= 0 {
			// default limit
			sql += fmt.Sprintf(" LIMIT %v", 10)
		}
		sql += fmt.Sprintf(" OFFSET %v", f.offset)
	}
	// query
	log.Debug(sql)
	tx, err := f.db.Begin()
	if err != nil {
		log.Fatal("get tx error when iterate all rows", "err", err)
	}
	stmt, err := tx.Prepare(sql)
	if err != nil {
		log.Fatal("stmt error when iterate all rows", "err", err)
	}
	result := []interface{}{}
	if rows, err := stmt.Query(values...); err != nil {
		log.Error("iterate data error", "err", err)
	} else {
		for rows.Next() {
			// new instance
			obj := reflect.New(reflect.TypeOf(f.instance).Elem()).Elem()
			numCols := reflect.TypeOf(f.instance).Elem().NumField()
			columns := make([]interface{}, numCols)
			for i := 0; i < numCols; i++ {
				field := obj.Field(i)
				columns[i] = field.Addr().Interface()
			}
			if err := rows.Scan(columns...); err != nil {
				log.Info("scan row error", "err", err)
			}
			result = append(result, obj.Interface())
		}
		if err := rows.Close(); err != nil {
			log.Error("got an error when close rows", "err", err)
		}
	}
	if err := tx.Commit(); err != nil {
		log.Error(err)
	}
	if err := stmt.Close(); err != nil {
		log.Error(err)
	}
	return result
}

func (f *filterSet) Offset(offset int) orm.FilterSet {
	if offset > 0 {
		f.offset = offset
	}
	return f
}
