package csv

import (
	"encoding/csv"
	"jira_go/internal/output"
	"os"
)

type CSV struct {
}

func (c *CSV) Dump(report output.Report) error {
	matrix := report.Normalize()

	writer := csv.NewWriter(os.Stdout)
	defer writer.Flush()

	err := writer.Write(matrix.Headers)
	if err != nil {
		return err
	}

	for _, record := range matrix.Rows {
		if err := writer.Write(record); err != nil {
			return err
		}
	}

	return nil
}
