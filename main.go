package main

import (
	"bytes"
	"fmt"
	"log"
	"os"

	"github.com/spf13/afero"
	"github.com/urfave/cli/v3"
	prettier "github.com/utopia-planitia/chart-prettier/pkg"
)

func main() {
	err := run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

func run(args []string) error {
	app := &cli.App{
		Name:  "chart-prettier",
		Usage: "sort files in chart directories",
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:    "truncate",
				Aliases: []string{"t"},
				Usage:   "only use stdin and remove all existing manifest files",
			},
			&cli.BoolFlag{
				Name:    "verbose",
				Aliases: []string{"v"},
				Usage:   "show progress",
			},
			&cli.BoolFlag{
				Name:        "stdin",
				Aliases:     []string{"i"},
				Usage:       "read from stdin",
				DefaultText: "auto",
			},
		},
		Action: cleanupChart,
	}

	err := app.Run(args)
	if err != nil {
		return err
	}

	return nil
}

func cleanupChart(c *cli.Context) error {
	verbose := c.Bool("verbose")

	stdin := c.Bool("stdin")
	var err error

	if !c.IsSet("stdin") {
		stdin, err = detectStdinPipe()
		if err != nil {
			return fmt.Errorf("detect pipe to stdin: %v", err)
		}
	}

	if stdin && c.Args().Len() != 1 {
		return fmt.Errorf("using stdin only supports exactly one chart directory")
	}

	truncate := c.Bool("truncate")

	appFs := afero.NewOsFs()

	for _, path := range c.Args().Slice() {
		chart := prettier.NewChart()

		if !truncate {
			logVerbose(verbose, "reading existing files")

			err := chart.LoadChart(appFs, path)
			if err != nil {
				return fmt.Errorf("loading manifests from existing chart: %v", err)
			}

			logVerbose(verbose, "done")
		}

		if stdin {
			logVerbose(verbose, "reading stdin")

			buf := bytes.Buffer{}

			_, err := buf.ReadFrom(os.Stdin)
			if err != nil {
				return fmt.Errorf("readming from stdin: %v", err)
			}

			err = chart.AddManifests(buf.String())
			if err != nil {
				return fmt.Errorf("loading manifests from stdin: %v", err)
			}

			logVerbose(verbose, "done")
		}

		logVerbose(verbose, "deleting files")

		err = chart.DeleteFiles(appFs, path)
		if err != nil {
			return fmt.Errorf("cleanup preexisting manifests: %v", err)
		}

		logVerbose(verbose, "done")

		logVerbose(verbose, "writing files")

		err = chart.WriteOut(appFs, path)
		if err != nil {
			return fmt.Errorf("create new manifests in chart: %v", err)
		}

		logVerbose(verbose, "done")
	}

	return nil
}

func detectStdinPipe() (bool, error) {
	stat, err := os.Stdin.Stat()
	if err != nil {
		return false, err
	}

	pipedToStdin := (stat.Mode() & os.ModeCharDevice) == 0

	return pipedToStdin, nil
}

func logVerbose(verbose bool, message string) {
	if !verbose {
		return
	}

	log.Println(message)
}
