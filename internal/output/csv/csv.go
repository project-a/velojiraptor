package csv

import (
	"encoding/csv"
	"jira_go/internal/output"
	"os"
)

type CSV struct {
}

func (c *CSV) Dump(r output.Report) error {
	m := r.Normalize()

	w := csv.NewWriter(os.Stdout)
	defer w.Flush()

	err := w.Write(m.Headers)
	if err != nil {
		return err
	}

	for _, record := range m.Rows {
		if err := w.Write(record); err != nil {
			return err
		}
	}

	return nil
}
