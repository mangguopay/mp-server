package task

import (
	"a.a/mp-server/common/constants"
	"a.a/mp-server/s1/p"
)

var (
	SaveMoneyTaskInst SaveMoneyTask
)

type SaveMoneyTask struct {
}

func (SaveMoneyTask) Init(ctx *TaskContext) {

}
func (SaveMoneyTask) Do(ctx *TaskContext) {
	p.Save("233", "244", "1", constants.CURRENCY_USD)
}
func (SaveMoneyTask) Next(ctx *TaskContext) {

}
func (SaveMoneyTask) IsSop(ctx *TaskContext) bool {
	return false
}
