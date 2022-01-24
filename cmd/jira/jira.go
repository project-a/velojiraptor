package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/andygrunwald/go-jira"
	"github.com/urfave/cli/v2"
	"jira_go/internal/output"
	"jira_go/internal/output/csv"
	"jira_go/internal/output/table"
	"jira_go/internal/report"
	"jira_go/internal/service"
	"log"
	"os"
)

func main() {
	app := &cli.App{
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "format",
				Value: "table",
				Usage: "output format",
			},
		},
		Commands: []*cli.Command{
			{
				Name:   "pull",
				Usage:  "pulls done(!) tickets which were updated within the given timeframe",
				Action: pullAction,
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:    "username",
						Usage:   "JIRA username",
						Value:   "",
						EnvVars: []string{"JIRA_USERNAME"},
					},
					&cli.StringFlag{
						Name:    "password",
						Usage:   "JIRA password",
						Value:   "",
						EnvVars: []string{"JIRA_PASSWORD"},
					},
					&cli.TimestampFlag{
						Name:     "from",
						Layout:   "2006-01-02",
						Usage:    "from date --from=2022-01-02",
						Required: true,
					},
					&cli.TimestampFlag{
						Name:     "to",
						Layout:   "2006-01-02",
						Usage:    "to date --to=2022-01-15",
						Required: true,
					},
					&cli.StringSliceFlag{
						Name:     "boards",
						Aliases:  []string{"b"},
						Usage:    "boards to inclode in the query eg.: -b GH -b GHDH",
						Required: true,
					},
				},
			},
			{
				Name:   "history",
				Usage:  "shows the changes of the given field",
				Action: historyAction,
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:     "field",
						Usage:    "field to show the history of, eg.: --field status",
						Required: true,
					},
				},
			},
			{
				Name:    "time-in-status",
				Aliases: []string{"tis"},
				Action:  timeInStatusAction,
				Flags: []cli.Flag{
					&cli.StringSliceFlag{
						Name:    "exclude",
						Aliases: []string{"e"},
						Usage:   "statuses to exclude eg.: -e TODO -e \"In Progress\"",
					},
				},
			},
			{
				Name:    "lead-time",
				Aliases: []string{"lt"},
				Action:  leadTimeAction,
				Flags: []cli.Flag{
					&cli.StringSliceFlag{
						Name:    "exclude",
						Aliases: []string{"e"},
						Usage:   "statuses to exclude eg.: -e TODO -e \"In Progress\"",
					},
				},
			},
		},
	}

	app.EnableBashCompletion = true
	err := app.Run(os.Args)

	if err != nil {
		log.Fatal(err)
	}
}

func foo(target string) output.Output {
	switch target {
	case "csv":
		return &csv.CSV{}
	default:
		return &table.Table{}
	}
}

func load(file string) (*[]jira.Issue, error) {
	issues, err := service.LoadIssuesFromFile(file)

	return issues, err
}

func pullAction(c *cli.Context) error {
	if len(c.String("username")) == 0 {
		return errors.New("username is empty")
	}

	if len(c.String("password")) == 0 {
		return errors.New("password is empty")
	}

	jiraService := service.NewJiraService(c.String("username"), c.String("password"))
	issues, err := jiraService.Pull(
		c.Timestamp("from"),
		c.Timestamp("to"),
		c.StringSlice("boards"),
	)

	if err != nil {
		return err
	}

	res, _ := json.Marshal(issues)

	if err != nil {
		return err
	}

	fmt.Println(string(res))

	return nil
}

func historyAction(c *cli.Context) error {
	issues, err := load("output.json")

	if err != nil {
		return err
	}

	r := report.History(issues, c.String("field"))

	return foo(c.String("target")).Dump(&r)
}

func leadTimeAction(c *cli.Context) error {
	issues, err := load("output.json")

	if err != nil {
		return err
	}

	r := report.LeadTime(issues, c.StringSlice("exclude"))

	return foo(c.String("target")).Dump(&r)
}

func timeInStatusAction(c *cli.Context) error {
	issues, err := load("output.json")

	if err != nil {
		return err
	}

	r := report.TimeInStatus(issues, c.StringSlice("exclude"))

	return foo(c.String("target")).Dump(&r)
}
