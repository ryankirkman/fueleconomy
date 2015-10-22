# srm (struct relational mapper)

Package for persisting structs to SQL database tables. Essentially a pared down version of [gorp](https://github.com/go-gorp/gorp).

Supports:
- PostgreSQL
- SQLite 3

## Usage

```go
package main

import (
    "database/sql"

    "github.com/teasherm/fueleconomy/srm"
)

func main() {
    conn, _ := sql.Open("postgres", "host=localhost dbname=db user=user password=pass sslmode=disable")

    type Model struct {
        Field string `db:"field"`
    }

    Db := &srm.DbMap{Conn: conn, Dialect: "postgres"}
    Db.InsertOne("models", &Model{Field: "value"})

    result := Model{}
    Db.SelectOne(&result, "SELECT * FROM models WHERE field = $1", "value")

    // Prints "value"
    fmt.Println(result.Field)
}
```
