package csv

import (
	"encoding/csv"
	"os"
	"velojiraptor/internal/output"
)

type CSV struct {
}

func (c *CSV) Dump(report output.Report) error {
	grid := report.Normalize()

	writer := csv.NewWriter(os.Stdout)
	defer writer.Flush()

	err := writer.Write(grid.Headers)
	if err != nil {
		return err
	}

	for _, record := range grid.Rows {
		if err := writer.Write(record); err != nil {
			return err
		}
	}

	return nil
}
