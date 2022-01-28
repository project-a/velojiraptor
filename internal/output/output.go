package output

type Output interface {
	Dump(r Report) error
}
