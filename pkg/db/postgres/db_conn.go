package postgres

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/jackc/pgx/v5/stdlib"
)

var myDb *sql.DB

// Open opens a database without having to take care of the logic, just calling
// the function should give you an error if any.
func open() error {
	var db *sql.DB
	db, err := sql.Open("pgx", os.Getenv("DATABASE_URL"))

	if err != nil {
		fmt.Printf("Error while trying to connect with the database: %v\n", err)
		return err
	}

	if err = db.Ping(); err != nil {
		fmt.Printf("Error while trying to verify the connection with the database: %v\n", err)
		return err
	}

	myDb = db
	return nil
}

func GetDb() *sql.DB {
	if myDb != nil {
		return myDb
	}
	_ = open()
	return myDb
}

func Disconnect() error {
	return myDb.Close()
}