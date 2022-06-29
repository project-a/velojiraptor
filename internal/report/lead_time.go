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
	avg := dataframe.NewSeriesFloat64("Average", nil, ltr.Average.Hours()/24)
	spr := dataframe.NewSeriesFloat64("Spread", nil, ltr.Spread.Hours()/24)
	ts := dataframe.NewSeriesTime("Timestamp", nil, time.Now())

	return dataframe.NewDataFrame(avg, spr, ts)
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
