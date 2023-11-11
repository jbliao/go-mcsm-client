package main

import (
	"fmt"
	"os"

	"github.com/jbliao/go-mcsm-client/internal"
)

func main() {
	app := internal.NewApp()
	if err := app.Run(os.Args); err != nil {
		fmt.Println("ERROR", err)
	}
}
