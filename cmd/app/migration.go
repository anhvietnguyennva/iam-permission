package app

import (
	"errors"

	"github.com/urfave/cli/v2"

	"iam-permission/internal/app/migration"
	"iam-permission/internal/pkg/config"
	"iam-permission/internal/pkg/db"
)

func MigrationCommand() *cli.Command {
	flagCfg := "config"
	flagUp := "up"
	flagDown := "down"

	defaultCfgFile := "internal/pkg/config/file/default.yaml"
	defaultMigrationDir := "file://./migration/postgres"

	return &cli.Command{
		Name:    "migration",
		Aliases: []string{},
		Usage:   "Run database migration",
		Flags: []cli.Flag{
			&cli.IntFlag{
				Name:  flagUp,
				Value: -1,
			},
			&cli.IntFlag{
				Name:  flagDown,
				Value: -1,
			},
			&cli.StringFlag{
				Name:  flagCfg,
				Value: defaultCfgFile,
			},
		},
		Action: func(c *cli.Context) (err error) {
			err = config.Load(c.String(flagCfg))
			if err != nil {
				return err
			}

			up := c.Int(flagUp)
			down := c.Int(flagDown)

			if up == -1 && down == -1 {
				return errors.New("no up or down migration declared")
			}

			if up != -1 && down != -1 {
				return errors.New("both up and down migration declared. stop the migration")
			}

			if err := db.InitDB(); err != nil {
				return err
			}

			m, err := migration.NewMigration(db.Instance(), defaultMigrationDir, "")
			if err != nil {
				return err
			}

			if up != -1 {
				return m.MigrateUp(up)
			}
			return m.MigrateDown(down)
		},
	}
}
