package handler

import (
	"a.a/cu/db"
	"a.a/cu/ss_log"
	"a.a/cu/ss_time"
	"a.a/cu/strext"
	"a.a/mp-server/common/global"
	"a.a/mp-server/common/ss_sql"
	"a.a/mp-server/tm-srv/dao"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"sync"
	"time"
)

// 管理所有事务的全局变量
var txManager = manager{txList: make(map[string]*txHandler)}

const (
	// 一个事务的总生存时间
	TxMaxLifeTime = time.Second * 600

	// 执行事务的类型
	SqlTypeQueryRow  int = 1 // 查询一条记录sql
	SqlTypeQueryRows int = 2 // 查询多条记录sql
	SqlTypeExec      int = 3 // 执行sql
	SqlTypeCommit    int = 4 // 提交事务
	SqlTypeRollback  int = 5 // 回滚事务

	// 结束事务的方式
	EndModeGetConnectFailed = "get_connect_failed" // 获取连接失败
	EndModeBeginTxFailed    = "begin_tx_failed"    // 开启事务失败
	EndModeCommit           = "commit"             // 提交事务
	EndModeRollback         = "rollback"           // 回滚事务
	EndModeTimeout          = "timeout"            // 事务超时
)

func init() {
	go intervalAddTmTxLog()
}

// 定时添加监控日志
func intervalAddTmTxLog() {
	ticker := time.NewTicker(time.Minute * 3)

	for t := range ticker.C {
		if t.IsZero() {
			// t必须接收，增加一个无效的处理来通过编译
		}

		dao.AddTmTxLog(txManager.Len())
	}
}

// 待执行的sql相关数据
type execData struct {
	Sql          string
	Args         []interface{}
	SqlType      int
	FromServerId string // 请求的服务id
}

type queryRes struct {
	Columns []string
	Rows    [][]string
}

// 待返回的相关数据
type ResData struct {
	Err  error
	Data queryRes // 事务的步骤
}

// 持有单个事务相关数据
type txHandler struct {
	ExecCh chan execData
	WaitCh chan ResData
}

// 管理事务结构体
type manager struct {
	txList map[string]*txHandler
	sync.Mutex
}

// 获取一个事务
func (m *manager) Get(txNo string) (*txHandler, bool) {
	m.Lock()
	defer m.Unlock()
	if th, ok := m.txList[txNo]; ok {
		return th, true
	}
	return nil, false
}

// 创建一个事务
func (m *manager) Create(txNoPrefix string) (string, *txHandler) {
	m.Lock()
	defer m.Unlock()

	// 生成全局事务id
	txNo := strext.GetDailyIdWithPrefix(txNoPrefix)
	th := &txHandler{ExecCh: make(chan execData, 1), WaitCh: make(chan ResData)}

	m.txList[txNo] = th
	return txNo, th
}

// 释放一个事务
func (m *manager) Release(txNo string) {
	m.Lock()
	defer m.Unlock()

	if th, ok := m.txList[txNo]; ok {
		// 释放资源
		close(th.ExecCh)
		close(th.WaitCh)
		// 删除记录
		delete(m.txList, txNo)
	}
}

// 获取未完成的事务数量
func (m *manager) Len() int {
	m.Lock()
	defer m.Unlock()
	return len(m.txList)
}

// 创建并启动一个全局事务
func createTx(dbName string, txNoPrefix string, fromServerId string) (string, error) {
	// 创建一个事务
	txNo, th := txManager.Create(txNoPrefix)

	// 创建一个chan等待事务是否启动成功
	startWait := make(chan error)

	// 开启一个新的goroutine处理全局事务
	go handleTx(dbName, startWait, txNo, th, fromServerId)

	return txNo, <-startWait
}

// 处理事务
func handleTx(dbName string, startWait chan error, txNo string, th *txHandler, fromServerId string) {
	startTime := ss_time.Now(global.Tz)
	endMode := ""
	sqlList := []string{fmt.Sprintf(`{"sqlType":0,"fromServerId":"%s","time":"%s"}`, fromServerId, ss_time.NowForPostgres(global.Tz))}

	defer func() {
		ss_log.Info("%s:事务结束,释放资源", txNo)
		txManager.Release(txNo)

		// 添加事务日志
		go dao.AddTmTimeLog(txNo, startTime, ss_time.Now(global.Tz), endMode, sqlList)
	}()

	// 开启数据库事务
	ss_log.Info("%s:开始开启事务", txNo)

	dbHandler := db.GetDB(dbName)
	if dbHandler == nil {
		ss_log.Error("%s:获取数据库连接失败", txNo)
		startWait <- errors.New("获取数据库连接失败")
		endMode = EndModeGetConnectFailed
		return
	}
	defer db.PutDB(dbName, dbHandler)

	// 如果超时会提示错误: transaction has already been committed or rolled back
	ctx, cancel := context.WithTimeout(context.TODO(), TxMaxLifeTime)
	defer cancel()
	tx, errTx := dbHandler.BeginTx(ctx, nil)

	//tx, errTx := dbHandler.BeginTx(context.TODO(), nil) // 不设置超时时间
	if errTx != nil {
		ss_log.Error("%s:开启数据库事务, err:%v", txNo, errTx)
		startWait <- errors.New("开启事务失败")
		endMode = EndModeBeginTxFailed
		return
	}

	startWait <- nil // 通知调用方，启动事务成功
	ss_log.Info("%s:开启事务成功", txNo)

	tch := time.After(TxMaxLifeTime)
	for {
		select {
		case data := <-th.ExecCh:
			//ss_log.Info("%s, SqlType:%d, Sql:[%s], Args:%v", txNo, data.SqlType, data.Sql, data.Args)

			sqlList = append(sqlList,
				fmt.Sprintf(`{"sqlType":%d,"sql":"%s","args":"%s","fromServerId":"%s","time":"%s"}`,
					data.SqlType, data.Sql, data.Args, data.FromServerId, ss_time.NowForPostgres(global.Tz),
				))

			if data.SqlType == SqlTypeQueryRow || data.SqlType == SqlTypeQueryRows { // 查询sql
				res, err := queryRows(tx, data.Sql, data.Args...)
				if err != nil {
					th.WaitCh <- ResData{Err: err}
				} else {
					th.WaitCh <- ResData{Err: nil, Data: res}
				}
				break
			} else if data.SqlType == SqlTypeExec { // 执行sql
				th.WaitCh <- ResData{Err: ss_sql.ExecTx(tx, data.Sql, data.Args...)}
				break
			} else if data.SqlType == SqlTypeCommit { // 提交事务
				th.WaitCh <- ResData{Err: tx.Commit()} // 提交数据库事务
				endMode = EndModeCommit
				return
			} else if data.SqlType == SqlTypeRollback { // 回滚事务
				th.WaitCh <- ResData{Err: tx.Rollback()}
				endMode = EndModeRollback
				return
			} else {
				th.WaitCh <- ResData{Err: errors.New("事务步骤字段错误")}
				break
			}
			th.WaitCh <- ResData{Err: nil}
		case <-tch: // 超时
			ss_log.Error("%s:超时回滚数事务", txNo)
			if err := tx.Rollback(); err != nil {
				ss_log.Error("%s:超时回滚事务失败, err:%v", txNo, err)
			}
			endMode = EndModeTimeout
			return
		}
	}
}

// 事务查询
func queryRows(tx *sql.Tx, sqlStr string, args ...interface{}) (queryRes, error) {
	rows, stmt, err := ss_sql.QueryTx(tx, sqlStr, args...)

	defer func() {
		if stmt != nil {
			stmt.Close()
		}

		if rows != nil {
			rows.Close()
		}
	}()

	if err != nil {
		return queryRes{}, err
	}

	columns, err := rows.Columns()
	if err != nil {
		return queryRes{}, err
	}

	res := queryRes{Columns: make([]string, 0), Rows: make([][]string, 0)}
	res.Columns = columns

	for rows.Next() {
		values := make([]sql.NullString, len(columns))
		valPtr := make([]interface{}, len(columns))

		for i, _ := range columns {
			valPtr[i] = &values[i]
		}

		err := rows.Scan(valPtr...)
		if err != nil { // 错误示例: converting NULL to string is unsupported
			return queryRes{}, err // 当中间某一行出错，前面scan成功的也忽略
		}

		// 数据转换为字符串
		row := make([]string, len(values))
		for i, v := range values {
			row[i] = v.String
		}
		res.Rows = append(res.Rows, row)
	}

	return res, nil
}
