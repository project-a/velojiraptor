package table

import (
	"github.com/olekukonko/tablewriter"
	"jira_go/internal/output"
	"os"
)

type Table struct {
}

func (c *Table) Dump(r output.Report) error {
	m := r.Normalize()

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader(m.Headers)
	table.AppendBulk(m.Rows)
	table.Render()

	return nil
}
