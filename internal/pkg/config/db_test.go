package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDBConfig(t *testing.T) {
	cfg := DBConfig{
		Host:     "localhost",
		Port:     5432,
		DBName:   "iam_permission",
		User:     "user",
		Password: "password",
	}
	assert.Equal(t, "host=localhost port=5432 user=user password=password dbname=iam_permission sslmode=disable", cfg.DSN())
}
