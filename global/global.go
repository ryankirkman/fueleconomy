package global

import (
	"database/sql"
	"encoding/json"
	"errors"
	"io"
	"log"
	"os"

	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"

	"github.com/teasherm/fueleconomy/srm"
)

var (
	Db     *srm.DbMap
	Logger *log.Logger
)

// Holds postgres connection string
type Config struct {
	Db string `json:"db"`
}

func GetDbConfig() (string, error) {
	var config Config

	file, _ := os.Open(os.Getenv("CONFIG_PATH"))
	decoder := json.NewDecoder(file)
	err := decoder.Decode(&config)
	if err != nil {
		return "", err
	}

	return config.Db, nil
}

func InitDb(driver, connString string) error {
	db, err := sql.Open(driver, connString)
	if err != nil {
		return err
	}

	var dialect srm.Dialect

	switch driver {
	case "postgres":
		dialect = srm.PostgresDialect{}
	case "sqlite3":
		dialect = srm.Sqlite3Dialect{}
	default:
		return errors.New("global.InitDb: Driver not supported")
	}

	Db = &srm.DbMap{Conn: db, Dialect: dialect}

	return nil
}

// TODO: glog?
func InitLogger(w io.Writer) {
	Logger = log.New(w, "INFO: ", log.LstdFlags)
}
