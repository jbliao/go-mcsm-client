package commands

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/jbliao/go-mcsm-client/pkg/mcsm"
	"github.com/jbliao/go-mcsm-client/pkg/mcsm/services"
	"github.com/urfave/cli/v2"
)

func NewBackupCommand() *cli.Command {
	return &cli.Command{
		Name:    "backup",
		Aliases: []string{"bak"},
		Usage:   "Backup a specific instance",
		UsageText: `mcsm instance backup [command options]
		This command will create a backup zip of (worlds, worlds_nether, worlds_the_end) with name /backups/FILE_NAME
		And download it to FILE_NAME in CWD.

		If you set FILE_NAME to '-', the backup will be created with default FILE_NAME, then download and write to STDOUT
		`,
		Action: DoBackup,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "file-name",
				Aliases:     []string{"f"},
				Usage:       "`FILE_NAME`",
				DefaultText: "worlds_yyyymmdd.zip",
			},
		},
	}
}

func DoBackup(ctx *cli.Context) error {
	url := ctx.String("url")
	apiKey := ctx.String("api-key")
	serviceUuid := ctx.String("service-uuid")
	instanceUuid := ctx.String("instance-uuid")

	mcsmClient, err := mcsm.NewClient(url, apiKey, serviceUuid, instanceUuid)
	if err != nil {
		return err
	}

	y, m, d := time.Now().Date()
	defaultFileName := fmt.Sprintf("worlds_%4d%2d%2d.zip", y, m, d)
	fileName := ctx.String("file-name")

	var saveTo *os.File
	if fileName == "-" {
		saveTo = os.Stdout
		fileName = defaultFileName
	} else {
		if !strings.HasSuffix(fileName, ".zip") {
			fileName = defaultFileName
		}

		saveTo, err = os.Create(fileName)
		if err != nil {
			return err
		}
	}

	fmt.Fprintln(os.Stderr, "Download worlds to", fileName)
	return services.BackupServerWorlds(ctx.Context, mcsmClient, fileName, saveTo)
}
