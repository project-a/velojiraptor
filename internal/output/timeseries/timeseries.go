package timeseries

import (
	"encoding/json"
	"fmt"
	"velojiraptor/internal/output"
)

type Timeseries struct {
	Timestamps    []interface{}
	NumericValues map[string][]interface{}
	TagValues     map[string][]interface{}
}

func (c *Timeseries) Dump(report output.Report) error {
	df := report.Normalize()
	df.Lock()

	ts := Timeseries{
		Timestamps:    []interface{}{},
		NumericValues: make(map[string][]interface{}),
		TagValues:     make(map[string][]interface{}),
	}

	for _, s := range df.Series {
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

		if "Timestamp" == s.Name() {
			ts.Timestamps = values

			continue
		}

		ts.NumericValues[s.Name()] = values
	}

	df.Unlock()

	data, _ := json.Marshal(ts)
	fmt.Println(string(data))

	return nil
}

//func unixMicro(timestamps []interface{}) []int64 {
//	var unixMicros []int64
//
//	for _, timestamp := range timestamps {
//		unixMicros = append(unixMicros, timestamp)
//	}
//
//	return unixMicros
//}
