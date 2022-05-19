package timeseries

import (
	"encoding/json"
	"fmt"
	"time"
	"velojiraptor/internal/output"
)

type Timeseries struct {
	Timestamps    []int64
	NumericValues map[string][]interface{}
	TagValues     map[string][]interface{}
}

func (c *Timeseries) Dump(report output.Report) error {
	df := report.Normalize()
	df.Lock()

	ts := Timeseries{
		Timestamps:    []int64{},
		NumericValues: make(map[string][]interface{}),
		TagValues:     make(map[string][]interface{}),
	}

	now := time.Now().UnixMicro()

	for _, s := range df.Series {
		ts.Timestamps = append(ts.Timestamps, now)
		var values []interface{}
		iterator := s.ValuesIterator()

		for {
			row, val, _ := iterator()

			if row == nil {
				break
			}

			values = append(values, val)
		}

		if "string" == s.Type() {
			ts.TagValues[s.Name()] = values

			continue
		}

		ts.NumericValues[s.Name()] = values
	}

	df.Unlock()

	data, _ := json.Marshal(ts)
	fmt.Println(string(data))

	return nil
}
