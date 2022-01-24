package report

import (
	"fmt"
	"jira_go/internal/output"
)

type Count struct {
	Count int
}

func (c *Count) ToRow() map[string]string {
	return map[string]string{"Count": fmt.Sprintf("%d", c.Count)}
}

func (c *Count) Normalize() output.Matrix {
	m := output.Matrix{
		Headers: []string{"Count"},
	}

	m.Add(c)

	return m
}
