package tables_test

import (
	"database/sql"
	"fmt"
	"os"
	"testing"

	_ "github.com/mattn/go-sqlite3"
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
	Username string `name:"username"`
	Password string `name:"password"`
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
}
