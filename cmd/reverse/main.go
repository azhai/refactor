package main

import (
	"os"

	"gitea.com/azhai/refactor"
	"gitea.com/azhai/refactor/cmd"
	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App {
		HideHelp: true,
		Version: cmd.VERSION,
		Usage: "从数据库导出对应的Model代码",
		Action: ReverseAction,
	}
	app.Flags = []cli.Flag{
		&cli.StringFlag{
			Name:    "file",
			Aliases: []string{"f"},
			Usage:   "配置文件路径",
			Value:   "settings.yml",
		},
	}
	err := app.Run(os.Args)
	if err != nil {
		panic(err)
	}
}

func ReverseAction(ctx *cli.Context) (err error) {
	names := ctx.Args().Slice()
	cfg := cmd.Prepare(ctx.String("file"))
	err = refactor.ExecReverseSettings(cfg, names...)
	return
}
