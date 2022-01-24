package main

import (
	"encoding/json"
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
	jiraFlags := []cli.Flag{
		&cli.StringFlag{
			Name:     "username",
			Usage:    "JIRA username",
			Value:    "",
			EnvVars:  []string{"JIRA_USERNAME"},
			Required: true,
		},
		&cli.StringFlag{
			Name:     "password",
			Usage:    "JIRA password",
			Value:    "",
			EnvVars:  []string{"JIRA_PASSWORD"},
			Required: true,
		},
		&cli.StringFlag{
			Name:     "url",
			Usage:    "JIRA url",
			Value:    "",
			EnvVars:  []string{"JIRA_URL"},
			Required: true,
		},
		&cli.StringFlag{
			Name:     "jql",
			Usage:    "JQL to filter the tickets",
			Value:    "",
			Required: true,
		},
	}

	excludeFlag := cli.StringSliceFlag{
		Name:    "exclude",
		Aliases: []string{"e"},
		Usage:   "statuses to exclude eg.: -e TODO -e \"In Progress\"",
	}

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
				Name:    "count",
				Usage:   "Counts how many issues match the given JQL",
				Aliases: []string{"c"},
				Action:  countAction,
				Flags:   jiraFlags,
			},
			{
				Name:    "history",
				Usage:   "Lists the changes made to the given field",
				Aliases: []string{"hi"},
				Action:  historyAction,
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:     "field",
						Usage:    "field to show the history of, eg.: --field status",
						Required: true,
					},
				},
			},
			{
				Name:    "lead-time",
				Usage:   "Shows the lead time report in the given format",
				Aliases: []string{"lt"},
				Action:  leadTimeAction,
				Flags:   []cli.Flag{&excludeFlag},
			},
			{
				Name:    "search",
				Usage:   "Searche issues matching the given JQL",
				Aliases: []string{"s"},
				Action:  searchAction,
				Flags:   jiraFlags,
			},
			{
				Name:    "time-in-status",
				Usage:   "Shows the time in status report in the given format",
				Aliases: []string{"tis"},
				Action:  timeInStatusAction,
				Flags:   []cli.Flag{&excludeFlag},
			},
		},
	}

	err := app.Run(os.Args)

	if err != nil {
		log.Fatal(err)
	}
}

func format(format string) output.Output {
	switch format {
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

func countAction(c *cli.Context) error {
	jiraService := service.NewJiraService(
		c.String("username"),
		c.String("password"),
		c.String("url"),
	)

	count, err := jiraService.Count(c.String("jql"))

	if err != nil {
		return err
	}

	r := report.Count{Count: count}

	return format(c.String("format")).Dump(&r)
}

func historyAction(c *cli.Context) error {
	issues, err := load("output.json")

	if err != nil {
		return err
	}

	r := report.History(issues, c.String("field"))

	return format(c.String("format")).Dump(&r)
}

func leadTimeAction(c *cli.Context) error {
	issues, err := load("output.json")

	if err != nil {
		return err
	}

	r := report.LeadTime(issues, c.StringSlice("exclude"))

	return format(c.String("format")).Dump(&r)
}

func searchAction(c *cli.Context) error {
	jiraService := service.NewJiraService(
		c.String("username"),
		c.String("password"),
		c.String("url"),
	)

	issues, err := jiraService.FindIssuesByJQL(c.String("jql"))

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

func timeInStatusAction(c *cli.Context) error {
	issues, err := load("output.json")

	if err != nil {
		return err
	}

	r := report.TimeInStatus(issues, c.StringSlice("exclude"))

	return format(c.String("format")).Dump(&r)
}
