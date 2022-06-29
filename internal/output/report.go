package output

import "github.com/rocketlaunchr/dataframe-go"

type Report interface {
	Normalize(df *dataframe.DataFrame) *dataframe.DataFrame
	NumericValues(df *dataframe.DataFrame) []dataframe.Series
	Tags(df *dataframe.DataFrame) []dataframe.Series
	Timestamp(df *dataframe.DataFrame) dataframe.Series
}
