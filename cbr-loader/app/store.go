package main

import (
	"database/sql"
	"log"

	"github.com/gchaincl/dotsql"

	_ "github.com/lib/pq"
)

type service struct {
	db      *sql.DB
	queries *dotsql.DotSql
}

var s *service

// Init performs DB initialization
func Init() {
	db := initDB()
	queries := initQueries()
	s = &service{db: db, queries: queries}
	initTables()
}

func initDB() *sql.DB {
	connString := "user=feardude dbname=feardude sslmode=disable"
	db, err := sql.Open("postgres", connString)
	check(err)
	return db
}

func initQueries() *dotsql.DotSql {
	queries, err := dotsql.LoadFromFile("./queries.sql")
	check(err)
	return queries
}

func check(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func initTables() {
	initTable("create-table-currencies")
	initTable("create-table-fx_rates")
}

func initTable(query string) {
	query, err := s.queries.Raw(query)
	_, err = s.db.Exec(query)
	check(err)
}

// Shutdown closes DB connection pool
func Shutdown() {
	err := s.db.Close()
	if err != nil {
		log.Fatal(err)
	}
}

// AddCurrency stores new currency in DB
func AddCurrency(c Currency) {
	query, err := s.queries.Raw("insert-currency")
	check(err)

	tx, err := s.db.Begin()
	defer tx.Rollback()
	check(err)

	_, err = tx.Exec(query, c.CodeCbr, c.CodeEng, c.NameRus, c.NameEng)
	check(err)
	tx.Commit()
}

// AddFxRate stores new FX rate
func AddFxRate(cbrCode string, fxRate FxRate) {
	query, err := s.queries.Raw("insert-fx_rate")
	check(err)

	tx, err := s.db.Begin()
	defer tx.Rollback()
	check(err)

	_, err = tx.Exec(query, cbrCode, fxRate.Date, fxRate.Value)
	check(err)
	tx.Commit()
}
