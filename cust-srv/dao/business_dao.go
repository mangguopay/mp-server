package dao

import (
	"a.a/cu/db"
	"a.a/cu/ss_log"
	"a.a/cu/strext"
	"a.a/mp-server/common/constants"
	"a.a/mp-server/common/model"
	go_micro_srv_cust "a.a/mp-server/common/proto/cust"
	"a.a/mp-server/common/ss_sql"
	"database/sql"
	"errors"
	"fmt"
)

type BusinessDao struct {
}

var BusinessDaoInst BusinessDao

func (BusinessDao) GetBusinessDetail(whereList []*model.WhereSqlCond) (data *go_micro_srv_cust.BusinessData, err error) {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	whereModel := ss_sql.SsSqlFactoryInst.InitWhereList(whereList)
	sqlIn := " SELECT bu.business_no, bu.create_time, bu.use_status, bu.full_name, bu.business_id, bu.business_type, " +
		" bu.main_industry, bi.up_code, bu.main_business, bu.contact_person, bu.contact_phone, bu.contact_phone_country_code, " +
		" bu.income_authorization, acc.business_auth_status " +
		" from business bu " +
		" left join business_industry bi on bi.code = bu.main_industry " +
		" left join account acc on acc.uid = bu.account_no " + whereModel.WhereStr
	rows, stmt, errT := ss_sql.QueryRowN(dbHandler, sqlIn, whereModel.Args...)
	if stmt != nil {
		defer stmt.Close()
	}
	if errT != nil {
		ss_log.Error("err=[%v]", errT)
		return nil, errT
	}

	dataT := &go_micro_srv_cust.BusinessData{}
	var upMainIndustry, authStatus sql.NullString
	errT = rows.Scan(
		&dataT.BusinessNo,
		&dataT.CreateTime,
		&dataT.UseStatus,
		&dataT.FullName,
		&dataT.BusinessId,

		&dataT.BusinessType,
		&dataT.MainIndustry,
		&upMainIndustry,
		&dataT.MainBusiness,
		&dataT.ContactPerson,

		&dataT.ContactPhone,
		&dataT.CountryCode, //联系人电话的国家码
		&dataT.IncomeAuthorization,
		&authStatus,
	)
	if errT != nil {
		ss_log.Error("err=[%v]", errT)
		return nil, errT
	}

	dataT.UpMainIndustry = upMainIndustry.String
	dataT.AuthStatus = authStatus.String

	return dataT, nil
}

func (BusinessDao) AddBusinessTx(tx *sql.Tx, accountNo, fullName, simplifyName, authNumber, businessType string) (err error, businessNo string) {
	//创建运营商
	businessNoT := strext.NewUUID()
	businessId := strext.GetDailyId()
	sqlStr := "insert into business(business_no, account_no, is_delete, use_status, full_name, simplify_name, auth_number, business_id , business_type, create_time) " +
		" values ($1,$2,$3,$4,$5,$6,$7,$8,$9,CURRENT_TIMESTAMP)"
	err = ss_sql.ExecTx(tx, sqlStr, businessNoT, accountNo, "0", "1", fullName, simplifyName, authNumber, businessId, businessType)
	if err != nil {
		ss_log.Error("err=[%v]", err)
	}
	return err, businessNoT
}

func (BusinessDao) UpdateBusinessInfo(accountNo, mainIndustry, mainBusiness, contactPerson, contactPhone, contactPhoneCountryCode string) (err error) {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	sqlStr := "update business set main_industry = $2, main_business = $3, contact_person = $4,contact_phone = $5, contact_phone_country_code = $6 where account_no = $1 "
	errT := ss_sql.Exec(dbHandler, sqlStr, accountNo, mainIndustry, mainBusiness, contactPerson, contactPhone, contactPhoneCountryCode)
	if errT != nil {
		ss_log.Error("err=[%v]", errT)
		return errT
	}

	return nil
}

func (*BusinessDao) GetCnt(whereList []*model.WhereSqlCond) (total string) {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)
	whereModel := ss_sql.SsSqlFactoryInst.InitWhereList(whereList)

	sqlCnt := "select count(1) " +
		" from business bu " +
		" LEFT JOIN business_industry bi ON bi.code = bu.main_industry " +
		" LEFT JOIN account acc ON acc.uid = bu.account_no " +
		" LEFT JOIN vaccount vacc1 ON vacc1.account_no = bu.account_no and vacc1.va_type = " + strext.ToStringNoPoint(constants.VaType_USD_BUSINESS_SETTLED) +
		" LEFT JOIN vaccount vacc2 ON vacc2.account_no = bu.account_no and vacc2.va_type =  " + strext.ToStringNoPoint(constants.VaType_KHR_BUSINESS_SETTLED) +
		" LEFT JOIN vaccount vacc3 ON vacc3.account_no = bu.account_no and vacc3.va_type =  " + strext.ToStringNoPoint(constants.VaType_USD_BUSINESS_UNSETTLED) +
		" LEFT JOIN vaccount vacc4 ON vacc4.account_no = bu.account_no and vacc4.va_type =  " + strext.ToStringNoPoint(constants.VaType_KHR_BUSINESS_UNSETTLED) +
		" " + whereModel.WhereStr
	var totalT sql.NullString
	if err := ss_sql.QueryRow(dbHandler, sqlCnt, []*sql.NullString{&totalT}, whereModel.Args...); err != nil {
		ss_log.Error("err=[%v]", err.Error())
		return "0"
	}

	return totalT.String
}

func (*BusinessDao) GetBusinessList(whereList []*model.WhereSqlCond, page, pageSize int) (datas []*go_micro_srv_cust.BusinessData, err error) {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	whereModel := ss_sql.SsSqlFactoryInst.InitWhereList(whereList)
	ss_sql.SsSqlFactoryInst.AppendWhereExtra(whereModel, " ORDER BY bu.create_time desc ")
	ss_sql.SsSqlFactoryInst.AppendWhereLimit(whereModel, pageSize, page)
	sqlStr := "SELECT bu.business_no, bu.create_time, bu.use_status, bu.full_name, bu.business_id, bu.business_type " +
		", bu.main_industry, bi.up_code, bu.main_business, bu.contact_person, bu.contact_phone, bu.income_authorization, bu.outgo_authorization	" +
		", acc.business_auth_status, acc.email, acc.country_code, acc.business_phone, acc.phone, acc.account, acc.uid " +
		", vacc1.balance, vacc2.balance, vacc3.balance, vacc4.balance " +
		" FROM business bu " +
		" LEFT JOIN business_industry bi ON bi.code = bu.main_industry " +
		" LEFT JOIN account acc ON acc.uid = bu.account_no " +
		" LEFT JOIN vaccount vacc1 ON vacc1.account_no = bu.account_no and vacc1.va_type = " + strext.ToStringNoPoint(constants.VaType_USD_BUSINESS_SETTLED) +
		" LEFT JOIN vaccount vacc2 ON vacc2.account_no = bu.account_no and vacc2.va_type =  " + strext.ToStringNoPoint(constants.VaType_KHR_BUSINESS_SETTLED) +
		" LEFT JOIN vaccount vacc3 ON vacc3.account_no = bu.account_no and vacc3.va_type =  " + strext.ToStringNoPoint(constants.VaType_USD_BUSINESS_UNSETTLED) +
		" LEFT JOIN vaccount vacc4 ON vacc4.account_no = bu.account_no and vacc4.va_type =  " + strext.ToStringNoPoint(constants.VaType_KHR_BUSINESS_UNSETTLED) +
		" " + whereModel.WhereStr
	rows, stmt, err2 := ss_sql.Query(dbHandler, sqlStr, whereModel.Args...)
	if stmt != nil {
		defer stmt.Close()
	}
	defer rows.Close()
	if err2 != nil {
		return nil, err2
	}

	var datasT []*go_micro_srv_cust.BusinessData
	for rows.Next() {
		data := go_micro_srv_cust.BusinessData{}
		var usdBalance, khrBalance sql.NullString
		var usdFrozenBalance, khrFrozenBalance sql.NullString
		var upMainIndustry, countryCode, businessPhone, phone sql.NullString
		err2 = rows.Scan(
			&data.BusinessNo,
			&data.CreateTime,
			&data.UseStatus,
			&data.FullName,
			&data.BusinessId,

			&data.BusinessType,
			&data.MainIndustry,
			&upMainIndustry,
			&data.MainBusiness,
			&data.ContactPerson,

			&data.ContactPhone,
			&data.IncomeAuthorization,
			&data.OutgoAuthorization,

			&data.AuthStatus,
			&data.Email,
			&countryCode,
			&businessPhone,
			&phone,
			&data.Account,
			&data.AccountNo,

			&usdBalance,
			&khrBalance,
			&usdFrozenBalance,
			&khrFrozenBalance,
		)

		if err2 != nil {
			ss_log.Error("err=[%v]", err2)
			continue
		}

		if usdBalance.String == "" {
			usdBalance.String = "0"
		}
		if khrBalance.String == "" {
			khrBalance.String = "0"
		}
		if usdFrozenBalance.String == "" {
			usdFrozenBalance.String = "0"
		}
		if khrFrozenBalance.String == "" {
			khrFrozenBalance.String = "0"
		}

		data.UsdBalance = usdBalance.String
		data.KhrBalance = khrBalance.String
		data.UsdFrozenBalance = usdFrozenBalance.String
		data.KhrFrozenBalance = khrFrozenBalance.String

		data.UpMainIndustry = upMainIndustry.String
		data.CountryCode = countryCode.String

		switch data.BusinessType {
		case constants.AccountType_PersonalBusiness:
			data.Phone = phone.String
		case constants.AccountType_EnterpriseBusiness:
			data.Phone = businessPhone.String
		}

		datasT = append(datasT, &data)
	}

	return datasT, nil
}

func (BusinessDao) UpdateBusinessStatus(businessNo, useStatus string) error {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	sqlStr := "update business set use_status = $2 where business_no = $1 and is_delete = '0' "
	errT := ss_sql.Exec(dbHandler, sqlStr, businessNo, useStatus)
	if errT != nil {
		ss_log.Error("err=[%v]", errT)
		return errT
	}

	return nil
}

//修改收款权限
func (BusinessDao) UpdateBusinessIncomeAuthorizationStatus(businessNo, Status string) error {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	sqlStr := "update business set income_authorization = $2 where business_no = $1 and is_delete = '0' "
	errT := ss_sql.Exec(dbHandler, sqlStr, businessNo, Status)
	if errT != nil {
		ss_log.Error("err=[%v]", errT)
		return errT
	}

	return nil
}

//修改出框权限
func (BusinessDao) UpdateBusinessOutgoAuthorizationStatus(businessNo, Status string) error {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	sqlStr := "update business set outgo_authorization = $2 where business_no = $1 and is_delete = '0' "
	errT := ss_sql.Exec(dbHandler, sqlStr, businessNo, Status)
	if errT != nil {
		ss_log.Error("err=[%v]", errT)
		return errT
	}

	return nil
}

//修改商家全称、简称、注册号/组织机构代码
func (BusinessDao) UpdateBusinessInfoTx(tx *sql.Tx, businessNo, fullName, simplifyName, authNumber string) error {
	sqlStr := "update business set full_name = $2, simplify_name = $3, auth_number = $4 where business_no = $1 and is_delete = '0' "
	if errT := ss_sql.ExecTx(tx, sqlStr, businessNo, fullName, simplifyName, authNumber); errT != nil {
		ss_log.Error("err=[%v]", errT)
		return errT
	}

	return nil
}

func (BusinessDao) UpdateBusinessSimplifyNameTx(tx *sql.Tx, businessNo, simplifyName string) error {
	sqlStr := "update business set simplify_name = $2 where business_no = $1 and is_delete = '0' "
	if errT := ss_sql.ExecTx(tx, sqlStr, businessNo, simplifyName); errT != nil {
		ss_log.Error("err=[%v]", errT)
		return errT
	}

	return nil
}

//获取商家数量
func (BusinessDao) GetBusinessAccountCnt(whereList []*model.WhereSqlCond) (total string, err error) {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	whereModel := ss_sql.SsSqlFactoryInst.InitWhereList(whereList)

	var totalT sql.NullString
	sqlCnt := "select count(1) from business bu " +
		" LEFT JOIN account acc ON acc.uid = bu.account_no " + whereModel.WhereStr
	cntErr := ss_sql.QueryRow(dbHandler, sqlCnt, []*sql.NullString{&totalT}, whereModel.Args...)
	if cntErr != nil {
		ss_log.Error("cntErr=[%v]", cntErr)
	}

	return totalT.String, nil
}

func (BusinessDao) GetBusinessAccount(whereList []*model.WhereSqlCond, page, pageSize int32, sortType string) (datas []*go_micro_srv_cust.BusinessAccountsData, err error) {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	whereModel := ss_sql.SsSqlFactoryInst.InitWhereList(whereList)

	switch sortType {
	case "usd_up": //usd余额正向排序(为null的最上面，余额按升序排，余额相同再按创建时间反序。)
		ss_sql.SsSqlFactoryInst.AppendWhereExtra(whereModel, ` order by vacc.balance is null desc,vacc.balance ASC, acc.create_time desc`)
	case "usd_down": //usd余额反向排序(为null的最下面，余额按降序排，余额相同再按创建时间反序。)
		ss_sql.SsSqlFactoryInst.AppendWhereExtra(whereModel, ` order by vacc.balance is null asc,vacc.balance desc, acc.create_time desc`)
	case "khr_up": //khr余额正向排序
		ss_sql.SsSqlFactoryInst.AppendWhereExtra(whereModel, ` order by vacc2.balance is null desc,vacc2.balance ASC, acc.create_time desc`)
	case "khr_down": //khr余额反向排序
		ss_sql.SsSqlFactoryInst.AppendWhereExtra(whereModel, ` order by vacc2.balance is null asc,vacc2.balance desc, acc.create_time desc`)
	default:
		ss_sql.SsSqlFactoryInst.AppendWhereExtra(whereModel, ` order by acc.create_time desc`)
	}

	//添加limit
	ss_sql.SsSqlFactoryInst.AppendWhereLimitI32(whereModel, pageSize, page)

	sqlStr := "select bu.account_no, bu.business_no, bu.business_type, acc.account, vacc.balance, vacc2.balance " +
		" from business bu " +
		" LEFT JOIN account acc ON acc.uid = bu.account_no " +
		" left join vaccount vacc on vacc.account_no = bu.account_no and vacc.va_type = " + strext.ToStringNoPoint(constants.VaType_USD_BUSINESS_SETTLED) +
		" left join vaccount vacc2 on vacc2.account_no = bu.account_no and vacc2.va_type = " + strext.ToStringNoPoint(constants.VaType_KHR_BUSINESS_SETTLED) + whereModel.WhereStr
	rows3, stmt3, errT := ss_sql.Query(dbHandler, sqlStr, whereModel.Args...)
	if stmt3 != nil {
		defer stmt3.Close()
	}
	defer rows3.Close()

	if errT != nil {
		ss_log.Error("errT=[%v]", errT)
		return nil, errT
	}

	for rows3.Next() {
		data := &go_micro_srv_cust.BusinessAccountsData{}
		var accountNo, businessNo, businessType, account, usdBalance, khrBalance sql.NullString
		errT = rows3.Scan(
			&accountNo,
			&businessNo,
			&businessType,
			&account,
			&usdBalance,
			&khrBalance,
		)
		if errT != nil {
			ss_log.Error("BusinessNo[%v],AccountNo[%v],err=[%v]", data.BusinessNo, data.AccountNo, errT)
			return nil, errT
		}

		data.AccountNo = accountNo.String
		data.BusinessNo = businessNo.String
		data.BusinessType = businessType.String
		data.Account = account.String
		data.UsdBalance = usdBalance.String
		data.KhrBalance = khrBalance.String

		if data.UsdBalance == "" {
			data.UsdBalance = "0"
		}

		if data.KhrBalance == "" {
			data.KhrBalance = "0"
		}

		datas = append(datas, data)
	}

	return datas, nil
}

type BusinessProfit struct {
	BusinessAccount string
	BusinessAccNo   string
	BusinessNo      string
	UsdAmount       string
	KhrAmount       string
	BusinessType    string
}

//查询商家收益
func (BusinessDao) GetBusinessProfit(whereList []*model.WhereSqlCond, page, pageSize int32, opType string, reason []string) ([]*BusinessProfit, error) {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	//操作类型：1(+,入账) 2(-,出账)
	if opType == "" {
		return nil, errors.New("opType can not be empty ")
	}

	//来源
	var reasonStr = ""
	if len(reason) > 1 {
		for k, v := range reason {
			if k == len(reason)-1 {
				reasonStr += fmt.Sprintf("%v", v)
				continue
			}
			reasonStr += fmt.Sprintf("%v, ", v)
		}
	} else if len(reason) == 1 {
		reasonStr += fmt.Sprintf("%v", reason[0])
	}

	//虚账类型：usd, khr已结算虚账
	appendWhere := " and va.va_type in (%v, %v) " +
		"group by acc.account, bu.account_no, bu.business_no "
	appendWhere = fmt.Sprintf(appendWhere, constants.VaType_USD_BUSINESS_SETTLED, constants.VaType_KHR_BUSINESS_SETTLED)

	whereModel := ss_sql.SsSqlFactoryInst.InitWhereList(whereList)
	ss_sql.SsSqlFactoryInst.AppendWhereExtra(whereModel, appendWhere)
	ss_sql.SsSqlFactoryInst.AppendWhereOrderBy(whereModel, "bu.create_time", false)
	ss_sql.SsSqlFactoryInst.AppendWhereLimitI32(whereModel, pageSize, page)
	sqlStr := "select acc.account, bu.account_no, bu.business_no, bu.business_type, sum(lva1.amount) usd, sum(lva2.amount) khr " +
		"from business bu " +
		"left join account acc on acc.uid = bu.account_no " +
		"left join vaccount va on va.account_no = acc.uid " +
		"left join log_vaccount lva1 on lva1.vaccount_no = va.vaccount_no " +
		"and va.balance_type = 'usd' and lva1.op_type = %v and lva1.reason in (%v) " +
		"left join log_vaccount lva2 on lva2.vaccount_no = va.vaccount_no " +
		"and va.balance_type = 'khr' and lva2.op_type = %v and lva2.reason in (%v) "

	sqlStr = fmt.Sprintf(sqlStr, opType, reasonStr, opType, reasonStr)
	sqlStr += whereModel.WhereStr

	rows, stmt, err := ss_sql.Query(dbHandler, sqlStr, whereModel.Args...)
	if err != nil {
		return nil, err
	}
	if stmt != nil {
		defer stmt.Close()
	}
	defer rows.Close()

	var list []*BusinessProfit
	for rows.Next() {
		var account, accountNo, businessNo, businessType, usdAmount, khrAmount sql.NullString
		err := rows.Scan(&account, &accountNo, &businessNo, &businessType, &usdAmount, &khrAmount)
		if err != nil {
			return nil, err
		}
		obj := new(BusinessProfit)
		obj.BusinessAccount = account.String
		obj.BusinessAccNo = accountNo.String
		obj.BusinessNo = businessNo.String
		obj.BusinessType = businessType.String
		obj.UsdAmount = usdAmount.String
		obj.KhrAmount = khrAmount.String
		list = append(list, obj)
	}
	return list, nil
}

func (BusinessDao) GetBusinessNo(accountNo, businessType string) (id, name string, err error) {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	sqlStr := "SELECT business_no, simplify_name FROM business WHERE account_no=$1 AND business_type = $2 AND use_status =$3 "
	var businessNo, simplifyName sql.NullString
	err = ss_sql.QueryRow(dbHandler, sqlStr, []*sql.NullString{&businessNo, &simplifyName}, accountNo, businessType, constants.BusinessUseStatusEnable)
	if err != nil {
		return "", "", err
	}

	return businessNo.String, simplifyName.String, nil
}
