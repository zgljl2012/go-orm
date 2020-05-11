package orm_test

import (
	"database/sql"
	"fmt"

	"os"
	"testing"

	_ "github.com/mattn/go-sqlite3"
	"github.com/zgljl2012/go-orm"
	"github.com/zgljl2012/go-orm/fields"
	"github.com/zgljl2012/go-orm/tables"
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
		fields.NewCharField("Username", fields.WithLength(20)),
		fields.NewCharField("Password", fields.WithLength(50)),
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
				"type":   "CHAR(20)",
				"pk":     false,
			},
			"Password": {
				"exists": false,
				"type":   "CHAR(50)",
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
				t.Errorf("field %v not found", name)
			}
		}
	}

}

func TestAddUpdateDelete(t *testing.T) {
	db := createTestDatabase()
	defer deleteTestDatabase()

	// create user table instance
	table, err := tables.NewTable(db, &User{})
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
	if rows, err := db.Query("SELECT ID, Username, Password FROM " + table.Name()); err != nil {
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

func TestFilterSet(t *testing.T) {
	db := createTestDatabase()
	defer deleteTestDatabase()

	// create user table instance
	table, err := tables.NewTable(db, &User{})
	if err != nil {
		t.Fatal(err)
	}

	if err := table.Create(true); err != nil {
		t.Error(err)
	}

	user := User{
		ID:       1,
		Username: "username1",
		Password: "pwd",
	}

	// Add
	if err := table.Add(&user); err != nil {
		t.Fatal(err)
	}

	// Filter
	if filter, err := table.Filter(); err != nil {
		t.Error(err)
	} else {
		rows := filter.All()
		id := 1
		for _, row := range rows {
			user := row.(User)
			if id != user.ID {
				t.Errorf("ID should be %v, but got %v", id, user.ID)
			}
			id += 1
		}
	}

	// Add
	for i := 1; i < 10; i++ {
		user.ID = i + 1
		user.Username = fmt.Sprintf("username%d", i+1)
		if err := table.Add(&user); err != nil {
			t.Fatal(err)
		}
	}

	// validate
	if filter, err := table.Filter(); err != nil {
		t.Error(err)
	} else {
		rows := filter.All()
		id := 1
		for _, row := range rows {
			user := row.(User)
			if id != user.ID {
				t.Errorf("ID should be %v, but got %v", id, user.ID)
			}
			if user.Username != fmt.Sprintf("username%d", id) {
				t.Errorf("ID should be %v, but got %v", fmt.Sprintf("username%d", id), user.Username)
			}
			id += 1
		}
	}

	// filter with id=1
	if filter, err := table.Filter(orm.WithParameter("ID", 1)); err != nil {
		t.Error(err)
	} else {
		rows := filter.All()
		if len(rows) != 1 {
			t.Error("You should only filter one row")
		}
		user1 := rows[0].(User)
		if user1.ID != 1 {
			t.Errorf("Expect id is 1, but got %v", user1.ID)
		}
	}

}
