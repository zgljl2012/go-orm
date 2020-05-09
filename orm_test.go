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
	if _, err := db.Exec(fmt.Sprintf("SELECT COUNT(*) FROM %s", table.Name())); err != nil {
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

}
