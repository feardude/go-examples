package main

import (
	"time"
)

// Currency represents currency structure
type Currency struct {
	CodeCbr string
	CodeEng string
	NameRus string
	NameEng string
}

// FxRate represents currency rate structure
type FxRate struct {
	EngCode string
	Date    time.Time
	Value   float32
}
