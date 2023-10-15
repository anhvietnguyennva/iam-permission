package main

import (
	"log"
	"os"
	"reflect"

	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/urfave/cli/v2"

	"iam-permission/cmd/app"
)

func main() {
	cmd := &cli.App{
		Name: "IAM Permission",
		Commands: []*cli.Command{
			app.APIServerCommand(),
			app.MigrationCommand(),
		},
	}

	err := cmd.Run(os.Args)
	if err != nil && !reflect.ValueOf(&err).IsNil() {
		log.Fatal(err)
	}
}
