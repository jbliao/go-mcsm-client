package services

import (
	"context"
	"io"
	"time"

	"github.com/jbliao/go-mcsm-client/pkg/mcsm"
)

// BackupServerWorlds issue zip command on a mcsm instance, then download it and write using writer.
func BackupServerWorlds(ctx context.Context, mcsmClient *mcsm.Client, fileNameOnServer string, writer io.Writer) (err error) {

	dirName := "backups"
	zipName := dirName + "/" + fileNameOnServer

	err = mcsmClient.SendCommand(ctx, "save-off")
	if err != nil {
		return
	}

	defer mcsmClient.SendCommand(ctx, "save-on")

	err = mcsmClient.ZipFiles(ctx, zipName, "world", "world_nether", "world_the_end")
	if err != nil {
		return
	}

	for {
		var result *mcsm.FileStatus
		result, err = mcsmClient.FileStatus(ctx)
		if err != nil {
			return
		}

		if result.InstanceFileTask == 0 {
			break
		}

		time.Sleep(time.Second)
	}

	reader, err := mcsmClient.DownloadFile(ctx, zipName)
	if err != nil {
		return
	}

	io.Copy(writer, reader)
	return
}
