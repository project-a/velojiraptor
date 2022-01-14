package output

type Matrix struct {
	Headers []string
	Rows    [][]string
}

type Normalizable interface {
	ToRow() map[string]string
}

func (m *Matrix) Add(normalizable Normalizable) {
	m.Rows = append(m.Rows, m.normalize(normalizable))
}

func (m *Matrix) normalize(normalizable Normalizable) []string {
	var normalizedRow []string
	row := normalizable.ToRow()

	for _, header := range m.Headers {
		if _, ok := row[header]; ok {
			normalizedRow = append(normalizedRow, row[header])

			continue
		}

		normalizedRow = append(normalizedRow, "")
	}

	return normalizedRow
}
