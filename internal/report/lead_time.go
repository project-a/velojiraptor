package report

import (
	"github.com/andygrunwald/go-jira"
	"github.com/rocketlaunchr/dataframe-go"
	"time"
)

type LeadTimeReport struct {
	Average time.Duration
	Spread  time.Duration
}

func (ltr *LeadTimeReport) Normalize() *dataframe.DataFrame {
	avg := dataframe.NewSeriesFloat64("Average", nil)
	spr := dataframe.NewSeriesFloat64("Spread", nil)

	df := dataframe.NewDataFrame(avg, spr)

	df.Append(nil, map[string]interface{}{
		"Average": ltr.Average.Hours() / 24,
		"Spread":  ltr.Spread.Hours() / 24,
	})

	return df
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
