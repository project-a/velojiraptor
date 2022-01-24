package output

type Matrix struct {
	Headers []string
	Rows    [][]string
}

func (m *Matrix) Add(row map[string]string) {
	m.Rows = append(m.Rows, m.normalize(row))
}

func (m *Matrix) normalize(row map[string]string) []string {
	var normalizedRow []string

	for _, header := range m.Headers {
		if _, ok := row[header]; ok {
			normalizedRow = append(normalizedRow, row[header])

			continue
		}

		normalizedRow = append(normalizedRow, "")
	}

	return normalizedRow
}
