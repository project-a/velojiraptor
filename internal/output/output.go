package output

type Output interface {
	Dump(report Report) error
}
