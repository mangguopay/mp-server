package util

import (
	"errors"

	custProto "a.a/mp-server/common/proto/cust"

	"a.a/cu/ss_log"
	"a.a/mp-server/common/constants"
	"github.com/tealeg/xlsx"
)

func CreateXlsxFile(fileName, billType, queryStr string, datas *custProto.XlsxFileContentData) (errT error) {
	//将结果写入xlsx文件中
	file := xlsx.NewFile()
	sheet, err1 := file.AddSheet("sheet1")
	if err1 != nil {
		ss_log.Error("err1=[%v]", err1)
		return err1
	}

	//添加筛选条件到xlsx文件中
	row := sheet.AddRow()
	row.SetHeightCM(1) //设置每行的高度
	cell := row.AddCell()
	cell.Value = queryStr

	switch billType {
	case constants.XLSX_FILE_TYPE_EXCHANGE:
		//添加标题
		row = sheet.AddRow()
		row.SetHeightCM(1) //设置每行的高度
		for _, title := range []string{
			"兑换订单流水号",
			"发起兑换的手机号",
			"发起金额类型",
			"发起金额",
			"到账金额类型",
			"到账金额",
			"创建时间",
			"平台汇率(%)",
			"订单状态",
			"来源",
			"完成时间",
			"手续费",
			"原因",
		} {
			cell := row.AddCell()
			cell.Value = title
		}

		//添加数据
		for _, data := range datas.ExchangeDatas {
			row := sheet.AddRow()
			row.SetHeightCM(1) //设置每行的高度

			dataStr := []string{
				data.LogNo,
				data.Phone,
				data.InType,
				data.Amount,
				data.OutType,
				data.TransAmount,
				data.CreateTime,
				data.Rate,
				data.OrderStatus,
				data.TransFrom,
				data.FinishTime,
				data.Fees,
				data.ErrReason,
			}

			for _, value := range dataStr {
				cell := row.AddCell()
				cell.Value = value
			}

		}
	case constants.XLSX_FILE_TYPE_INCOME:
		//添加标题
		row = sheet.AddRow()
		row.SetHeightCM(1) //设置每行的高度
		for _, title := range []string{
			"存款订单流水号",
			"存款手机号",
			"收款手机号",
			"币种",
			"金额",
			"手续费",
			"订单状态",
			"发起时间",
			"结束时间",
			"收款的服务商账号",
			"核销码",
		} {
			cell := row.AddCell()
			cell.Value = title
		}

		//添加数据
		for _, data := range datas.IncomeDatas {
			row := sheet.AddRow()
			row.SetHeightCM(1) //设置每行的高度

			dataStr := []string{
				data.LogNo,
				data.IncomePhone,
				data.RecvPhone,
				data.BalanceType,
				data.Amount,
				data.Fees,
				data.OrderStatus,
				data.CreateTime,
				data.FinishTime,
				data.Account,
				data.WriteOff,
			}

			for _, value := range dataStr {
				cell := row.AddCell()
				cell.Value = value
			}
		}
	case constants.XLSX_FILE_TYPE_OUTGO:
		//添加标题
		row = sheet.AddRow()
		row.SetHeightCM(1) //设置每行的高度

		for _, title := range []string{
			"取款订单流水号",
			"取款手机号",
			"币种",
			"金额",
			"手续费",
			"订单状态",
			"发起时间",
			"结束时间",
			"服务商账号",
			"核销码",
		} {
			cell := row.AddCell()
			cell.Value = title
		}

		//添加数据
		for _, data := range datas.OutgoDatas {
			row := sheet.AddRow()
			row.SetHeightCM(1) //设置每行的高度

			dataStr := []string{
				data.LogNo,
				data.Phone,
				data.BalanceType,
				data.Amount,
				data.Fees,
				data.OrderStatus,
				data.CreateTime,
				data.FinishTime,
				data.Account,
				data.WriteOff,
			}

			for _, value := range dataStr {
				cell := row.AddCell()
				cell.Value = value
			}

		}
	case constants.XLSX_FILE_TYPE_TRANSFER:
		//添加标题
		row = sheet.AddRow()
		row.SetHeightCM(1) //设置每行的高度

		for _, title := range []string{
			"转账订单流水号",
			"付款账号",
			"收款账号",
			"币种",
			"金额",
			"手续费",
			"订单状态",
			"发起时间",
			"结束时间",
			"核销码",
		} {
			cell := row.AddCell()
			cell.Value = title
		}

		//添加数据
		for _, data := range datas.TransferDatas {
			row := sheet.AddRow()
			row.SetHeightCM(1) //设置每行的高度

			dataStr := []string{
				data.LogNo,
				data.FromAccount,
				data.ToAccount,
				data.BalanceType,
				data.Amount,
				data.Fees,
				data.OrderStatus,
				data.CreateTime,
				data.FinishTime,
				data.WriteOff,
			}

			for _, value := range dataStr {
				cell := row.AddCell()
				cell.Value = value
			}

		}
	case constants.XLSX_FILE_TYPE_COLLECTION:
		//添加标题
		row = sheet.AddRow()
		row.SetHeightCM(1) //设置每行的高度

		//添加标题
		for _, title := range []string{
			"收款订单流水号",
			"付款账号",
			"收款账号",
			"币种",
			"金额",
			"手续费",
			"订单状态",
			"发起时间",
			"结束时间",
		} {
			cell := row.AddCell()
			cell.Value = title
		}

		//添加数据
		for _, data := range datas.CollectionDatas {
			row := sheet.AddRow()
			row.SetHeightCM(1) //设置每行的高度

			dataStr := []string{
				data.LogNo,
				data.FromAccount,
				data.ToAccount,
				data.BalanceType,
				data.Amount,
				data.Fees,
				data.OrderStatus,
				data.CreateTime,
				data.FinishTime,
			}

			for _, value := range dataStr {
				cell := row.AddCell()
				cell.Value = value
			}

		}
	case constants.XLSX_FILE_TYPE_TO_HEADQUARTERS:
		//添加标题
		row = sheet.AddRow()
		row.SetHeightCM(1) //设置每行的高度

		//添加标题
		for _, title := range []string{
			"服务商充值订单流水号",
			"服务商账号",
			"币种",
			"金额",
			"收款方式",
			"收款人/收款卡号",
			"发起时间",
			"订单类型",
			"订单状态",
		} {
			cell := row.AddCell()
			cell.Value = title
		}

		//添加数据
		for _, data := range datas.ToHeadquartersDatas {
			row := sheet.AddRow()
			row.SetHeightCM(1) //设置每行的高度

			dataStr := []string{
				data.LogNo,
				data.Account,
				data.CurrencyType,
				data.Amount,
				data.CollectionType,

				data.Name + "-" + data.CardNumber,
				data.CreateTime,
				data.OrderType,
				data.OrderStatus,
			}

			for _, value := range dataStr {
				cell := row.AddCell()
				cell.Value = value
			}

		}
	case constants.XLSX_FILE_TYPE_TO_SERVICER:
		//添加标题
		row = sheet.AddRow()
		row.SetHeightCM(1) //设置每行的高度

		//添加标题
		for _, title := range []string{
			"服务商请款订单流水号",
			"服务商账号",
			"服务商昵称",
			"币种",
			"金额",
			"收款方式",
			"收款人-收款卡号",
			"发起时间",
			"订单类型",
			"订单状态",
		} {
			cell := row.AddCell()
			cell.Value = title
		}

		//添加数据
		for _, data := range datas.ToServicerDatas {
			row := sheet.AddRow()
			row.SetHeightCM(1) //设置每行的高度

			dataStr := []string{
				data.LogNo,
				data.Account,
				data.Nickname,
				data.CurrencyType,
				data.Amount,
				data.CollectionType,

				data.Name + "_" + data.CardNumber,
				data.CreateTime,
				data.OrderType,
				data.OrderStatus,
			}

			for _, value := range dataStr {
				cell := row.AddCell()
				cell.Value = value
			}
		}
	case constants.XLSX_FILE_TYPE_VACCOUNT_LOG:
		//添加标题
		row = sheet.AddRow()
		row.SetHeightCM(1) //设置每行的高度

		//添加标题
		for _, title := range []string{
			"虚拟账户日志流水号",
			"业务流水号",
			"虚拟账户账号",
			"创建时间",
			"币种",
			"变化的金额",
			"操作类型",
			"冻结金额",
			"现有余额",
			"原因",
			"每小时对账单号",
			"每日对账单号",
		} {
			cell := row.AddCell()
			cell.Value = title
		}

		//添加数据
		for _, data := range datas.LogVaccountDatas {
			row := sheet.AddRow()
			row.SetHeightCM(1) //设置每行的高度

			dataStr := []string{
				data.LogNo,
				data.BizLogNo,
				data.Nickname,
				data.CreateTime,
				data.BalanceType,
				data.Amount,

				data.OpType,
				data.FrozenBalance,
				data.Balance,
				data.Reason,
				data.SettleHourlyLogNo,
				data.SettleDailyLogNo,
			}

			for _, value := range dataStr {
				cell := row.AddCell()
				cell.Value = value
			}

		}

	default:
		ss_log.Error("订单类型错误[%v]", billType)
		return errors.New("订单类型错误")
	}

	//保存xlsx文件
	//errSave := file.Save(pathStr + "/" + xlsxTaskNo + ".xlsx")
	errSave := file.Save(fileName)
	if errSave != nil {
		ss_log.Error("保存文件失败，errSave=[%v]", errSave)
		return errSave
	}

	//
	ss_log.Info("产生文件成功[%v]", fileName)
	return nil
}
