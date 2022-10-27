package report

import (
	"github.com/andygrunwald/go-jira"
	"velojiraptor/internal/output"
)

type HeaderListReport struct {
	UniqueHeaders []string
}

func (h *HeaderListReport) Normalize() output.Grid {
	return output.Grid{
		Headers: h.UniqueHeaders,
	}
}

func HeaderList(issues *[]jira.Issue) HeaderListReport {
	headerMap := make(map[string]bool)
	// creating a map because in golang there are no set
	for _, issue := range *issues {
		headerMap[issue.Fields.Status.Name] = true
	}
	uniqueHeaders := make([]string, 0, len(headerMap))
	for key := range headerMap {
		uniqueHeaders = append(uniqueHeaders, key)
	}
	return HeaderListReport{
		UniqueHeaders: uniqueHeaders,
	}
}
