package tables_test

import (
	"database/sql"
	"fmt"
	"os"
	"testing"

	_ "github.com/mattn/go-sqlite3"
	"github.com/zgljl2012/go-orm"
	"github.com/zgljl2012/go-orm/tables"
	log "github.com/zgljl2012/slog"
)

var (
	testDB = "./test.db"
)

func createTestDatabase() *sql.DB {
	db, err := sql.Open("sqlite3", testDB)
	if err != nil {
		log.Fatal(err)
	}
	return db
}

func deleteTestDatabase() {
	if _, err := os.Stat(testDB); err == nil {
		if err := os.Remove(testDB); err != nil {
			log.Fatal(err)
		}
	}
}

// User is a test table
type User struct {
	ID       int    `name:"id" primaryKey:"true"`
	Username string `name:"username" length:"20"`
	Password string `name:"password" length:"50"`
	Active   bool
	// Age       float32
	// CreatedAt time.Time
	// Count     uint64
}

func TestCreateTable(t *testing.T) {
	db := createTestDatabase()
	defer deleteTestDatabase()

	// create user table instance
	table, err := tables.NewStructTagsTable(db, &User{})
	if err != nil {
		t.Fatal(err)
	}

	// table's type is wrong
	if _, err := tables.NewStructTagsTable(db, 1); err == nil {
		t.Fatal("should got an error, but is normal")
	}

	// table's do not have any fields with name tag
	if _, err := tables.NewStructTagsTable(db, &struct{}{}); err == nil {
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
			"id": {
				"exists": false,
				"type":   "INT",
				"pk":     true,
				"null":   false,
			},
			"username": {
				"exists": false,
				"type":   "CHAR(20)",
				"pk":     false,
				"null":   true,
			},
			"password": {
				"exists": false,
				"type":   "CHAR(50)",
				"pk":     false,
				"null":   true,
			},
			// "Active": {
			// 	"exists": false,
			// 	"type":   "BOOL",
			// 	"pk":     false,
			// 	"null":   false,
			// },
			// "CreatedAt": {
			// 	"exists": false,
			// 	"type":   "DATETIME",
			// 	"pk":     false,
			// 	"null":   true,
			// },
			// "Age": {
			// 	"exists": false,
			// 	"type":   "FLOAT",
			// 	"pk":     false,
			// 	"null":   true,
			// },
			// "Count": {
			// 	"exists": false,
			// 	"type":   "BIGINT",
			// 	"pk":     false,
			// 	"null":   true,
			// },
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
			// t.Log(cid, name, _type, notnull, dflt_value, pk, notnull)
			// validate field
			if field, ok := fields[name]; ok {
				field["exists"] = true
				if field["type"].(string) != _type {
					t.Errorf("Field %v's type is wrong, expect %v, but got %v", name, field["type"], _type)
				}
				if field["pk"].(bool) != pk {
					t.Errorf("Field %v's pk is wrong, expect %v, but got %v", name, field["pk"], pk)
				}
				if field["null"].(bool) != !notnull {
					t.Errorf("Field %v's null is wrong, expect %v, but got %v", name, field["null"], !notnull)
				}
			} else {
				t.Errorf("There is a undefined field, name:%v, type:%v, pk:%v", name, _type, pk)
			}
		}

		// iterate fields
		for name, field := range fields {
			if !field["exists"].(bool) {
				t.Errorf("field %v not found", name)
			}
		}
	}
}

func TestAddUpdateDelete(t *testing.T) {
	db := createTestDatabase()
	defer deleteTestDatabase()

	// create user table instance
	table, err := tables.NewStructTagsTable(db, &User{})
	if err != nil {
		t.Fatal(err)
	}

	if err := table.Create(true); err != nil {
		t.Error(err)
	}

	user := User{
		ID:       1,
		Username: "username",
		Password: "pwd",
		Active:   false,
		// CreatedAt: time.Now(),
	}

	// Add
	if err := table.Add(&user); err != nil {
		t.Fatal(err)
	} else {
		checkUser(t, db, table, &user)
	}

	// Update
	user.Username = "username1"

	if err := table.Update(&user); err != nil {
		t.Error(err)
	} else {
		checkUser(t, db, table, &user)
	}

	// delete
	if err := table.Delete(&user); err != nil {
		t.Error(err)
	} else {
		// count should be zero
		if cnt, err := table.Count(&user); err != nil {
			t.Error(err)
		} else {
			if cnt != 0 {
				t.Errorf("count of user should be zero, but got %v", cnt)
			}
		}
	}

}

func checkUser(t *testing.T, db *sql.DB, table orm.Table, expect *User) {
	// check
	if rows, err := db.Query("SELECT id, username, password FROM " + table.Name()); err != nil {
		t.Error(err)
	} else {
		if rows.Next() {
			user1 := User{}
			if err := rows.Scan(&user1.ID, &user1.Username, &user1.Password); err != nil {
				t.Error(err)
			}
			if user1.ID != expect.ID {
				t.Errorf("user's ID is wrong: %v", user1.ID)
			}
			if user1.Username != expect.Username {
				t.Errorf("user's username is wrong: %v, expect: %v", user1.Username, expect.Username)
			}
			if user1.Password != expect.Password {
				t.Errorf("user's Password is wrong: %v", user1.Password)
			}
		} else {
			t.Error("no user found")
		}
		rows.Close()
	}
}
