# Golang ORM Framework

+ Create table via a custom struct
+ Base CRUD

## Usage

```bash

go get "github.com/zgljl2012/go-orm"

```

### Define Table

You need to add tags for your field of the struct just like `json:"..."`. The `name` tag is required, if not specify `name` tag, the field won't be parsed as an orm.field.

And you should specify the `length` tag for all of your `char` field or an error will be reported.

The struct at least has one field that specifies `primaryKey` tag.

```golang

// User is a test table
type User struct {
    ID        int       `name:"id" primaryKey:"true"`
    Username  string    `name:"username" length:"20"`
    Password  string    `name:"password" length:"50"`
    Active    bool      `name:"active" null:"false"`
    Age       float32   `name:"age"`
    CreatedAt time.Time `name:"created_at"`
    Count     uint64    `name:"count"`
}

```

Supported Type:

+ `Int`
+ `Float`
+ `Bool`
+ `Datetime`
+ `Char`
+ `Uint64` (`BigInt`)

### Create Table

```golang

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
    ID        int       `name:"id" primaryKey:"true"`
    Username  string    `name:"username" length:"20"`
    Password  string    `name:"password" length:"50"`
    Active    bool      `name:"active" null:"false"`
    Age       float32   `name:"age"`
    CreatedAt time.Time `name:"created_at"`
    Count     uint64    `name:"count"`
}

// create test db
func createTestDatabase() *sql.DB {
    db, err := sql.Open("sqlite3", testDB)
    if err != nil {
        log.Fatal(err)
    }
    return db
}

// delete test db
func deleteTestDatabase() {
    if err := os.Remove(testDB); err != nil {
        log.Fatal(err)
    }
}

func TestCreateTable(t *testing.T) {
    db := createTestDatabase()
    defer deleteTestDatabase()

    // create user table instance
    table, err := tables.NewStructTagsTable(db, &User{})
    if err != nil {
        t.Fatal(err)
    }

    // create table in database, name is the same as struct
    if err := table.Create(false); err != nil {
        t.Error(err)
    }
}

```

### Add/Update/Delete

```golang

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
    }

    // Add
    if err := table.Add(&user); err != nil {
        t.Fatal(err)
    }

    // Update
    user.Username = "username1"

    if err := table.Update(&user); err != nil {
        t.Error(err)
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

```

### Filter

Do not support `in`, `like` and so on now.

```golang

func TestFilterSet(t *testing.T) {
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
        ID:        1,
        Username:  "username1",
        Password:  "pwd",
        Active:    false,
        CreatedAt: time.Now(),
    }

    // Add
    if err := table.Add(&user); err != nil {
        t.Fatal(err)
    }

    // Filter
    filter := table.Filter()
    rows := filter.All()
    id := 1
    for _, row := range rows {
        user := row.(User)
        if id != user.ID {
            t.Errorf("ID should be %v, but got %v", id, user.ID)
        }
        id += 1
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
    filter = table.Filter()
    rows = filter.All()
    id = 1
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

    // filter with id=1
    filter = table.Filter(orm.WithParameter("ID", 1))
    rows = filter.All()
    if len(rows) != 1 {
        t.Error("You should only filter one row")
    }
    user1 := rows[0].(User)
    if user1.ID != 1 {
        t.Errorf("Expect id is 1, but got %v", user1.ID)
    }

    // orderby
    user1 = table.Filter().OrderBy("-ID").All()[0].(User)
    if user1.ID != 10 {
        t.Errorf("ID of this user should be 10, but got %v", user1.ID)
    }

    // limit
    rows = table.Filter().Limit(5).All()
    if len(rows) != 5 {
        t.Errorf("rows'cnt should be 5, but got %v", len(rows))
    }

    // offset
    rows = table.Filter().Offset(2).All()
    if rows[0].(User).ID != 3 {
        t.Errorf("expected 3, but got %v", rows[0].(User).ID)
    }
}

```
