package store

import (
	"database/sql"
	"log"

	"github.com/gchaincl/dotsql"

	// Imports Postgresql driver
	_ "github.com/lib/pq"
)

// Init performs DB initialization
func Init() {
	connString := "user=feardude dbname=feardude sslmode=disable"
	db, err := sql.Open("postgres", connString)
	defer db.Close()

	if err != nil {
		log.Fatal(err)
	}

	initTables(db)
}

func initTables(db *sql.DB) {
	queries, err := dotsql.LoadFromFile("./store/queries.sql")

	if err != nil {
		log.Fatal(err)
	}

	initTable(db, queries, "create-table-currencies")
	initTable(db, queries, "create-table-fx_rates")
}

func initTable(db *sql.DB, queries *dotsql.DotSql, query string) {
	query, err := queries.Raw(query)
	_, err = db.Exec(query)

	if err != nil {
		log.Fatal(err)
	}
}
