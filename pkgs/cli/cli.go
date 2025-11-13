package cli

import (
	"context"
	"log"
	"os"

	"github.com/urfave/cli/v3"

	"kavigo/pkgs/globvars"
)

func RunCli() {

	cmd := &cli.Command{
		Name:    "KaviGo",
		Usage:   "KaviGo is a simple Go-based CLI tool that automatically renames files to match the naming conventions required by Kavita",
		Version: "v4.0.0",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "d",
				Value:       "",
				Usage:       "input directory",
				Destination: &globvars.D,
				Required:    false,
			},
			&cli.StringFlag{
				Name:        "o",
				Value:       "",
				Usage:       "Output directory",
				Destination: &globvars.O,
				Required:    false,
			},
			&cli.StringFlag{
				Name:        "r",
				Value:       globvars.R,
				Usage:       "Path to the Volume ranges file (comma-delimited)",
				Destination: &globvars.R,
				Required:    false,
			},
			&cli.StringFlag{
				Name:        "c",
				Value:       globvars.C,
				Usage:       "Yaml config file",
				Destination: &globvars.C,
				Required:    false,
			},
			&cli.BoolFlag{
				Name:        "v",
				Value:       false,
				Usage:       "Verbose output",
				Destination: &globvars.V,
				Required:    false,
			},
			&cli.BoolFlag{
				Name:        "p",
				Value:       false,
				Usage:       "Preserve original files",
				Destination: &globvars.P,
				Required:    false,
			},
		},
		Action: func(ctx context.Context, cmd *cli.Command) error {
			return nil
		},
	}

	if err := cmd.Run(context.Background(), os.Args); err != nil {
		log.Fatal(err)
	}

}
