package table

import (
	"github.com/olekukonko/tablewriter"
	"jira_go/internal/output"
	"os"
)

type Table struct {
}

func (c *Table) Dump(report output.Report) error {
	matrix := report.Normalize()

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader(matrix.Headers)
	table.AppendBulk(matrix.Rows)
	table.Render()

	return nil
}
