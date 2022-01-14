package report

import (
	"fmt"
	"github.com/andygrunwald/go-jira"
	"jira_go/internal/output"
	"time"
)

type LeadTimeReport struct {
	Average time.Duration
	Spread  time.Duration
}

func (ltr *LeadTimeReport) ToRow() map[string]string {
	return map[string]string{
		"Average": fmt.Sprintf("%.2f", ltr.Average.Hours()/24),
		"Spread":  fmt.Sprintf("%.2f", ltr.Spread.Hours()/24),
	}
}

func LeadTime(issues *[]jira.Issue, excludedStatuses []string) LeadTimeReport {
	timeInStatusReport := TimeInStatus(issues, excludedStatuses)
	min, max, grandTotal := time.Duration(0), time.Duration(0), time.Duration(0)

	for _, record := range timeInStatusReport.Summaries {
		total := time.Duration(0)

		for _, duration := range record.Statuses {
			total += duration
		}

		if min > total || min == 0 {
			min = total
		}

		if max < total {
			max = total
		}

		grandTotal += total
	}

	return LeadTimeReport{
		Average: time.Duration(grandTotal.Nanoseconds() / int64(len(timeInStatusReport.Summaries))),
		Spread:  max - min,
	}
}

func (ltr *LeadTimeReport) Normalize() output.Matrix {
	m := output.Matrix{
		Headers: []string{"Average", "Spread"},
	}

	m.Add(ltr)

	return m
}
