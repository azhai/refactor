package main

import (
	"os"

	"gitea.com/azhai/refactor"
	"gitea.com/azhai/refactor/cmd"
	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Version: cmd.VERSION,
		Usage:   "从数据库导出对应的Model代码",
		Action:  ReverseAction,
	}
	app.Flags = []cli.Flag{
		&cli.BoolFlag{
			Name:    "verbose",
			Aliases: []string{"vv"},
			Usage:   "输出详细信息",
		},
		&cli.StringFlag{
			Name:    "file",
			Aliases: []string{"c", "f"},
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
	cfg, verbose := cmd.Prepare(ctx.String("file")), cmd.Verbose()
	names := ctx.Args().Slice()
	err = refactor.ExecReverseSettings(cfg, verbose, names...)
	return
}
