package report

import (
	"github.com/rocketlaunchr/dataframe-go"
	"time"
)

type Count struct {
	Count int
}

func (c *Count) Normalize() *dataframe.DataFrame {
	count := dataframe.NewSeriesFloat64("Count", nil, c.Count)
	ts := dataframe.NewSeriesTime("Timestamp", nil, time.Now())

	return dataframe.NewDataFrame(count, ts)
}

func NumericValues(df *dataframe.DataFrame) []dataframe.Series {
	return
}

func Tags(df *dataframe.DataFrame) []dataframe.Series {

}

func Timestamp(df *dataframe.DataFrame) dataframe.Series {

}
