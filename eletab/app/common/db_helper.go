package common

import (
	_ "github.com/mattn/go-sqlite3"
	"path"
	"github.com/jmoiron/sqlx"
	"database/sql"
	"os"
)

// Database file name
const DBFile = "eletab_lite.db"

// Database path
var DBPath = os.Getenv("ELETAB_PATH")

type DatabaseHelper interface {
	CreateTable(database *Database) error
}

type Database struct {
	*sqlx.DB
}

// Opens database and it i returns it
func OpenDB(helper DatabaseHelper) (*Database, error) {
	var db *Database
	sqlDB, err := sqlx.Open("sqlite3", path.Join(DBPath, DBFile))
	//sqlDB, err := sqlx.Open("sqlite3", DBFile)
	if err == nil {
		db = &Database{sqlDB}
		err = helper.CreateTable(db)
	}

	return db, err
}

// Execute queries with params
func (db *Database) Exec(qry string, data... interface{}) (sql.Result, error) {
	statement, err := db.Prepare(qry)
	if err != nil {
		return nil, err
	}

	var res sql.Result
	if len(data) > 0 {
		res, err = statement.Exec(data...)
	} else {
		res, err = statement.Exec()
	}
	if err != nil {
		return nil, err
	}

	return res, nil
}

// Executes DB query with data from struct
func (db *Database) NamedExec(qry string, obj interface{}) (sql.Result, error) {
	tx := db.MustBegin()
	res, err := tx.NamedExec(qry, obj)
	if err != nil {
		return nil, err
	}
	return res, tx.Commit()
}

// Closes the database
func (db *Database) Close() {
	db.DB.Close()
}
