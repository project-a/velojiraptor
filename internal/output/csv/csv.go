package csv

import (
	"context"
	"github.com/rocketlaunchr/dataframe-go"
	"github.com/rocketlaunchr/dataframe-go/exports"
	"os"
	"velojiraptor/internal/output"
)

type CSV struct {
}

func (csv *CSV) Dump(report output.Report) error {
	nullString := ""

	options := exports.CSVExportOptions{
		NullString: &nullString,
		Range:      dataframe.Range{},
		Separator:  ',',
		UseCRLF:    false,
	}

	return exports.ExportToCSV(context.TODO(), os.Stdout, report.Normalize(), options)
}
