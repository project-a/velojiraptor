package report

import (
	"github.com/rocketlaunchr/dataframe-go"
)

type Count struct {
	Count int
}

func (c *Count) Normalize() *dataframe.DataFrame {
	return dataframe.NewDataFrame(dataframe.NewSeriesFloat64("Count", nil, c.Count))
}
