package report

import (
	"github.com/andygrunwald/go-jira"
	"sort"
	"time"
	"velojiraptor/internal/output"
)

type HistoryReport struct {
	Changes []Change
}

type Change struct {
	Issue     string
	Field     string
	From      string
	To        string
	ChangedAt time.Time
}

func (hr *HistoryReport) Normalize() output.Grid {
	grid := output.Grid{
		Headers: []string{
			"Issue",
			"Field",
			"From",
			"To",
			"Changed At",
		},
	}

	for _, change := range hr.Changes {
		grid.Add(map[string]string{
			"Issue":      change.Issue,
			"Field":      change.Field,
			"From":       change.From,
			"To":         change.To,
			"Changed At": change.ChangedAt.Format(time.RFC3339),
		})
	}

	return grid
}

func History(issues *[]jira.Issue, field string) HistoryReport {
	hr := HistoryReport{
		Changes: []Change{},
	}

	for _, issue := range *issues {
		for _, history := range issue.Changelog.Histories {
			for _, historyItem := range history.Items {
				if field != historyItem.Field {
					continue
				}

				t, _ := history.CreatedTime()

				hr.Changes = append(hr.Changes, Change{
					issue.Key,
					historyItem.Field,
					historyItem.FromString,
					historyItem.ToString,
					t,
				})
			}
		}
	}

	sort.Slice(hr.Changes, func(i, j int) bool {
		return hr.Changes[i].ChangedAt.Before(hr.Changes[j].ChangedAt)
	})

	return hr
}

func (change *Change) ToArray() []string {
	return []string{
		change.Issue,
		change.Field,
		change.From,
		change.To,
		change.ChangedAt.Format(time.RFC3339),
	}
}

func FilterByIssue(changes []Change, issue string) []Change {
	var results []Change

	for _, change := range changes {
		if issue != change.Issue {
			continue
		}

		results = append(results, change)
	}

	return results
}
