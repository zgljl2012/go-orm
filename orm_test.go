package orm_test

import (
	"database/sql"
	"fmt"
	"orm"
	"orm/fields"
	"orm/tables"
	"os"
	"testing"

	_ "github.com/mattn/go-sqlite3"
	log "github.com/zgljl2012/slog"
)

var (
	testDB = "./test.db"
)

// User is a test table
type User struct {
	ID       int
	Username string
	Password string
}

// Fields return all fields to want to bind with database
func (u *User) Fields() []orm.Field {
	return []orm.Field{
		fields.NewIntField("ID", fields.WithPrimaryKey(true)),
	}
}

func createTestDatabase() *sql.DB {
	db, err := sql.Open("sqlite3", testDB)
	if err != nil {
		log.Fatal(err)
	}
	return db
}

func deleteTestDatabase() {
	if err := os.Remove(testDB); err != nil {
		log.Fatal(err)
	}
}

func TestCreateTable(t *testing.T) {
	db := createTestDatabase()
	defer deleteTestDatabase()

	// create user table instance
	table, err := tables.NewTable(db, &User{})
	if err != nil {
		t.Fatal(err)
	}

	// table's type is wrong
	if _, err := tables.NewTable(db, 1); err == nil {
		t.Fatal("should got an error, but is normal")
	}

	// table's is not implements ModelFields
	if _, err := tables.NewTable(db, &struct{}{}); err == nil {
		t.Fatal("should got an error, but is normal")
	}

	// create table in database, name is the same as struct
	if err := table.Create(false); err != nil {
		t.Error(err)
	}

	// Check if the table has been created
	if _, err := db.Query(fmt.Sprintf("SELECT COUNT(*) FROM %s", table.Name())); err != nil {
		t.Fatal(err)
	}

	// If you create again, you will get an error because the table already exists
	if err := table.Create(false); err == nil {
		t.Error("you should get an error because the table already exists")
	}

	// But if you skip creation, you won't get the error above.
	if err := table.Create(true); err != nil {
		t.Error(err)
	}

	// check primary key
	if result, err := db.Query(fmt.Sprintf("PRAGMA table_info(%s)", table.Name())); err != nil {
		t.Fatal(err)
	} else {
		if cols, err := result.Columns(); err != nil {
			t.Error(err)
		} else {
			for _, col := range cols {
				// t.Log(col)
				_ = col
			}
		}

		fields := map[string]map[string]interface{}{
			"ID": {
				"exists": false,
				"type":   "INT",
				"pk":     true,
			},
			"Username": {
				"exists": false,
				"type":   "String",
				"pk":     false,
			},
			"Password": {
				"exists": false,
				"type":   "String",
				"pk":     false,
			},
		}

		for result.Next() {
			var (
				cid        int
				name       string
				_type      string
				notnull    bool
				dflt_value interface{}
				pk         bool
			)
			if err := result.Scan(&cid, &name, &_type, &notnull, &dflt_value, &pk); err != nil {
				t.Error(err)
			}
			t.Log(cid, name, _type, notnull, dflt_value, pk)
			// validate field
			if field, ok := fields[name]; ok {
				field["exists"] = true
				if field["type"].(string) != _type {
					t.Errorf("Field %v's type is wrong, expect %v, but got %v", name, field["type"], _type)
				}
				if field["pk"].(bool) != pk {
					t.Errorf("Field %v's pk is wrong, expect %v, but got %v", name, field["pk"], pk)
				}
			} else {
				t.Errorf("There is a undefined field, name:%v, type:%v, pk:%v", name, _type, pk)
			}
		}

		// iterate fields
		for name, field := range fields {
			if !field["exists"].(bool) {
				t.Errorf("field %v is not found", name)
			}
		}
	}

}
