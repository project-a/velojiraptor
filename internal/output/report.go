package output

import "github.com/rocketlaunchr/dataframe-go"

type Report interface {
	Normalize() *dataframe.DataFrame
}
