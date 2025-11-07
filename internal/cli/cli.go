package cli

import (
	"fmt"
	"log"
	"os"

	"github.com/urfave/cli/v2"

	"mkubasz/quanto/internal/session"
)

type QuantoCLI struct {
	Session *session.QuantoSession
}

func New() *QuantoCLI {
	return &QuantoCLI{
		Session: session.New(),
	}
}

func (q *QuantoCLI) Run() error {
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
			q.Session.AppName = ctx.String("name")
			switch ctx.String("mode") {
			case "local":
				q.Session.Mode = session.Local
			case "cluster":
				q.Session.Mode = session.Cluster
			default:
				q.Session.Mode = session.Local
			}
			return nil
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}

	return nil
}
