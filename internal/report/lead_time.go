package report

import (
	"fmt"
	"github.com/andygrunwald/go-jira"
	"time"
	"velojiraptor/internal/output"
)

type LeadTimeReport struct {
	Average time.Duration
	Spread  time.Duration
}

func (ltr *LeadTimeReport) Normalize() output.Grid {
	grid := output.Grid{
		Headers: []string{"Average", "Spread"},
	}

	grid.Add(map[string]string{
		"Average": fmt.Sprintf("%.2f", ltr.Average.Hours()/24),
		"Spread":  fmt.Sprintf("%.2f", ltr.Spread.Hours()/24),
	})

	return grid
}

func LeadTime(issues *[]jira.Issue, excludedStatuses []string) LeadTimeReport {
	timeInStatusReport := TimeInStatus(issues, excludedStatuses)
	min, max, grandTotal := time.Duration(0), time.Duration(0), time.Duration(0)

	for i, record := range timeInStatusReport.Summaries {
		total := time.Duration(0)

		for _, duration := range record.Statuses {
			total += duration
		}

		if min > total || i == 0 {
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
