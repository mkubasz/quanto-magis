package main

import (
	"fmt"
	"log"
	"os"

	"mkubasz/quanto/quantoCore"

	"github.com/urfave/cli/v2"
)

type QuantoCLI struct {
	quantosession quantoCore.QuantoSession
}

func NewQuantoCLI() *QuantoCLI {
	return &QuantoCLI{
		quantosession: *quantoCore.NewQuantoSession(),
	}
}

func Init() *QuantoCLI {
	quantoCli := NewQuantoCLI()
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
			quantoCli.quantosession.AppName = ctx.String("name")
			switch ctx.String("mode") {
			case "local":
				quantoCli.quantosession.Mode = quantoCore.Local
			case "cluster":
				quantoCli.quantosession.Mode = quantoCore.Cluster
			default:
				quantoCli.quantosession.Mode = quantoCore.Local
			}
			return nil
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}

	return quantoCli
}

func main() {
	c := quantoCore.NewDataFrame([]interface{}{[]interface{}{"A", "B", "A", "D", "E"}, []interface{}{1, 2, 3, 4, 5}}, []string{"col1", "col2"})
	fmt.Println(c.HasColumn("col1"))
	quantoCli := Init()
	fmt.Println(quantoCli.quantosession)
}
