package main

import (
	"os"

	"app/common/config"
	// "app/cronjob"
	"app/services"
	"app/web"

	"github.com/urfave/cli"
)

var (
	cfg = config.GetConfig()
)

func main() {
	app := cli.NewApp()
	app.Commands = []cli.Command{
		{
			Name:    "web",
			Aliases: []string{"s"},
			Action: func(c *cli.Context) error {
				return web.NewAppSrv().Start()
			},
		},
		{
			// Bộ api cho user
			Name: "api.autofarmer.net",
			Action: func(c *cli.Context) error {
				services.NewLoadServiceToQueue().Run()
				return web.NewAutofarmerNetAPISrv().Start()
			},
		},
		{
			// Bộ api cho employee
			Name: "admin.api.autofarmer.net",
			Action: func(c *cli.Context) error {
				return web.NewAdminAutofarmerNetAPISrv().Start()
			},
		},
		{
			// Bộ api cho employee
			Name: "api8.autofarmer.net",
			Action: func(c *cli.Context) error {
				// cronjob.NewUpdateLive().AddFuncCron()
				return web.NewAutofarmerNetAPI8().Start()
			},
		},
	}

	if len(os.Args) == 1 {
		os.Args = append(os.Args, "s")
	}

	app.Run(os.Args)
}
