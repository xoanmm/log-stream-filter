// Package main contains code tu run dir-cleaner as a CLI command.
package main

import (
	"fmt"
	"github.com/xoanmm/log-stream-filter/pkg/filter"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/urfave/cli/v2"
)

const dateLayout = "01/02/2006 15:04:05"

var version = "1.0.0"
var date = time.Now().Format(time.RFC3339)
var now = time.Now().UTC()
var nowDate = now.Format(dateLayout)
var nowDateLessEightHours = now.Add(-8 * time.Hour).Format(dateLayout)

func main() {
	cmd := buildCLI(&filter.App{})

	if err := cmd.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

// buildCLI creates a CLI app
func buildCLI(app *filter.App) *cli.App {
	d, _ := time.Parse(time.RFC3339, date)
	return &cli.App{
		Name:      "aws-log-filter",
		Usage:     "retrieves all event logs from all streamLogGroup of a specific logGroup of AWS",
		Version:   version,
		Compiled:  d,
		UsageText: "aws-log-filter [--log-group <log-group-name>] [--log-stream-filter <filter>] " +
			"[--log-stream-filter-position <position>]" +
			"[--aws-profile <aws-profile>] [--aws-region <aws-region>] " +
			"[--path <path>] [--start-date <date>] [--end-date <date>]",
		Authors: []*cli.Author{
			{
				Name:  "Xoan Mallon",
				Email: "xoanmallon@gmail.com",
			},
		},
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "log-group",
				Usage:       "log group name on which all logStreams will be obtained and will apply the filtering",
				Value:       "my-lo-group",
				Aliases:     []string{"n"},
			},

			&cli.StringFlag{
				Name:        "log-stream-filter",
				Usage:       "filter to apply on logStreams name to retrieve eventLogs or not",
				Value:       "service-name-1",
				Aliases:     []string{"l"},
			},

			&cli.IntFlag{
				Name:        "log-stream-filter-position",
				Usage:       "position in which to apply the log-stream-filter in the logStreams of the logGroup by splitting by the character / (Example of logStreamGroup: " +
					"log-group/log-stream-group-prefix/ccc7b271-83ee-4487-b8f0-4246ce2d90ad)",
				Value:       1,
				Aliases:     []string{"f"},
			},

			&cli.StringFlag{
				Name:        "aws-profile",
				Usage:       "aws-profile to use for credentials",
				Value:       "my-profile",
				Aliases:     []string{"a"},
			},

			&cli.StringFlag{
				Name:        "aws-region",
				Usage:       "aws region to use for call operations to aws sdk",
				Value:       "us-east-1",
				Aliases:     []string{"r"},
			},

			&cli.StringFlag{
				Name:        "path",
				Usage:       "path where to store the logs",
				Value:       "/tmp/",
				Aliases:     []string{"p"},
			},

			&cli.StringFlag{
				Name:    "start-date",
				Usage:   "filter only from a date specified (UTC format)",
				Value:   nowDateLessEightHours,
				Aliases: []string{"s"},
			},

			&cli.StringFlag{
				Name:    "end-date",
				Usage:   "filter only until a date specified (UTC format)",
				Value:   nowDate,
				Aliases: []string{"e"},
			},
		},
		Action: func(c *cli.Context) error {
			path, _ := filepath.Abs(c.String("path"))
			logGroup := c.String("log-group")

			logsFileGenerated := app.FilterLogs(&filter.Options{
				LogGroup:	logGroup,
				AwsProfile: c.String("aws-profile"),
				AwsRegion:	c.String("aws-region"),
				LogStreamFilter: c.String("log-stream-filter"),
				LogStreamFilterPosition: c.Int("log-stream-filter-position"),
				Path:       path,
				StartDate:	c.String("start-date"),
				EndDate:   c.String("end-date"),
			})
			fmt.Println(len(logsFileGenerated), "files generated for logs of logStreams filtered for logGroup", logGroup)
			for k := range logsFileGenerated {
				fmt.Printf("logStreamName %s Location of file where its logs were stored at %s\n", k, logsFileGenerated[k][0])
			}
			return nil
		},
	}
}