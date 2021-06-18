package transform

import (
	"a.a/mp-server/etl-srv/m"
)

type RemoveFieldsTrans struct {
}

var RemoveFieldsTransInst RemoveFieldsTrans

func (RemoveFieldsTrans) Do(ctx *m.TaskContext) {
	for _, v := range ctx.Transform[ctx.TransformPc].Keys {
		delete(ctx.DataMap, v)
	}
	return
}
