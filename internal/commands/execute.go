package commands

import (
	"fmt"
	"os"
	"strings"

	"github.com/jbliao/go-mcsm-client/pkg/mcsm"
	"github.com/urfave/cli/v2"
)

func NewExecuteCommand() *cli.Command {
	return &cli.Command{
		Name:   "execute",
		Action: DoExecute,
	}
}

func DoExecute(ctx *cli.Context) error {
	client, err := mcsm.NewClient(
		ctx.String("url"),
		ctx.String("api-key"),
		ctx.String("service-uuid"),
		ctx.String("instance-uuid"),
	)
	if err != nil {
		return err
	}

	var mcCommand = strings.Join(ctx.Args().Slice(), " ")
	fmt.Fprintln(os.Stderr, "Executing", mcCommand)
	if err := client.SendCommand(ctx.Context, mcCommand); err != nil {
		return err
	}
	fmt.Fprintln(os.Stderr, "Executed", mcCommand)
	return nil
}
