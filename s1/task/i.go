package task

type TaskContext struct {
	AccountNo string
}

type ITask interface {
	Init(ctx *TaskContext)
	Do(ctx *TaskContext)
	Next(ctx *TaskContext)
	IsSop(ctx *TaskContext) bool
}

var (
	TaskInst Task
)

type Task struct {
	Pipeline map[string][]ITask
}

func (r Task) Init() {
	r.Pipeline = make(map[string][]ITask, 1)
}

func (r Task) Run(tag string, ctx *TaskContext) {
	if ctx == nil {
		ctx = new(TaskContext)
	}
	for _, v := range r.Pipeline[tag] {
		v.Init(ctx)
		v.Do(ctx)
		if v.IsSop(ctx) {
			break
		}
		v.Next(ctx)
	}
}
