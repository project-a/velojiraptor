package service

import (
	"encoding/json"
	"fmt"
	"github.com/andygrunwald/go-jira"
	"os"
	"strings"
	"time"
)

const DateLayout = "2006-01-02"

type JiraService struct {
	client jira.Client
}

func NewJiraService(username string, password string) *JiraService {
	// @TODO use env vars here
	tp := jira.BasicAuthTransport{
		Password: password,
		Username: username,
	}

	c, _ := jira.NewClient(tp.Client(), "https://gartenhaus-gmbh.atlassian.net")

	return &JiraService{client: *c}
}

// Query will implement pagination of api and get all the issues.
// Jira API has limitation as to maxResults it can return at one time.
// You may have usecase where you need to get all the issues according to jql
// This is where this example comes in.
func (js *JiraService) Query(searchString string) ([]jira.Issue, error) {
	last := 0
	var issues []jira.Issue
	for {
		opt := &jira.SearchOptions{
			MaxResults: 1000, // Max results can go up to 1000
			StartAt:    last,
			Expand:     "changelog",
		}

		chunk, resp, err := js.client.Issue.Search(searchString, opt)

		if err != nil {
			return nil, err
		}

		total := resp.Total

		if issues == nil {
			issues = make([]jira.Issue, 0, total)
		}

		issues = append(issues, chunk...)
		last = resp.StartAt + len(chunk)

		if last >= total {
			return issues, nil
		}
	}
}

func LoadIssuesFromFile(file string) (*[]jira.Issue, error) {
	f, err := os.Open(file)

	if err != nil {
		return nil, err
	}

	defer func(f *os.File) {
		err := f.Close()

		if err != nil {
			panic(err)
		}
	}(f)

	issues := new([]jira.Issue)
	err = json.NewDecoder(f).Decode(issues)

	return issues, err
}

func (js *JiraService) Pull(from *time.Time, to *time.Time, boards []string) (*[]jira.Issue, error) {
	jql := fmt.Sprintf(
		"project IN (%s) AND updated > %s AND updated < %s AND statusCategory IN (Done)",
		strings.Join(boards, ", "),
		from.Format(DateLayout),
		to.Format(DateLayout),
	)

	issues, err := js.Query(jql)

	return &issues, err
}
