package table

import (
	"github.com/olekukonko/tablewriter"
	"os"
	"velojiraptor/internal/output"
)

type Table struct {
}

func (c *Table) Dump(report output.Report) error {
	grid := report.Normalize()

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader(grid.Headers)
	table.AppendBulk(grid.Rows)
	table.Render()

	return nil
}
