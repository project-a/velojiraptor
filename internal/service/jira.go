package service

import (
	"encoding/json"
	"github.com/andygrunwald/go-jira"
	"os"
)

type JiraService struct {
	client jira.Client
}

func NewJiraService(username string, apiToken string, baseURL string) *JiraService {
	tp := jira.BasicAuthTransport{
		Password: apiToken,
		Username: username,
	}

	c, _ := jira.NewClient(tp.Client(), baseURL)

	return &JiraService{client: *c}
}

// FindIssuesByJQL hadles pagination to get all the issues.
// Jira API has limitation as to maxResults it can return at one time.
func (js *JiraService) FindIssuesByJQL(searchString string) ([]jira.Issue, error) {
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

func (js *JiraService) Count(searchString string) (int, error) {
	opt := &jira.SearchOptions{}

	_, resp, err := js.client.Issue.Search(searchString, opt)

	return resp.Total, err
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
