package util

import (
	"a.a/cu/ss_log"
	"a.a/cu/ss_time"
	"a.a/cu/strext"
	"a.a/mp-server/common/constants"
	"a.a/mp-server/common/global"
	custProto "a.a/mp-server/common/proto/cust"
	"a.a/mp-server/common/ss_count"
	"fmt"
	"github.com/tealeg/xlsx"
)

type FileCommonInfo struct {
	FileName          string
	FilePath          string
	Head              string
	Account           string
	Total             string
	USDCnt            string
	USDSum            string
	USDRefundCnt      string
	USDRefundSum      string
	USDRefundFailCnt  string
	USDRefundFailSum  string
	USDRefundBeingCnt string
	USDRefundBeingSum string
	KHRCnt            string
	KHRSum            string
	KHRRefundCnt      string
	KHRRefundSum      string
	KHRRefundFailCnt  string
	KHRRefundFailSum  string
	KHRRefundBeingCnt string
	KHRRefundBeingSum string
	QueryStr          []string
	Titles            []string
	OrderList         []*custProto.BusinessBillData
	RefundOrder       []*custProto.RefundBill
}

func CreateBillFile(f *FileCommonInfo) (errT error) {

	//将结果写入xlsx文件中
	file := xlsx.NewFile()
	sheet, err1 := file.AddSheet("sheet1")
	if err1 != nil {
		ss_log.Error("err1=[%v]", err1)
		return err1
	}

	style := xlsx.NewStyle()
	style.Font = *xlsx.NewFont(11, "宋体")
	style.ApplyFont = true

	//表头：#modernpay待结算/已结算
	row := sheet.AddRow()
	cell := row.AddCell()
	cell.SetStyle(style)
	cell.Value = f.Head

	//账号
	row2 := sheet.AddRow()
	cell2 := row2.AddCell()
	cell2.SetStyle(style)
	cell2.Value = fmt.Sprintf("#账号：%s", f.Account)

	//筛选条件：#起始日期：[2020年08月21日 00:00:00]
	if f.QueryStr != nil {
		row3 := sheet.AddRow()
		for _, v := range f.QueryStr {
			cell3 := row3.AddCell()
			cell3.SetStyle(style)
			cell3.Value = v
		}
	}

	sheet.AddRow()

	//标题
	row4 := sheet.AddRow()
	for _, title := range f.Titles {
		cell4 := row4.AddCell()
		cell4.SetStyle(style)
		cell4.Value = title
	}

	//添加数据
	for _, data := range f.OrderList {
		row := sheet.AddRow()
		if data.CurrencyType == constants.CURRENCY_UP_USD {
			data.Amount = ss_count.Div(data.Amount, "100").StringFixedBank(2)
			data.Fee = ss_count.Div(data.Fee, "100").StringFixedBank(2)
		}
		switch data.OrderStatus {
		case constants.BusinessOrderStatusPending:
			data.OrderStatus = "待支付"
		case constants.BusinessOrderStatusPay:
			data.OrderStatus = "支付成功"
		case constants.BusinessOrderStatusPayTimeOut:
			data.OrderStatus = "支付失败"
		case constants.BusinessOrderStatusRefund:
			data.OrderStatus = "已全额退款"
		case constants.BusinessOrderStatusRebatesRefund:
			data.OrderStatus = "已部分退款"
		}
		dataStr := []string{
			data.OrderNo,
			data.OutOrderNo,
			ss_time.ParseTimeFromPostgres(data.CreateTime, global.Tz).Format(ss_time.DateTimeDashFormat),
			data.CurrencyType,
			data.Amount,
			data.Fee,
			data.AppName,
			data.Subject,
			data.OrderStatus,
		}

		for _, value := range dataStr {
			cell := row.AddCell()
			cell.SetStyle(style)
			cell.Value = value
		}

	}

	sheet.AddRow()

	//合计：3笔
	row = sheet.AddRow()
	cell = row.AddCell()
	cell.SetStyle(style)
	cell.Value = fmt.Sprintf("合计：%v笔", f.Total)

	//美金
	row = sheet.AddRow()
	cell = row.AddCell()
	cell.SetStyle(style)
	cell.Value = fmt.Sprintf("美金：%v笔", strext.ToInt64(f.USDCnt))
	cell = row.AddCell()
	cell.SetStyle(style)
	cell.Value = fmt.Sprintf("共计:＄%v", ss_count.Div(f.USDSum, "100").StringFixedBank(2))

	row = sheet.AddRow()
	cell = row.AddCell()
	cell.SetStyle(style)
	cell.Value = fmt.Sprintf("美金退款：%v笔", strext.ToInt64(f.USDRefundCnt))
	cell = row.AddCell()
	cell.SetStyle(style)
	cell.Value = fmt.Sprintf("共计:＄%v", ss_count.Div(f.USDRefundSum, "100").StringFixedBank(2))

	//瑞尔
	row = sheet.AddRow()
	cell = row.AddCell()
	cell.SetStyle(style)
	cell.Value = fmt.Sprintf("瑞尔：%v笔", strext.ToInt64(f.KHRCnt))
	cell = row.AddCell()
	cell.SetStyle(style)
	cell.Value = fmt.Sprintf("共计:៛%v", ss_count.Div(f.KHRSum, "100").StringFixedBank(2))

	row = sheet.AddRow()
	cell = row.AddCell()
	cell.SetStyle(style)
	cell.Value = fmt.Sprintf("瑞尔退款：%v笔", strext.ToInt64(f.KHRRefundCnt))
	cell = row.AddCell()
	cell.SetStyle(style)
	cell.Value = fmt.Sprintf("共计:៛%v", ss_count.Div(f.KHRRefundSum, "100").StringFixedBank(2))

	//保存xlsx文件
	errSave := file.Save(f.FilePath)
	if errSave != nil {
		ss_log.Error("保存文件失败，errSave=[%v]", errSave)
		return errSave
	}

	//
	ss_log.Info("产生文件成功[%v]", f.FilePath)
	return nil
}

func CreateRefundFile(f *FileCommonInfo) (errT error) {

	//将结果写入xlsx文件中
	file := xlsx.NewFile()
	sheet, err1 := file.AddSheet("sheet1")
	if err1 != nil {
		return err1
	}

	style := xlsx.NewStyle()
	style.Font = *xlsx.NewFont(11, "宋体")
	style.ApplyFont = true

	//表头：#Mango Pay待结算/已结算
	row := sheet.AddRow()
	cell := row.AddCell()
	cell.SetStyle(style)
	cell.Value = f.Head

	//账号
	row2 := sheet.AddRow()
	cell2 := row2.AddCell()
	cell2.SetStyle(style)
	cell2.Value = fmt.Sprintf("#账号：%s", f.Account)

	//筛选条件：#起始日期：[2020年08月21日 00:00:00]
	if f.QueryStr != nil {
		row3 := sheet.AddRow()
		for _, v := range f.QueryStr {
			cell3 := row3.AddCell()
			cell3.SetStyle(style)
			cell3.Value = v
		}
	}

	sheet.AddRow()

	//标题
	row4 := sheet.AddRow()
	for _, title := range f.Titles {
		cell4 := row4.AddCell()
		cell4.SetStyle(style)
		cell4.Value = title
	}

	//添加数据
	for _, data := range f.RefundOrder {
		row := sheet.AddRow()
		if data.CurrencyType == constants.CURRENCY_UP_USD {
			data.Amount = ss_count.Div(data.Amount, "100").StringFixedBank(2)
			data.TransAmount = ss_count.Div(data.TransAmount, "100").StringFixedBank(2)
		}
		switch data.OrderStatus {
		case constants.BusinessRefundStatusPending:
			data.OrderStatus = "处理中"
		case constants.BusinessRefundStatusSuccess:
			data.OrderStatus = "退款成功"
		case constants.BusinessRefundStatusFail:
			data.OrderStatus = "退款失败"
		}
		dataStr := []string{
			data.RefundNo,
			data.TransOrderNo,
			ss_time.ParseTimeFromPostgres(data.CreateTime, global.Tz).Format(ss_time.DateTimeDashFormat),
			data.CurrencyType,
			data.TransAmount,
			data.Amount,
			data.AppName,
			data.Subject,
			data.OrderStatus,
		}

		for _, value := range dataStr {
			cell := row.AddCell()
			cell.SetStyle(style)
			cell.Value = value
		}

	}

	sheet.AddRow()

	//合计：3笔
	row = sheet.AddRow()
	cell = row.AddCell()
	cell.SetStyle(style)
	cell.Value = fmt.Sprintf("合计：%v笔", f.Total)

	//美金
	row = sheet.AddRow()
	cell = row.AddCell()
	cell.SetStyle(style)
	cell.Value = fmt.Sprintf("美金退款成功：%v笔", strext.ToInt64(f.USDRefundCnt))
	cell = row.AddCell()
	cell.SetStyle(style)
	cell.Value = fmt.Sprintf("共计:＄%v", ss_count.Div(f.USDRefundSum, "100").StringFixedBank(2))

	row = sheet.AddRow()
	cell = row.AddCell()
	cell.SetStyle(style)
	cell.Value = fmt.Sprintf("美金退款失败：%v笔", strext.ToInt64(f.USDRefundFailCnt))
	cell = row.AddCell()
	cell.SetStyle(style)
	cell.Value = fmt.Sprintf("共计:＄%v", ss_count.Div(f.USDRefundFailSum, "100").StringFixedBank(2))

	row = sheet.AddRow()
	cell = row.AddCell()
	cell.SetStyle(style)
	cell.Value = fmt.Sprintf("美金退款处理中：%v笔", strext.ToInt64(f.USDRefundBeingCnt))
	cell = row.AddCell()
	cell.SetStyle(style)
	cell.Value = fmt.Sprintf("共计:＄%v", ss_count.Div(f.USDRefundBeingSum, "100").StringFixedBank(2))

	//瑞尔
	row = sheet.AddRow()
	cell = row.AddCell()
	cell.SetStyle(style)
	cell.Value = fmt.Sprintf("瑞尔退款成功：%v笔", strext.ToInt64(f.KHRRefundCnt))
	cell = row.AddCell()
	cell.SetStyle(style)
	cell.Value = fmt.Sprintf("共计:៛%v", ss_count.Div(f.KHRRefundSum, "100").StringFixedBank(2))

	row = sheet.AddRow()
	cell = row.AddCell()
	cell.SetStyle(style)
	cell.Value = fmt.Sprintf("瑞尔退款失败：%v笔", strext.ToInt64(f.KHRRefundFailCnt))
	cell = row.AddCell()
	cell.SetStyle(style)
	cell.Value = fmt.Sprintf("共计:៛%v", ss_count.Div(f.KHRRefundFailSum, "100").StringFixedBank(2))

	row = sheet.AddRow()
	cell = row.AddCell()
	cell.SetStyle(style)
	cell.Value = fmt.Sprintf("瑞尔退款处理中：%v笔", strext.ToInt64(f.KHRRefundBeingCnt))
	cell = row.AddCell()
	cell.SetStyle(style)
	cell.Value = fmt.Sprintf("共计:៛%v", ss_count.Div(f.KHRRefundBeingSum, "100").StringFixedBank(2))

	//保存xlsx文件
	errSave := file.Save(f.FilePath)
	if errSave != nil {
		ss_log.Error("保存文件失败，errSave=[%v]", errSave)
		return errSave
	}

	//
	ss_log.Info("产生文件成功[%v]", f.FilePath)

	return nil
}
