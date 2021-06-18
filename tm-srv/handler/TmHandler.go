package handler

import (
	"a.a/cu/ss_log"
	"a.a/mp-server/common/constants"
	go_micro_srv_tm "a.a/mp-server/common/proto/tm"
	"context"
	"fmt"
	"github.com/micro/go-micro/v2/server"
	"net"
)

type TmHandler struct {
	Server server.Server
}

// 获取服务的端口
func (t *TmHandler) getServerPort() string {
	_, port, _ := net.SplitHostPort(t.Server.Options().Address)
	return port
}

// 获取事物的前缀
func (t *TmHandler) getTxNoPrefix() string {
	//return fmt.Sprintf("%s-%s###", t.Server.Options().Name, t.Server.Options().Id)
	return fmt.Sprintf("%s%s", t.Server.Options().Id, constants.TMSERVER_TX_ID_SEPARATOR)
}

// 启动一个事务
func (t *TmHandler) TxBegin(ctx context.Context, req *go_micro_srv_tm.TxBeginRequest, reply *go_micro_srv_tm.TxBeginReply) error {
	//ss_log.Info("------>BeginTx-req")

	if req.FromServerId == "" { // 来源服务id为空
		ss_log.Error("来源服务id为空")
		reply.Err = "来源服务id为空"
		return nil
	}

	// 创建并启动事务
	txNo, err := createTx(constants.DB_CRM, t.getTxNoPrefix(), req.FromServerId)
	if err != nil {
		ss_log.Error("创建事务失败, err:%v", err.Error())
		reply.Err = err.Error()
		return nil
	}

	reply.TxNo = txNo
	return nil
}

// 单行记录查询
func (t *TmHandler) TxQueryRow(ctx context.Context, req *go_micro_srv_tm.TxQueryRowRequest, reply *go_micro_srv_tm.TxQueryRowReply) error {
	//ss_log.Info("------>QueryRowTx-req, TxNo:%s, Sql:%s, Args:%v", req.TxNo, req.Sql, req.Args)

	if req.FromServerId == "" { // 来源服务id为空
		ss_log.Error("%s:来源服务id为空", req.TxNo)
		reply.Err = "来源服务id为空"
		return nil
	}

	if req.TxNo == "" { // 事务id为空
		ss_log.Error("%s:事务id为空", req.TxNo)
		reply.Err = "事务id为空"
		return nil
	}

	// 事务id进行返回
	reply.TxNo = req.TxNo

	// 获取对应全局事务的处理channel
	th, exist := txManager.Get(req.TxNo)
	if !exist {
		ss_log.Error("%s:事务id不正确,全局事务不存在", req.TxNo)
		reply.Err = "事务id不正确"
		return nil
	}

	// 将args参数转换为interface类型
	args := []interface{}{}
	for _, v := range req.Args {
		args = append(args, interface{}(v))
	}

	// 将sql和参数传入goroutine进行执行
	th.ExecCh <- execData{Sql: req.Sql, Args: args, SqlType: SqlTypeQueryRow, FromServerId: req.FromServerId}

	// 等待执行结果
	res := <-th.WaitCh
	if res.Err != nil {
		ss_log.Error("%s:执行失败, err:%v", req.TxNo, res.Err)
		reply.Err = res.Err.Error()
		return nil
	}

	columnLen := len(res.Data.Columns) // 列的数量
	data := make(map[string]string)

	ss_log.Info("res.Data.Columns:%v", res.Data.Columns)
	ss_log.Info("res.Data.Rows:%v", res.Data.Rows)

	if len(res.Data.Rows) > 0 && len(res.Data.Rows[0]) > 0 { // 取第一条记录
		if columnLen == len(res.Data.Rows[0]) { // 列的个数和查询的数据个数是一样的
			for i := 0; i < columnLen; i++ {
				data[res.Data.Columns[i]] = res.Data.Rows[0][i]
			}
		}
	}

	ss_log.Info("data--------------->:%v", data)

	reply.Datas = data
	ss_log.Info("reply.Datas--------------->:%v", reply)
	return nil
}

// 多行记录查询
func (t *TmHandler) TxQueryRows(ctx context.Context, req *go_micro_srv_tm.TxQueryRowsRequest, reply *go_micro_srv_tm.TxQueryRowsReply) error {
	//ss_log.Info("------>QueryRowsTx-req, TxNo:%s, Sql:%s, Args:%v", req.TxNo, req.Sql, req.Args)

	if req.FromServerId == "" { // 来源服务id为空
		ss_log.Error("%s:来源服务id为空", req.TxNo)
		reply.Err = "来源服务id为空"
		return nil
	}

	if req.TxNo == "" { // 事务id为空
		ss_log.Error("%s:事务id为空", req.TxNo)
		reply.Err = "事务id为空"
		return nil
	}

	// 事务id进行返回
	reply.TxNo = req.TxNo

	// 获取对应全局事务的处理channel
	th, exist := txManager.Get(req.TxNo)
	if !exist {
		ss_log.Error("%s:事务id不正确,全局事务不存在", req.TxNo)
		reply.Err = "事务id不正确"
		return nil
	}

	// 将args参数转换为interface类型
	args := []interface{}{}
	for _, v := range req.Args {
		args = append(args, interface{}(v))
	}

	// 将sql和参数传入goroutine进行执行
	th.ExecCh <- execData{Sql: req.Sql, Args: args, SqlType: SqlTypeQueryRows, FromServerId: req.FromServerId}

	// 等待执行结果
	res := <-th.WaitCh
	if res.Err != nil {
		ss_log.Error("%s:执行失败, err:%v", req.TxNo, res.Err)
		reply.Err = res.Err.Error()
		return nil
	}

	if len(res.Data.Columns) > 0 {
		reply.Columns = res.Data.Columns
	}

	if len(res.Data.Rows) > 0 {
		for _, r := range res.Data.Rows {
			reply.Rows = append(reply.Rows, &go_micro_srv_tm.TxQueryRowsReply_Row{Data: r})
		}
	}

	return nil
}

// 执行sql
func (t *TmHandler) TxExec(ctx context.Context, req *go_micro_srv_tm.TxExecRequest, reply *go_micro_srv_tm.TxExecReply) error {
	//ss_log.Info("------>ExecTx-req, TxNo:%s, Sql:%s, Args:%v", req.TxNo, req.Sql, req.Args)

	if req.FromServerId == "" { // 来源服务id为空
		ss_log.Error("%s:来源服务id为空", req.TxNo)
		reply.Err = "来源服务id为空"
		return nil
	}

	if req.TxNo == "" { // 事务id为空
		ss_log.Error("%s:事务id为空", req.TxNo)
		reply.Err = "事务id为空"
		return nil
	}

	// 事务id进行返回
	reply.TxNo = req.TxNo

	// 获取对应全局事务的处理channel
	th, exist := txManager.Get(req.TxNo)
	if !exist {
		ss_log.Error("%s:事务id不正确,全局事务不存在", req.TxNo)
		reply.Err = "事务id不正确"
		return nil
	}

	// 将args参数转换为interface类型
	args := []interface{}{}
	for _, v := range req.Args {
		args = append(args, interface{}(v))
	}

	// 将sql和参数传入goroutine进行执行
	th.ExecCh <- execData{Sql: req.Sql, Args: args, SqlType: SqlTypeExec, FromServerId: req.FromServerId}

	// 等待执行结果
	res := <-th.WaitCh
	if res.Err != nil {
		ss_log.Error("%s:执行失败, err:%v", req.TxNo, res.Err)
		reply.Err = res.Err.Error()
		return nil
	}

	return nil
}

// 提交事务
func (t *TmHandler) TxCommit(ctx context.Context, req *go_micro_srv_tm.TxCommitRequest, reply *go_micro_srv_tm.TxCommitReply) error {
	//ss_log.Info("------>CommitTx-req, TxNo:%s", req.TxNo)

	if req.FromServerId == "" { // 来源服务id为空
		ss_log.Error("%s:来源服务id为空", req.TxNo)
		reply.Err = "来源服务id为空"
		return nil
	}

	if req.TxNo == "" { // 事务id为空
		ss_log.Error("%s:事务id为空", req.TxNo)
		reply.Err = "事务id为空"
		return nil
	}

	// 获取对应全局事务的处理channel
	th, exist := txManager.Get(req.TxNo)
	if !exist {
		ss_log.Error("%s:事务id不正确,全局事务不存在", req.TxNo)
		reply.Err = "事务id不正确"
		return nil
	}

	// 传入goroutine进行执行
	th.ExecCh <- execData{Sql: "", Args: nil, SqlType: SqlTypeCommit, FromServerId: req.FromServerId}

	// 等待执行结果
	res := <-th.WaitCh
	if res.Err != nil {
		ss_log.Error("%s:提交失败, err:%v", req.TxNo, res.Err)
		reply.Err = res.Err.Error()
		return nil
	}

	return nil
}

// 回滚事务
func (t *TmHandler) TxRollback(ctx context.Context, req *go_micro_srv_tm.TxRollbackRequest, reply *go_micro_srv_tm.TxRollbackReply) error {
	//ss_log.Info("------>RollbackTx-req, TxNo:%s", req.TxNo)

	if req.FromServerId == "" { // 来源服务id为空
		ss_log.Error("%s:来源服务id为空", req.TxNo)
		reply.Err = "来源服务id为空"
		return nil
	}

	if req.TxNo == "" { // 事务id为空
		ss_log.Error("%s:事务id为空", req.TxNo)
		reply.Err = "事务id为空"
		return nil
	}

	// 获取对应全局事务的处理channel
	th, exist := txManager.Get(req.TxNo)
	if !exist {
		ss_log.Error("%s:事务id不正确,全局事务不存在", req.TxNo)
		reply.Err = "事务id不正确"
		return nil
	}

	// 传入goroutine进行执行
	th.ExecCh <- execData{Sql: "", Args: nil, SqlType: SqlTypeRollback, FromServerId: req.FromServerId}

	// 等待执行结果
	res := <-th.WaitCh
	if res.Err != nil {
		ss_log.Error("%s:回滚失败, err:%v", req.TxNo, res.Err)
		reply.Err = res.Err.Error()
		return nil
	}

	return nil
}
