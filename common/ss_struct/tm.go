package ss_struct

import (
	"context"
	"errors"
	"runtime/debug"
	"strings"
	"time"

	"a.a/cu/ss_log"
	"a.a/mp-server/common/constants"
	tmProto "a.a/mp-server/common/proto/tm"
	"github.com/micro/go-micro/v2/client"
)

// 获取一个新的tm服务
// fromServerId 值示例: go.micro.srv.auth-bd53d88d-9515-401e-85e7-5d3f3b6ab142
func NewTmServerProxy(fromServerId string) (*TmServerProxy, error) {
	return getTmServerProxy(fromServerId, "")
}

// 通过事物id获取对应的tm服务
// fromServerId 值示例: go.micro.srv.auth-bd53d88d-9515-401e-85e7-5d3f3b6ab142
// tmServerNodeId 值示例: 60163dc5-0515-4ae6-ba68-1734aa40ba74
func GetTmServerProxyFromTxNo(fromServerId string, txNo string) (*TmServerProxy, error) {
	arr := strings.Split(txNo, constants.TMSERVER_TX_ID_SEPARATOR)
	if len(arr) != 2 {
		return nil, errors.New("事物txNo参数不正确")
	}

	if arr[0] == "" {
		return nil, errors.New("获取服务节点ID失败")
	}

	p, err := getTmServerProxy(fromServerId, arr[0])
	if err != nil {
		return nil, err
	}

	p.SetTxNo(txNo) // 设置事务id
	return p, nil
}

// 通过tm服务的节点id获取对应的tm服务
// fromServerId 值示例: go.micro.srv.auth-bd53d88d-9515-401e-85e7-5d3f3b6ab142
// tmServerNodeId 值示例: 60163dc5-0515-4ae6-ba68-1734aa40ba74
func getTmServerProxy(fromServerId string, tmServerNodeId string) (*TmServerProxy, error) {
	if fromServerId == "" {
		return nil, errors.New("fromServerId参数不能为空")
	}

	servname := constants.ServerNameTm

	tmClient := &TmClient{Client: client.DefaultClient}

	if err := tmClient.Init(); err != nil {
		return nil, err
	}

	if tmServerNodeId != "" { // 在id的前面拼接上服务名称
		tmServerNodeId = servname + "-" + tmServerNodeId
	}

	if err := tmClient.SelectNode(servname, tmServerNodeId); err != nil {
		return nil, err
	}

	return &TmServerProxy{tmServer: tmProto.NewTmService(servname, tmClient), fromServerId: fromServerId}, nil
}

type TmServerProxy struct {
	tmServer     tmProto.TmService
	fromServerId string
	txNo         string
}

// 开启事务
func (t *TmServerProxy) TxBegin() error {
	rsp, err := t.tmServer.TxBegin(context.TODO(), &tmProto.TxBeginRequest{FromServerId: t.fromServerId})
	if err != nil {
		return err
	}

	if rsp.Err != "" {
		return errors.New(rsp.Err)
	}

	t.SetTxNo(rsp.TxNo) // 设置事务id

	return nil
}

// 提交事务
func (t *TmServerProxy) TxCommit() error {
	rsp, err := t.tmServer.TxCommit(context.TODO(), &tmProto.TxCommitRequest{FromServerId: t.fromServerId, TxNo: t.txNo})
	if err != nil {
		return err
	}

	if rsp.Err != "" {
		return errors.New(rsp.Err)
	}
	return nil
}

// 回滚事务
func (t *TmServerProxy) TxRollback() error {
	rsp, err := t.tmServer.TxRollback(context.TODO(), &tmProto.TxRollbackRequest{FromServerId: t.fromServerId, TxNo: t.txNo})
	if err != nil {
		return err
	}

	if rsp.Err != "" {
		return errors.New(rsp.Err)
	}
	return nil
}

// 获取tm服务对象
func (t *TmServerProxy) GetTmServer() tmProto.TmService {
	return t.tmServer
}

// 获取来源服务id
func (t *TmServerProxy) GetFromServerId() string {
	return t.fromServerId
}

// 获取事务id
func (t *TmServerProxy) SetTxNo(txNo string) {
	t.txNo = txNo
}

// 获取事务id
func (t *TmServerProxy) GetTxNo() string {
	return t.txNo
}

type TmClient struct {
	client.Client
	nodeId   string
	nodeAddr string
}

// PrintNodeInfo--> NodeId:go.micro.srv.tm-60163dc5-0515-4ae6-ba68-1734aa40ba74, NodeAddr:10.41.6.7:5000
func (l *TmClient) PrintNodeInfo() {
	ss_log.Info("PrintNodeInfo--> nodeId:%v, nodeAddr:%v", l.nodeId, l.nodeAddr)
}

// 选择一个服务节点
func (l *TmClient) SelectNode(serverName string, tmServerNodeId string) error {
	//defer l.PrintNodeInfo()

	if tmServerNodeId == "" {
		// 根据底层的算法选择一个节点
		next, err := l.Client.Options().Selector.Select(serverName, l.Client.Options().CallOptions.SelectOptions...)
		if err != nil {
			ss_log.Error("select-err:%v", err)
			return err
		}

		node, nErr := next()
		if nErr != nil {
			ss_log.Error("select-next-err:%v", nErr)
			return nErr
		}

		l.nodeId = node.Id
		l.nodeAddr = node.Address
		return nil
	}

	// 2. 获取指定的服务地址
	serverList, err := l.Client.Options().Registry.GetService(serverName)
	if err != nil {
		ss_log.Error("select-err:%v", err)
		return err
	}

	for _, s := range serverList {
		for _, n := range s.Nodes {
			if n.Id == tmServerNodeId {
				l.nodeId = n.Id
				l.nodeAddr = n.Address
				return nil
			}
		}
	}

	return errors.New("选择服务节点失败,未找到对应的服务节点")
}

func (l *TmClient) Call(ctx context.Context, req client.Request, rsp interface{}, opts ...client.CallOption) error {
	ss_log.Info("[rpc_cli]send=>[%s][%v]", req.Method(), req.Body())
	defer func() {
		if r := recover(); r != nil {
			debug.PrintStack()
			ss_log.Info("[rpc_cli]%s|recovered|r=%+v", req.Method(), r)
		}
	}()

	// 调用options
	var opss []client.CallOption
	opss = append(opss, func(o *client.CallOptions) {
		o.RequestTimeout = time.Second * 30
		o.DialTimeout = time.Second * 30
		o.Retries = 0 // 不进行重试，因为重试的时候是在底层重新选择一个服务，有可能导致请求不在同一台机器了
	})

	if l.nodeAddr != "" {
		//ss_log.Info("[rpc_cli]send-host=>:[%v]", l.nodeAddr)
		opss = append(opss, client.WithAddress(l.nodeAddr))
	}

	opss = append(opss, opts...)

	tb := time.Now()
	err := l.Client.Call(ctx, req, rsp, opss...)
	te := time.Now()
	diff := te.Sub(tb).Milliseconds()
	if err != nil {
		ss_log.Error("[rpc_cli]recv<=|%s|err=[%v]|cost=[%v]ms", req.Method(), err, diff)
	} else {
		ss_log.Info("[rpc_cli]recv<=|%s|[%v]|cost=[%v]ms", req.Method(), rsp, diff)
	}

	return err
}
