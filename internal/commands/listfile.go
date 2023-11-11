package commands

import (
	"fmt"

	"github.com/jbliao/go-mcsm-client/pkg/mcsm"
	"github.com/urfave/cli/v2"
)

func NewListFileCommand() *cli.Command {
	return &cli.Command{
		Name: "listfile",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "path",
				Usage: "`PATH` to the direcory you want to list",
				Value: "/",
			},
			&cli.StringFlag{
				Name:  "include",
				Usage: "Filter out filenames with `INCLUDE`",
			},
		},
		Action: DoListFile,
	}
}

func DoListFile(ctx *cli.Context) error {
	client, err := mcsm.NewClient(
		ctx.String("url"),
		ctx.String("api-key"),
		ctx.String("service-uuid"),
		ctx.String("instance-uuid"),
	)
	if err != nil {
		return err
	}

	out, err := client.ListFile(ctx.Context, ctx.String("path"), ctx.String("include"))
	if err != nil {
		return err
	}

	for _, item := range out.Items {
		if item.Type == 1 {
			fmt.Print("F\t")
		} else {
			fmt.Print("D\t")
		}

		fmt.Println(item.Name, "\t\t", item.Size)
	}
	return nil
}
