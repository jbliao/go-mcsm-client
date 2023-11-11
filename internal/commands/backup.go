package commands

import (
	"github.com/jbliao/go-mcsm-client/pkg/mcsm"
	"github.com/jbliao/go-mcsm-client/pkg/mcsm/services"
	"github.com/urfave/cli/v2"
)

func NewBackupCommand() *cli.Command {
	return &cli.Command{
		Name:    "backup",
		Aliases: []string{"bak"},
		Usage:   "Backup a specific instance",
		Action:  DoBackup,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "file-name",
				Aliases: []string{"f"},
			},
		},
	}
}

func DoBackup(cCtx *cli.Context) error {
	url := cCtx.String("url")
	apiKey := cCtx.String("api-key")
	serviceUuid := cCtx.String("service-uuid")
	instanceUuid := cCtx.String("instance-uuid")

	mcsmClient, err := mcsm.NewClient(url, apiKey, serviceUuid, instanceUuid)
	if err != nil {
		return err
	}

	services.BackupServerWorlds(cCtx.Context, mcsmClient)
	return nil
}
