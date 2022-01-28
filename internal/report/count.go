package report

import (
	"fmt"
	"velojiraptor/internal/output"
)

type Count struct {
	Count int
}

func (c *Count) Normalize() output.Grid {
	grid := output.Grid{
		Headers: []string{"Count"},
	}

	grid.Add(map[string]string{"Count": fmt.Sprintf("%d", c.Count)})

	return grid
}
