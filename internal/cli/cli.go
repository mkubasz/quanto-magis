// Package cli provides command-line interface functionality for Quanto.
package cli

import (
	"log"
	"os"

	"github.com/urfave/cli/v2"

	"mkubasz/quanto/internal/session"
)

// QuantoCLI represents the command-line interface for Quanto applications.
type QuantoCLI struct {
	Session *session.QuantoSession
}

// New creates a new QuantoCLI instance with an initialized session.
func New() *QuantoCLI {
	return &QuantoCLI{
		Session: session.New(),
	}
}

// Run starts the CLI application and processes command-line arguments.
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
			log.Println("Quanto is running")
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
