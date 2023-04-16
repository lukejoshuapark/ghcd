package main

import (
	"fmt"
	"os"

	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Name:  "ghcd",
		Usage: "A CLI tool that enables change detection in GitHub workflows.",
		Commands: []*cli.Command{
			{
				Name:  "detect",
				Usage: "Detect changes between a historical commit and the present",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:  "f",
						Value: "ghcd.yml",
						Usage: "the configuration file to use, ghcd.yml by default",
					},
					&cli.StringFlag{
						Name:  "token",
						Usage: "the GitHub token supplied by GitHub workflow for EnvironmentDiff mode",
					},
					&cli.StringFlag{
						Name:  "repository",
						Usage: "the GitHub repository e.g. someOwner/someRepo",
					},
					&cli.StringFlag{
						Name:  "start",
						Usage: "the starting, exclusive commit for FilesDiff mode",
					},
					&cli.StringFlag{
						Name:     "end",
						Usage:    "the ending, inclusive commit",
						Required: true,
					},
				},
				Action: detect,
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
}
