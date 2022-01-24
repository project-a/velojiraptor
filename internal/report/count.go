package report

import (
	"fmt"
	"jira_go/internal/output"
)

type Count struct {
	Count int
}

func (c *Count) Normalize() output.Matrix {
	m := output.Matrix{
		Headers: []string{"Count"},
	}

	m.Add(map[string]string{"Count": fmt.Sprintf("%d", c.Count)})

	return m
}
