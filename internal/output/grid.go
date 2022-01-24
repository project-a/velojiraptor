package output

type Grid struct {
	Headers []string
	Rows    [][]string
}

func (grid *Grid) Add(row map[string]string) {
	grid.Rows = append(grid.Rows, grid.normalize(row))
}

func (grid *Grid) normalize(row map[string]string) []string {
	var normalizedRow []string

	for _, header := range grid.Headers {
		if _, ok := row[header]; ok {
			normalizedRow = append(normalizedRow, row[header])

			continue
		}

		normalizedRow = append(normalizedRow, "")
	}

	return normalizedRow
}
