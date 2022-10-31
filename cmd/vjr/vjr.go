package main

import (
	"encoding/json"
	"fmt"
	"github.com/urfave/cli/v2"
	"log"
	"os"
	"velojiraptor/internal/output"
	"velojiraptor/internal/output/csv"
	"velojiraptor/internal/output/table"
	"velojiraptor/internal/report"
	"velojiraptor/internal/service"
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
			Name:     "token",
			Usage:    "JIRA token",
			Value:    "",
			EnvVars:  []string{"JIRA_TOKEN"},
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

	inputFlag := cli.StringFlag{
		Name:     "input",
		Usage:    "Input JSON to process",
		Value:    "",
		Required: true,
	}

	app := &cli.App{
		Usage: "Pulls and generates metrics from Jira",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "format",
				Value: "table",
				Usage: "output format (\"table\", \"csv\")",
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
					&inputFlag,
				},
			},
			{
				Name:    "lead-time",
				Usage:   "Shows the lead time report in the given format",
				Aliases: []string{"lt"},
				Action:  leadTimeAction,
				Flags:   []cli.Flag{&excludeFlag, &inputFlag},
			},
			{
				Name:    "search",
				Usage:   "Search issues matching the given JQL",
				Aliases: []string{"s"},
				Action:  searchAction,
				Flags:   jiraFlags,
			},
			{
				Name:    "time-in-status",
				Usage:   "Shows the time in status report in the given format",
				Aliases: []string{"tis"},
				Action:  timeInStatusAction,
				Flags:   []cli.Flag{&excludeFlag, &inputFlag},
			},
			{
				Name:    "header-list",
				Usage:   "Shows available headers in the given imported file from Jira",
				Aliases: []string{"hl"},
				Action:  getHeaderListAction,
				Flags:   []cli.Flag{&inputFlag},
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

func countAction(c *cli.Context) error {
	jiraService := service.NewJiraService(
		c.String("username"),
		c.String("token"),
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
	issues, err := service.LoadIssuesFromFile(c.String("input"))

	if err != nil {
		return err
	}

	r := report.History(issues, c.String("field"))

	return format(c.String("format")).Dump(&r)
}

func leadTimeAction(c *cli.Context) error {
	issues, err := service.LoadIssuesFromFile(c.String("input"))

	if err != nil {
		return err
	}

	r := report.LeadTime(issues, c.StringSlice("exclude"))

	return format(c.String("format")).Dump(&r)
}

func searchAction(c *cli.Context) error {
	jiraService := service.NewJiraService(
		c.String("username"),
		c.String("token"),
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
	issues, err := service.LoadIssuesFromFile(c.String("input"))

	if err != nil {
		return err
	}

	r := report.TimeInStatus(issues, c.StringSlice("exclude"))

	return format(c.String("format")).Dump(&r)
}

func getHeaderListAction(c *cli.Context) error {
	issues, err := service.LoadIssuesFromFile(c.String("input"))

	if err != nil {
		return err
	}
	r := report.HeaderList(issues)
	for _, header := range r.UniqueHeaders {
		fmt.Println(fmt.Sprintf("\"%s\"", header))
	}
	return nil
}
