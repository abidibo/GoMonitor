package db

import (
	"database/sql"
	"errors"
	"os"
	"path/filepath"
	"sync"

	_ "github.com/mattn/go-sqlite3"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

type DBClient struct {
	C *sql.DB
}

var (
	db   *DBClient
	once sync.Once
)

// DB returns a singleton pointer to pgx connection instance
func DB() *DBClient {
	if db == nil {
		once.Do(func() {
			dbFile := filepath.Join(viper.GetString("app.homePath"), "gomonitor.db")
			conn, err := sql.Open("sqlite3", dbFile)
			if err != nil {
				zap.S().Fatal("Cannot connect to sqlite3 db ", dbFile)
			} else {
				zap.S().Info("Succesfully connected to sqlite3 db on ", dbFile)
			}

			db = &DBClient{conn}
		})
	}

	return db
}

func InitDatabase() {
	dbFile := filepath.Join(viper.GetString("app.homePath"), "gomonitor.db")
	if _, err := os.Stat(dbFile); errors.Is(err, os.ErrNotExist) {
		os.OpenFile(dbFile, os.O_RDONLY|os.O_CREATE, 0666)
	}

	// sqlite3
	// create a new log table if not exists with the following columns:
	// id: auto increment
	// timestamp: timestamp of the log
	// user string
	// partial_time_minutes: total time usage for current day
	stmt := `
CREATE TABLE IF NOT EXISTS log (
id INTEGER PRIMARY KEY AUTOINCREMENT,
timestamp DATETIME DEFAULT CURRENT_TIMESTAMP,
user TEXT,
partial_time_minutes INTEGER
);
		`
	client := DB()
	_, err := client.C.Exec(stmt)
	if err != nil {
		panic(err)
	}

	// sqlite3
	// create a new process log table if not exists with the following columns:
	// id: auto increment
	// log_id: foreign key to log table
	// name: process Name
	// cpu_percent: process CPU usage
	// memory_percent: process memory usage
	stmt = `
CREATE TABLE IF NOT EXISTS log_process (
id INTEGER PRIMARY KEY AUTOINCREMENT,
log_id INTEGER,
name TEXT,
cpu_percent REAL,
memory_percent REAL
);
		`
	_, err = client.C.Exec(stmt)
	if err != nil {
		panic(err)
	}
}
