package services

import (
	"context"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/jbliao/go-mcsm-client/pkg/mcsm"
)

func BackupServerWorlds(ctx context.Context, mcsmClient *mcsm.Client) (err error) {

	y, m, d := time.Now().Date()
	fileName := fmt.Sprintf("worlds_%2d%2d%2d.zip", y, m, d)
	dirName := "backups"
	zipName := dirName + "/" + fileName

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

	f, _ := os.Create(fileName)
	io.Copy(f, reader)
	return
}
