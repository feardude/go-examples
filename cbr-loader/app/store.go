package main

import (
	"database/sql"
	"log"
	"time"

	"github.com/gchaincl/dotsql"

	_ "github.com/lib/pq"
)

type service struct {
	db      *sql.DB
	queries *dotsql.DotSql
}

var s *service

// InitDB performs DB initialization
func InitDB() {
	db := initDB()
	queries := initQueries()
	s = &service{db: db, queries: queries}
}

func initDB() *sql.DB {
	connString := "host=postgres user=postgres dbname=postgres sslmode=disable"
	db, err := sql.Open("postgres", connString)
	check(err)
	db.SetMaxOpenConns(100)
	return db
}

func initQueries() *dotsql.DotSql {
	queries, err := dotsql.LoadFromFile("./queries.sql")
	check(err)
	return queries
}

// ShutdownDB closes DB connection pool
func ShutdownDB() {
	err := s.db.Close()
	if err != nil {
		log.Fatal(err)
	}
}

// AddFxRate stores new FX rate
func AddFxRate(fxRate FxRate) {
	query, err := s.queries.Raw("insert-fx_rate")
	check(err)

	tx, err := s.db.Begin()
	defer tx.Rollback()
	check(err)

	_, err = tx.Exec(query, fxRate.CbrCode, fxRate.Date, fxRate.Value)
	check(err)
	tx.Commit()
}

// GetLastDate finds last FX rate date for cbrCode currency
func GetLastDate(cbrCode string) time.Time {
	query, err := s.queries.Raw("select-last-date")
	check(err)

	result, err := s.db.Query(query, cbrCode)
	check(err)

	var lastDate time.Time
	result.Next()
	result.Scan(&lastDate)

	return lastDate
}

// GetCurrencies returns all currencies currently stored
func GetCurrencies() []Currency {
	query, err := s.queries.Raw("select-currencies")
	check(err)

	rows, err := s.db.Query(query)
	check(err)

	currencies := make([]Currency, 0)
	for rows.Next() {
		var currency Currency
		err := rows.Scan(&currency.CodeCbr, &currency.CodeEng, &currency.NameRus, &currency.NameEng)
		check(err)
		currencies = append(currencies, currency)
	}
	return currencies
}
