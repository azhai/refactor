package main

import (
	"fmt"
	"os"

	"gitea.com/azhai/refactor"
	"gitea.com/azhai/refactor/config"
	_ "github.com/denisenkom/go-mssqldb"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
	"github.com/urfave/cli/v2"
)

const VERSION = "0.9.1"

var ReverseFlags = []cli.Flag{
	&cli.StringFlag{
		Name:    "file",
		Aliases: []string{"f"},
		Usage:   "yml file to apply for reverse",
		Value:   "settings.yml",
	},
}

func ReverseAction(ctx *cli.Context) error {
	configFile := ctx.String("file")
	if configFile == "" {
		return fmt.Errorf("need reverse file")
	}
	cfg, err := config.ReadSettings(configFile)
	if err != nil {
		return err
	}

	names := ctx.Args().Slice()
	return refactor.ExecReverseSettings(cfg, names...)
}

func main() {
	app := &cli.App{
		Name:      "reverse",
		Version:   VERSION,
		Usage:     "Reverse is a database reverse command line tool",
		UsageText: "reverse [global options] [arguments...]",
		Flags:     ReverseFlags,
		Action:    ReverseAction,
	}

	if err := app.Run(os.Args); err != nil {
		panic(err)
	}
}
