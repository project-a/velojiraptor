package report

import (
	"fmt"
	"github.com/andygrunwald/go-jira"
	"jira_go/internal/output"
	"time"
)

type TimeInStatusSummary struct {
	Issue    string
	Statuses map[string]time.Duration
}

type TimeInStatusReport struct {
	UniqueStatuses []string
	Summaries      []TimeInStatusSummary
}

func (tisr *TimeInStatusReport) Normalize() output.Matrix {
	m := output.Matrix{
		Headers: tisr.UniqueStatuses,
	}

	for _, summary := range tisr.Summaries {
		m.Add(summary.ToRow())
	}

	return m
}

func TimeInStatus(issues *[]jira.Issue, excludedStatuses []string) TimeInStatusReport {
	historyReport := History(issues, "status")
	var excludes = make(map[string]struct{})

	for _, ignoredStatus := range excludedStatuses {
		excludes[ignoredStatus] = struct{}{}
	}

	timeInStatusReport := TimeInStatusReport{
		UniqueStatuses: uniqueStatuses(historyReport.Changes, excludes),
		Summaries:      []TimeInStatusSummary{},
	}

	for _, issue := range *issues {
		timeInStatusReport.Summaries = append(
			timeInStatusReport.Summaries,
			summarize(issue, historyReport, excludes),
		)
	}

	return timeInStatusReport
}

func uniqueStatuses(changes []Change, excludes map[string]struct{}) []string {
	statuses := make(map[string]struct{})

	for _, change := range changes {
		if _, ok := excludes[change.From]; ok {
			continue
		}

		statuses[change.From] = struct{}{}
	}

	headers := []string{"Issue"}
	for status := range statuses {
		headers = append(headers, status)
	}

	return headers
}

func summarize(issue jira.Issue, historyReport HistoryReport, exclude map[string]struct{}) TimeInStatusSummary {
	summary := TimeInStatusSummary{
		Issue:    issue.Key,
		Statuses: make(map[string]time.Duration),
	}

	createdAt := time.Time(issue.Fields.Created)
	ticketChanges := FilterByIssue(historyReport.Changes, issue.Key)

	for i, change := range ticketChanges {
		if _, ok := exclude[change.From]; ok {
			continue
		}

		if 0 == i {
			summary.add(change.From, change.ChangedAt.Sub(createdAt))

			continue
		}

		summary.add(change.From, change.ChangedAt.Sub(ticketChanges[i-1].ChangedAt))
	}

	return summary
}

func (summary *TimeInStatusSummary) getDurationByStatus(status string) time.Duration {
	if _, ok := summary.Statuses[status]; ok {
		return summary.Statuses[status]
	}

	return time.Duration(0)
}

func (summary *TimeInStatusSummary) add(status string, duration time.Duration) {
	if _, ok := summary.Statuses[status]; ok {
		summary.Statuses[status] += duration

		return
	}

	summary.Statuses[status] = duration
}

func (summary *TimeInStatusSummary) ToRow() map[string]string {
	m := map[string]string{
		"Issue": summary.Issue,
	}

	for status, duration := range summary.Statuses {
		m[status] = fmt.Sprintf("%.2f", duration.Hours()/24)
	}

	return m
}
