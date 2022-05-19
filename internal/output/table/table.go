package table

import (
	"fmt"
	"velojiraptor/internal/output"
)

type Table struct {
}

func (c *Table) Dump(report output.Report) error {
	fmt.Printf(report.Normalize().Table())

	return nil
}
