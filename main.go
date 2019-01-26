package main

import (
	"log"
	"os"

	"github.com/urfave/cli"
	"github.com/fox-kazuki-matsumoto/cli-sample/api"
)

func main() {
	app := cli.NewApp()
	app.Name = "cli-sample"
	app.Usage = "search qiita articles"
	app.Action = func(c *cli.Context) error {
		api.RunQiitaSearch()
		return nil
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
