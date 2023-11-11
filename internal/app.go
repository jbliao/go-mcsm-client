package internal

import (
	"github.com/jbliao/go-mcsm-client/internal/commands"
	"github.com/urfave/cli/v2"
)

type App struct {
	cli.App
}

func NewApp() *App {
	return &App{
		App: cli.App{
			Name:   "mcsm",
			Usage:  "A cli tool for MCSManager",
			Action: cli.ShowAppHelp,
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:     "url",
					Usage:    "The `url` to your MCSM panel",
					EnvVars:  []string{"MCSM_URL"},
					Required: true,
				},
				&cli.StringFlag{
					Name:     "api-key",
					Usage:    "The `apikey` to your MCSM panel",
					EnvVars:  []string{"MCSM_API_KEY"},
					Required: true,
				},
			},
			Commands: []*cli.Command{
				{
					Name: "instance",
					Flags: []cli.Flag{
						&cli.StringFlag{
							Name:     "service-uuid",
							Aliases:  []string{"sid"},
							Usage:    "The service `UUID` of the instance (Named GUID on UI)",
							EnvVars:  []string{"MCSM_SERVICE_UUID"},
							Required: true,
						},
						&cli.StringFlag{
							Name:     "instance-uuid",
							Aliases:  []string{"iid"},
							Usage:    "The `UUID` of the instance",
							EnvVars:  []string{"MCSM_INSTANCE_UUID"},
							Required: true,
						},
					},
					Usage: "Instance specific commands",
					Subcommands: []*cli.Command{
						commands.NewBackupCommand(),
						commands.NewLiseFileCommand(),
						commands.NewExecuteCommand(),
					},
				},
			},
		},
	}
}
