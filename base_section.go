package ptvvisum

// BaseSection implements common section functionality
type BaseSection struct {
	name    string
	headers []string
	rows    [][]string
}

func (s *BaseSection) Name() string        { return s.name }
func (s *BaseSection) Headers() []string   { return s.headers }
func (s *BaseSection) Rows() [][]string    { return s.rows }
func (s *BaseSection) AddRow(row []string) { s.rows = append(s.rows, row) }
