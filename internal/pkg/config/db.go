package config

import "fmt"

// DBConfig contains all configuration of database (PostgreSQL)
type DBConfig struct {
	Host         string `default:"localhost"`
	Port         int    `default:"5432"`
	DBName       string `default:"iam_permission"`
	User         string `default:"iam_permission"`
	Password     string `default:"123456"`
	ConnLifeTime int    `default:"300"`
	MaxIdleConns int    `default:"100"`
	MaxOpenConns int    `default:"200"`
	LogLevel     int    `default:"1"`
}

func (c *DBConfig) DSN() string {
	return fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		c.Host, c.Port, c.User, c.Password, c.DBName,
	)
}
