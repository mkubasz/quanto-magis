package quantoCore

import (
	"fmt"
	"log"
	"os"

	"github.com/urfave/cli/v2"
)

type QuantoCLI struct {
	quantosession QuantoSession
}

func NewQuantoCLI() *QuantoCLI {
	return &QuantoCLI{
		quantosession: *NewQuantoSession(),
	}
}

func main() {
	NewQuantoCLI()
	app := &cli.App{
		Name:  "Quanto",
		Usage: "quanto manager for computing data",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "name",
				Value: "Default",
				Usage: "name of the app to run quanto in",
			},
			&cli.StringFlag{
				Name:  "mode",
				Value: "local",
				Usage: "mode to run quanto in",
			},
		},
		Action: func(ctx *cli.Context) error {
			fmt.Println("Quanto is running")
			fmt.Println("Mode:", ctx.String("mode"))
			return nil
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
