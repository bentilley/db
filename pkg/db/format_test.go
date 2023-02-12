package db_test

import (
	"testing"

	"github.com/bentilley/db/pkg/db"
	"github.com/stretchr/testify/assert"
)

func TestColumnFormat(t *testing.T) {
	input := [][]string{
		{"1", "postgres://localhost:5432", "Postgres database"},
		{"2", "mysql://some.host:12345", "MYSQL database"},
		{"30", "mssql://another.host.io:5555", "MSSQL database"},
	}
	expected := []string{
		"1  postgres://localhost:5432    Postgres database",
		"2  mysql://some.host:12345      MYSQL database",
		"30 mssql://another.host.io:5555 MSSQL database",
	}
	actual := db.ColumnFormat(input)
	assert.Equal(t, expected, actual)
}
