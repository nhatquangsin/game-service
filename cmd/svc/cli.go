package main

import (
	"os"

	"github.com/urfave/cli/v2"
)

// List of common keys used in service.
const (
	EnvKeyServiceName      = "SERVICE_NAME"
	EnvKeyServiceComponent = "SERVICE_COMPONENT"
	EnvKeyServiceVersion   = "SERVICE_VERSION"
)

// Service info. These values will be injected via "ldflags" in build time.
var (
	ServiceName    = "item"
	ServiceVersion = "0.0.1"
)

func newApp() *cli.App {
	app := cli.NewApp()

	app.Name = ServiceName
	app.Authors = []*cli.Author{
		{
			Name: "Sky Mavis Backend Team",
		},
	}
	app.Commands = []*cli.Command{
		runCmd,
	}

	return app
}

var runCmd = &cli.Command{
	Name:  "run",
	Usage: "Runs a component in service",
	Subcommands: []*cli.Command{
		runAPICmd,
	},
}

var runAPICmd = &cli.Command{
	Name:  "api",
	Usage: "Runs api component in service",
	Flags: []cli.Flag{},
	Action: func(c *cli.Context) error {
		if err := os.Setenv(EnvKeyServiceComponent, "api"); err != nil {
			return err
		}
		app := newAPIApp()
		app.Run()

		return nil
	},
}
