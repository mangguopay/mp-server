package m

type Task interface {
	Do(ctx *TaskContext)
}

type TaskParam struct {
	Task      Task          // func
	SqlStr    string        // sql
	Dbname    string        // dbname
	Args      []interface{} // args
	Data      []interface{} // data
	IntsParam []int         // param
	DataCnt   int           // data len
	ArgsType  int           // args type
	Keys      []string
}

type TaskContext struct {
	DataMap  map[string]interface{}
	DataList []interface{}

	Extract     *TaskParam
	Transform   []*TaskParam
	TransformPc int
	Load        []*TaskParam
	LoadPc      int
}

type TaskGroup interface {
	Do() *TaskContext
}

const (
	ArgsType_DataList = 1
	ArgsType_DataMap  = 2
)
