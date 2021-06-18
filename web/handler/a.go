package handler

type Op struct {
	Type     string
	Rollback string
	// pipeline
	Steps string
	// if
	Condition string // op_no
	Y         []*Op
	N         []*Op
	// switch
	// -- Condition string // op_no
	Ops map[string][]*Op
}
