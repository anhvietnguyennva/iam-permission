package migration

import (
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"gorm.io/gorm"
)

type Migration interface {
	MigrateUp(up int) error
	MigrateDown(down int) error
}

type postgresMigration struct {
	m *migrate.Migrate
}

func NewMigration(db *gorm.DB, dir string, migrationTable string) (Migration, error) {
	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	if len(migrationTable) == 0 {
		migrationTable = "schema_migrations"
	}

	driver, err := postgres.WithInstance(sqlDB, &postgres.Config{
		MigrationsTable: migrationTable,
	})
	if err != nil {
		return nil, err
	}

	m, err := migrate.NewWithDatabaseInstance(
		dir,
		"iam_permission",
		driver,
	)
	if err != nil {
		return nil, err
	}

	return &postgresMigration{m: m}, nil

}

func (t *postgresMigration) MigrateUp(up int) error {
	if up == 0 {
		return t.m.Up()
	}
	return t.m.Steps(up)
}

func (t *postgresMigration) MigrateDown(down int) error {
	if down == 0 {
		return t.m.Down()
	}
	return t.m.Steps(-down)
}
