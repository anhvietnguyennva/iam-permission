package app

import (
	"errors"
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	"github.com/urfave/cli/v2"

	apiutil "iam-permission/internal/app/api/component"
	apiserver "iam-permission/internal/app/api/server"
	"iam-permission/internal/app/migration"
	"iam-permission/internal/pkg/config"
	"iam-permission/internal/pkg/db"
)

func APIServerCommand() *cli.Command {
	flagCfg := "config"
	flagAutoMigration := "auto-migration"
	flagMigrationDir := "migration-dir"
	flagAdmin := "admin"

	defaultCfgFile := "internal/pkg/config/file/default.yaml"
	defaultMigrationDir := "file://./migration/postgres"

	return &cli.Command{
		Name:    "api",
		Aliases: []string{},
		Usage:   "Run api server",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  flagCfg,
				Value: defaultCfgFile,
			},
			&cli.BoolFlag{
				Name:  flagAutoMigration,
				Value: false,
			},
			&cli.StringFlag{
				Name:  flagMigrationDir,
				Value: defaultMigrationDir,
			},
			&cli.BoolFlag{
				Name:  flagAdmin,
				Value: false,
			},
		},
		Action: func(c *cli.Context) (err error) {
			err = config.Load(c.String(flagCfg))
			if err != nil {
				return err
			}

			// auto migration
			if c.Bool(flagAutoMigration) {
				fmt.Println("-------- Run migration --------")
				err := db.InitDB()
				if err != nil {
					return err
				}
				m, err := migration.NewMigration(db.Instance(), defaultMigrationDir, "")
				if err != nil {
					fmt.Println("Can not create migration " + err.Error())
					return err
				}
				err = m.MigrateUp(0)
				if err != nil && !errors.Is(err, migrate.ErrNoChange) {
					return err
				}
			}

			err = apiutil.InitComponents()
			if err != nil {
				return err
			}

			server := apiserver.NewAPIServer(c.Bool(flagAdmin))
			if err != nil {
				return err
			}
			return server.Run()
		},
	}
}
