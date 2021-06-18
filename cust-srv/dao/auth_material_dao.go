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
)

type AuthMaterialDao struct {
}

//商家认证材料
type BusinessAuthMaterial struct {
	LicenseImgNo string //营业执照图片id
	AuthName     string //公司名称
	AuthNumber   string //注册号/组织机构代码
	AccountUid   string
	StartDate    string //营业期限起始时间

	EndDate      string //营业期限结束时间
	TermType     string //营业期限类型(1-短期，2长期)
	Addr         string //公司地址
	SimplifyName string //简称
	IndustryNo   string //行业类目id
}

var AuthMaterialDaoInst AuthMaterialDao

//=========================用户实名认证============================
func (AuthMaterialDao) GetCnt(whereModelStr string, whereModelArgs []interface{}) (total string, err error) {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	var totalT sql.NullString
	sqlStr := "select count(1) from auth_material am " +
		" left join account acc on acc.uid = am.account_uid " + whereModelStr
	errT := ss_sql.QueryRow(dbHandler, sqlStr, []*sql.NullString{&totalT}, whereModelArgs...)
	if errT != nil {
		ss_log.Error("err=[%v]", errT)
		return "", errT
	}

	return totalT.String, nil
}

func (AuthMaterialDao) GetAuthMaterials(whereModelStr string, whereModelArgs []interface{}) (datas []*go_micro_srv_cust.AuthMaterialData, errR error) {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	sqlStr := "select am.auth_material_no, am.front_img_no, am.back_img_no, am.auth_name, am.auth_number, am.create_time, am.account_uid, am.status" +
		",acc.account " +
		" from auth_material am " +
		" left join account acc on acc.uid = am.account_uid " + whereModelStr

	rows, stmt, err := ss_sql.Query(dbHandler, sqlStr, whereModelArgs...)
	if stmt != nil {
		defer stmt.Close()
	}
	defer rows.Close()
	datasT := []*go_micro_srv_cust.AuthMaterialData{}
	if err != nil {
		ss_log.Error("err=[%v]", err)
		return nil, err
	}
	for rows.Next() {
		data := &go_micro_srv_cust.AuthMaterialData{}
		err = rows.Scan(
			&data.AuthMaterialNo,
			&data.FrontImgNo,
			&data.BackImgNo,
			&data.AuthName,
			&data.AuthNumber,

			&data.CreateTime,
			&data.AccountUid,
			&data.Status,
			&data.Account,
		)
		if err != nil {
			ss_log.Error("err=[%v],AuthMaterialNo=[%v]", err, data.AuthMaterialNo)
			continue
		}
		datasT = append(datasT, data)
	}

	return datasT, nil
}

//根据用户uid查询账号的实名认证信息
func (AuthMaterialDao) GetAuthMaterialDetailByAccountUid(whereModelStr string, whereModelArgs []interface{}) (data *go_micro_srv_cust.AuthMaterialData, errR error) {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	sqlStr := "select am.auth_material_no, am.front_img_no, am.back_img_no, am.auth_name, am.auth_number, am.create_time, am.account_uid, acc.individual_auth_status" +
		" from account acc " +
		" left join auth_material am on am.auth_material_no = acc.individual_auth_material_no " + whereModelStr

	rows, stmt, err := ss_sql.QueryRowN(dbHandler, sqlStr, whereModelArgs...)
	if stmt != nil {
		defer stmt.Close()
	}
	if err != nil {
		ss_log.Error("err=[%v]", err)
		return nil, err
	}
	dataT := &go_micro_srv_cust.AuthMaterialData{}
	var authMaterialNo, frontImgNo, backImgNo, authName, authNumber sql.NullString
	var createTime, accountUid, status sql.NullString
	err = rows.Scan(
		&authMaterialNo,
		&frontImgNo,
		&backImgNo,
		&authName,
		&authNumber,

		&createTime,
		&accountUid,
		&status,
	)
	if err != nil {
		ss_log.Error("err=[%v],AuthMaterialNo=[%v]", err, dataT.AuthMaterialNo)
		return nil, err
	}

	dataT.Status = status.String
	if dataT.Status != constants.AuthMaterialStatus_UnAuth {
		dataT.AuthMaterialNo = authMaterialNo.String
		dataT.FrontImgNo = frontImgNo.String
		dataT.BackImgNo = backImgNo.String
		dataT.AuthName = authName.String
		dataT.AuthNumber = authNumber.String

		dataT.CreateTime = createTime.String
		dataT.AccountUid = accountUid.String
	}

	return dataT, nil
}

func (AuthMaterialDao) ModifyAuthMaterialStatus(tx *sql.Tx, authMaterialNo, status, oldStatus string) (errR error) {
	sqlStr := "update auth_material set status = $2, auth_time=CURRENT_TIMESTAMP where auth_material_no = $1 and status = $3 "
	err := ss_sql.ExecTx(tx, sqlStr, authMaterialNo, status, oldStatus)
	if err != nil {
		ss_log.Error("err=[%v]", err)
		return err
	}
	return nil
}

//审核认证材料后的修改账号认证状态
func (AuthMaterialDao) ModifyAccountIndividualAuthStatus(tx *sql.Tx, uid, status, oldStatus string) (errR error) {
	sqlStr := "update account set individual_auth_status = $2 where uid = $1  and individual_auth_status = $3 and is_delete = '0' "
	err := ss_sql.ExecTx(tx, sqlStr, uid, status, oldStatus)
	if err != nil {
		ss_log.Error("err=[%v]", err)
		return err
	}
	return nil
}

//上传认证材料后的修改账号认证状态和认证资料id
func (AuthMaterialDao) ModifyAccountAuthStatusAndAuthMaterialNo(tx *sql.Tx, uid, status, authMaterialNo string) (errR error) {
	sqlStr := "update account set individual_auth_status = $2, individual_auth_material_no = $3 where uid = $1 and is_delete = '0' "
	err := ss_sql.ExecTx(tx, sqlStr, uid, status, authMaterialNo)
	if err != nil {
		ss_log.Error("err=[%v]", err)
		return err
	}
	return nil
}

//确认账号的个人身份审核认证状态是审核不通过或未认证（只有这两种情况才可以上传认证材料）
func (AuthMaterialDao) CheckAccountAuthStatusByUid(uid string) (total string, errR error) {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	var totalT sql.NullString
	sqlStr := "select count(1) from account where uid = $1 and individual_auth_status in ($2,$3) and is_delete = '0' "
	err := ss_sql.QueryRow(dbHandler, sqlStr, []*sql.NullString{&totalT}, uid, constants.AuthMaterialStatus_Deny, constants.AuthMaterialStatus_UnAuth)
	if err != nil {
		ss_log.Error("err=[%v]", err)
		return "0", err
	}
	return totalT.String, nil
}

//添加个人账号身份认证资料
func (AuthMaterialDao) AddAuthMaterialInfo(tx *sql.Tx, frontImgNo, backImgNo, authName, authNumber, accountUid string) (authMaterialNoR string, errR error) {
	//插入认证资料信息
	authMaterialNo := strext.NewUUID()
	sqlStr := "insert into auth_material(auth_material_no ,front_img_no ,back_img_no ,auth_name ,auth_number ,account_uid , status, create_time) " +
		" values($1,$2,$3,$4,$5,$6,$7,current_timestamp)"
	err := ss_sql.ExecTx(tx, sqlStr, authMaterialNo, frontImgNo, backImgNo, authName, authNumber, accountUid, constants.AuthMaterialStatus_Pending)
	if err != nil {
		ss_log.Error("err=[%v]", err)
		return "", err
	}

	return authMaterialNo, nil
}

//=====================商家认证============================

//获取个人商家认证材料数量
func (AuthMaterialDao) GetBusinessMaterialCnt(whereList []*model.WhereSqlCond) (total string, err error) {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	whereModel := ss_sql.SsSqlFactoryInst.InitWhereList(whereList)

	var totalT sql.NullString
	sqlStr := "select count(1) from auth_material_business amb " +
		" left join account acc on acc.uid = amb.account_uid " + whereModel.WhereStr
	errT := ss_sql.QueryRow(dbHandler, sqlStr, []*sql.NullString{&totalT}, whereModel.Args...)
	if errT != nil {
		ss_log.Error("err=[%v]", errT)
		return "", errT
	}

	return totalT.String, nil
}

func (AuthMaterialDao) GetBusinessMaterials(whereList []*model.WhereSqlCond, page, pageSize int) (datas []*go_micro_srv_cust.AuthMaterialBusinessData, errR error) {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	whereModel := ss_sql.SsSqlFactoryInst.InitWhereList(whereList)

	//添加排序
	ss_sql.SsSqlFactoryInst.AppendWhereExtra(whereModel, "order by case amb.status when "+constants.AuthMaterialStatus_Pending+" then 1 end, amb.create_time desc ")
	ss_sql.SsSqlFactoryInst.AppendWhereLimit(whereModel, pageSize, page)

	sqlStr := "SELECT amb.auth_material_no, amb.license_img_no, amb.auth_name, amb.auth_number, amb.create_time, amb.account_uid, amb.status " +
		", amb.start_date, amb.end_date, amb.term_type, amb.simplify_name, amb.industry_no, bui.name_ch, acc.account " +
		" FROM auth_material_business amb " +
		" LEFT JOIN account acc ON acc.uid = amb.account_uid " +
		" LEFT JOIN business_industry bui ON bui.code = amb.industry_no " + whereModel.WhereStr

	rows, stmt, err := ss_sql.Query(dbHandler, sqlStr, whereModel.Args...)
	if stmt != nil {
		defer stmt.Close()
	}
	defer rows.Close()
	if err != nil {
		ss_log.Error("err=[%v]", err)
		return nil, err
	}
	datasT := []*go_micro_srv_cust.AuthMaterialBusinessData{}
	for rows.Next() {
		data := &go_micro_srv_cust.AuthMaterialBusinessData{}
		var authMaterialNo, licenseImgNo, authName, authNumber, createTime,
			accountUid, status, startDate, endDate, termType,
			simplifyName, industryNo, industryName, account sql.NullString
		err = rows.Scan(
			&authMaterialNo,
			&licenseImgNo,
			&authName,
			&authNumber,
			&createTime,

			&accountUid,
			&status,
			&startDate,
			&endDate,
			&termType,

			&simplifyName,
			&industryNo,
			&industryName,
			&account,
		)
		if err != nil {
			ss_log.Error("err=[%v],AuthMaterialNo=[%v]", err, data.AuthMaterialNo)
			continue
		}

		data.AuthMaterialNo = authMaterialNo.String
		data.LicenseImgNo = licenseImgNo.String
		data.AuthName = authName.String
		data.AuthNumber = authNumber.String
		data.CreateTime = createTime.String

		data.AccountUid = accountUid.String
		data.Status = status.String
		data.StartDate = startDate.String
		data.EndDate = endDate.String
		data.TermType = termType.String

		data.SimplifyName = simplifyName.String
		data.IndustryNo = industryNo.String
		data.IndustryName = industryName.String
		data.Account = account.String

		datasT = append(datasT, data)
	}

	return datasT, nil
}

//获取企业商家认证材料数量
func (AuthMaterialDao) GetEnterpriseMaterialCnt(whereList []*model.WhereSqlCond) (total string, err error) {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	whereModel := ss_sql.SsSqlFactoryInst.InitWhereList(whereList)

	var totalT sql.NullString
	sqlStr := "select count(1) from auth_material_enterprise am " +
		" left join account acc on acc.uid = am.account_uid " + whereModel.WhereStr
	errT := ss_sql.QueryRow(dbHandler, sqlStr, []*sql.NullString{&totalT}, whereModel.Args...)
	if errT != nil {
		ss_log.Error("err=[%v]", errT)
		return "", errT
	}

	return totalT.String, nil
}

func (AuthMaterialDao) GetEnterpriseMaterials(whereList []*model.WhereSqlCond, page, pageSize int) (datas []*go_micro_srv_cust.AuthMaterialEnterpriseData, errR error) {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	whereModel := ss_sql.SsSqlFactoryInst.InitWhereList(whereList)

	//添加排序
	ss_sql.SsSqlFactoryInst.AppendWhereExtra(whereModel, "order by case am.status when "+constants.AuthMaterialStatus_Pending+" then 1 end, am.create_time desc ")
	ss_sql.SsSqlFactoryInst.AppendWhereLimit(whereModel, pageSize, page)

	sqlStr := "select am.auth_material_no, am.license_img_no, am.auth_name, am.auth_number, am.create_time, am.account_uid, am.status " +
		", am.start_date, am.end_date, am.term_type, am.simplify_name, acc.account " +
		" from auth_material_enterprise am" +
		" left join account acc on acc.uid = am.account_uid " + whereModel.WhereStr

	rows, stmt, err := ss_sql.Query(dbHandler, sqlStr, whereModel.Args...)
	if stmt != nil {
		defer stmt.Close()
	}
	defer rows.Close()
	if err != nil {
		ss_log.Error("err=[%v]", err)
		return nil, err
	}
	datasT := []*go_micro_srv_cust.AuthMaterialEnterpriseData{}
	for rows.Next() {
		data := &go_micro_srv_cust.AuthMaterialEnterpriseData{}
		var startDate, endDate sql.NullString
		err = rows.Scan(
			&data.AuthMaterialNo,
			&data.LicenseImgNo,
			&data.AuthName,
			&data.AuthNumber,

			&data.CreateTime,
			&data.AccountUid,
			&data.Status,
			&startDate,
			&endDate,
			&data.TermType,
			&data.SimplifyName,
			&data.Account,
		)
		data.StartDate = startDate.String
		data.EndDate = endDate.String

		if err != nil {
			ss_log.Error("err=[%v],AuthMaterialNo=[%v]", err, data.AuthMaterialNo)
			return nil, err
		}
		datasT = append(datasT, data)
	}

	return datasT, nil
}

//修改个人商家认证材料状态
func (AuthMaterialDao) ModifyAuthMaterialBusinessStatus(tx *sql.Tx, authMaterialNo, status, oldStatus string) (errR error) {
	sqlStr := "update auth_material_business set status = $2,auth_time=CURRENT_TIMESTAMP where auth_material_no = $1 and status = $3 "
	err := ss_sql.ExecTx(tx, sqlStr, authMaterialNo, status, oldStatus)
	if err != nil {
		ss_log.Error("err=[%v]", err)
		return err
	}
	return nil
}

//修改企业商家认证材料状态
func (AuthMaterialDao) ModifyAuthMaterialEnterpriseStatus(tx *sql.Tx, authMaterialNo, status, oldStatus string) (errR error) {
	sqlStr := "update auth_material_enterprise set status = $2, auth_time = CURRENT_TIMESTAMP where auth_material_no = $1 and status = $3 "
	err := ss_sql.ExecTx(tx, sqlStr, authMaterialNo, status, oldStatus)
	if err != nil {
		ss_log.Error("err=[%v]", err)
		return err
	}
	return nil
}

//如果账号的个人商家认证有审核中或审核通过的则不允许再添加
func (AuthMaterialDao) CheckAccountIndividualBusinessAuthStatusByUid(uid string) (total string, errR error) {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	var totalT sql.NullString
	sqlStr := "select count(1) from auth_material_business where account_uid = $1 and status in ($2,$3) "
	err := ss_sql.QueryRow(dbHandler, sqlStr, []*sql.NullString{&totalT}, uid, constants.AuthMaterialStatus_Pending, constants.AuthMaterialStatus_Passed)
	if err != nil {
		ss_log.Error("err=[%v]", err)
		return "-1", err
	}
	return totalT.String, nil
}

//如果账号的企业商家认证有审核中或审核通过的则不允许再添加
func (AuthMaterialDao) CheckAccountBusinessAuthStatusByUid(uid string) (total string, errR error) {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	var totalT sql.NullString
	sqlStr := "select count(1) from auth_material_enterprise where account_uid = $1 and status in ($2,$3) "
	err := ss_sql.QueryRow(dbHandler, sqlStr, []*sql.NullString{&totalT}, uid, constants.AuthMaterialStatus_Pending, constants.AuthMaterialStatus_Passed)
	if err != nil {
		ss_log.Error("err=[%v]", err)
		return "-1", err
	}
	return totalT.String, nil
}

//添加账号个人商家认证资料
func (AuthMaterialDao) AddBusinessAuthMaterial(tx *sql.Tx, data BusinessAuthMaterial) (authMaterialNoR string, errR error) {
	//插入认证资料信息
	authMaterialNo := strext.NewUUID()
	switch data.TermType {
	case constants.TermType_Short: //短期会录入期限开始、结束时间
		//插入认证资料信息
		sqlStr := " insert into auth_material_business(auth_material_no, license_img_no, auth_name, auth_number, account_uid," +
			" status, start_date, end_date, term_type, industry_no," +
			" simplify_name, create_time) "
		sqlStr += " values($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,current_timestamp)"
		err := ss_sql.ExecTx(tx, sqlStr, authMaterialNo, data.LicenseImgNo, data.AuthName, data.AuthNumber, data.AccountUid,
			constants.AuthMaterialStatus_Pending, data.StartDate, data.EndDate, data.TermType, data.IndustryNo,
			data.SimplifyName)
		if err != nil {
			ss_log.Error("err=[%v]", err)
			return "", err
		}
	case constants.TermType_Long: //长期不会录入期限开始、结束时间
		//插入认证资料信息
		sqlStr := " insert into auth_material_business(auth_material_no, license_img_no, auth_name, auth_number, account_uid, " +
			" status, term_type, industry_no, simplify_name, create_time) "
		sqlStr += " values($1,$2,$3,$4,$5,$6,$7,$8,$9,current_timestamp)"
		err := ss_sql.ExecTx(tx, sqlStr, authMaterialNo, data.LicenseImgNo, data.AuthName, data.AuthNumber, data.AccountUid,
			constants.AuthMaterialStatus_Pending, data.TermType, data.IndustryNo, data.SimplifyName)
		if err != nil {
			ss_log.Error("err=[%v]", err)
			return "", err
		}
	default:
		ss_log.Error("TermType参数[%v]错误", data.TermType)
		return "", errors.New("TermType参数错误")
	}

	return authMaterialNo, nil
}

//添加账号企业商家认证资料
func (AuthMaterialDao) AddEnterpriseBusinessAuthMaterial(tx *sql.Tx, data BusinessAuthMaterial) (authMaterialNoR string, errR error) {
	authMaterialNo := strext.NewUUID()
	switch data.TermType {
	case constants.TermType_Short: //短期会录入期限开始、结束时间
		//插入认证资料信息
		sqlStr := " insert into auth_material_enterprise(auth_material_no ,license_img_no ,auth_name ,auth_number ,account_uid , status, start_date, end_date, term_type, simplify_name, create_time) "
		sqlStr += " values($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,current_timestamp)"
		err := ss_sql.ExecTx(tx, sqlStr, authMaterialNo, data.LicenseImgNo, data.AuthName, data.AuthNumber, data.AccountUid, constants.AuthMaterialStatus_Pending, data.StartDate, data.EndDate, data.TermType, data.SimplifyName)
		if err != nil {
			ss_log.Error("err=[%v]", err)
			return "", err
		}
	case constants.TermType_Long: //长期不会录入期限开始、结束时间
		//插入认证资料信息
		sqlStr := " insert into auth_material_enterprise(auth_material_no ,license_img_no ,auth_name ,auth_number ,account_uid , status, term_type, simplify_name, create_time) "
		sqlStr += " values($1,$2,$3,$4,$5,$6,$7,$8,current_timestamp)"
		err := ss_sql.ExecTx(tx, sqlStr, authMaterialNo, data.LicenseImgNo, data.AuthName, data.AuthNumber, data.AccountUid, constants.AuthMaterialStatus_Pending, data.TermType, data.SimplifyName)
		if err != nil {
			ss_log.Error("err=[%v]", err)
			return "", err
		}
	default:
		ss_log.Error("TermType参数[%v]错误", data.TermType)
		return "", errors.New("TermType参数错误")
	}

	return authMaterialNo, nil
}

//查询账号的个人商家认证信息
func (AuthMaterialDao) GetAuthMaterialBusinessDetail(whereList []*model.WhereSqlCond) (data *go_micro_srv_cust.AuthMaterialBusinessData, errR error) {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	whereModel := ss_sql.SsSqlFactoryInst.InitWhereList(whereList)
	ss_sql.SsSqlFactoryInst.AppendWhereExtra(whereModel, " ORDER BY amb.create_time DESC ")
	ss_sql.SsSqlFactoryInst.AppendWhereExtra(whereModel, " LIMIT 1 ")

	sqlStr := "select amb.auth_material_no, amb.license_img_no, amb.auth_name, amb.auth_number, amb.create_time," +
		" amb.auth_time, amb.account_uid, amb.status, amb.start_date, amb.end_date," +
		" amb.term_type, amb.simplify_name, amb.industry_no, bi.up_code " +
		" from auth_material_business amb " +
		" LEFT JOIN business_industry bi ON bi.code = amb.industry_no " + whereModel.WhereStr

	rows, stmt, err := ss_sql.QueryRowN(dbHandler, sqlStr, whereModel.Args...)
	if stmt != nil {
		defer stmt.Close()
	}
	if err != nil {
		ss_log.Error("err=[%v]", err)
		return nil, err
	}
	dataT := &go_micro_srv_cust.AuthMaterialBusinessData{}
	var authMaterialNo, licenseImgNo, authName, authNumber sql.NullString
	var createTime, authTime, accountUid, status sql.NullString
	var startDate, endDate, termType, simplifyName, industryNo, upIndustryNo sql.NullString
	err = rows.Scan(
		&authMaterialNo,
		&licenseImgNo,
		&authName,
		&authNumber,
		&createTime,

		&authTime,
		&accountUid,
		&status,
		&startDate,
		&endDate,

		&termType,
		&simplifyName,
		&industryNo,
		&upIndustryNo,
	)
	if err != nil && err != sql.ErrNoRows {
		ss_log.Error("err=[%v],AuthMaterialNo=[%v]", err, dataT.AuthMaterialNo)
		return nil, err
	}

	dataT.Status = status.String

	if dataT.Status == "" {
		dataT.Status = constants.AuthMaterialStatus_UnAuth
	} else {
		dataT.AuthMaterialNo = authMaterialNo.String
		dataT.LicenseImgNo = licenseImgNo.String
		dataT.AuthName = authName.String
		dataT.AuthNumber = authNumber.String
		dataT.CreateTime = createTime.String
		dataT.AuthTime = authTime.String
		dataT.AccountUid = accountUid.String
		dataT.StartDate = startDate.String
		dataT.EndDate = endDate.String
		dataT.TermType = termType.String
		dataT.SimplifyName = simplifyName.String
		dataT.IndustryNo = industryNo.String
		dataT.UpIndustryNo = upIndustryNo.String

		if dataT.Status == constants.AuthMaterialStatus_Passed { //只有通过的认证材料才能修改
			//修改简称的审核状态
			dataT.UpdateStatus, err = AuthMaterialDaoInst.GetAuthMaterialUpdateInfoStatus(dataT.AuthMaterialNo, dataT.AccountUid, constants.AccountType_PersonalBusiness)
			if err != nil && err != sql.ErrNoRows {
				ss_log.Error("查询修改认证材料信息状态失败，err=[%v]", err)
			}
		}

	}

	return dataT, nil
}

//查询企业商家认证信息（单个）
func (AuthMaterialDao) GetAuthMaterialEnterpriseDetail(whereList []*model.WhereSqlCond) (data *go_micro_srv_cust.AuthMaterialEnterpriseData, errR error) {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	whereModel := ss_sql.SsSqlFactoryInst.InitWhereList(whereList)
	ss_sql.SsSqlFactoryInst.AppendWhereExtra(whereModel, " ORDER BY create_time DESC ")
	ss_sql.SsSqlFactoryInst.AppendWhereExtra(whereModel, " LIMIT 1 ")

	//sqlStr := "select am.auth_material_no, am.license_img_no, am.auth_name, am.auth_number, am.create_time, am.auth_time, " +
	//	"am.account_uid, acc.business_auth_status, am.start_date, am.end_date, am.term_type, am.simplify_name " +
	//	" from account acc " +
	//	" left join auth_material_enterprise am on am.auth_material_no = acc.business_auth_material_no " + whereModel.WhereStr
	sqlStr := "select auth_material_no, license_img_no, auth_name, auth_number, create_time, auth_time, " +
		" account_uid, status, start_date, end_date, term_type, simplify_name " +
		" from auth_material_enterprise  " + whereModel.WhereStr

	rows, stmt, err := ss_sql.QueryRowN(dbHandler, sqlStr, whereModel.Args...)
	if stmt != nil {
		defer stmt.Close()
	}
	if err != nil {
		ss_log.Error("err=[%v]", err)
		return nil, err
	}
	dataT := &go_micro_srv_cust.AuthMaterialEnterpriseData{}
	var authMaterialNo, licenseImgNo, authName, authNumber sql.NullString
	var createTime, authTime, accountUid, status sql.NullString
	var startDate, endDate, termType, simplifyName sql.NullString
	err = rows.Scan(
		&authMaterialNo,
		&licenseImgNo,
		&authName,
		&authNumber,

		&createTime,
		&authTime,
		&accountUid,
		&status,
		&startDate,
		&endDate,
		&termType,
		&simplifyName,
	)
	if err != nil && err != sql.ErrNoRows { //查询不到说明是没有提交过认证材料
		ss_log.Error("err=[%v],AuthMaterialNo=[%v]", err, dataT.AuthMaterialNo)
		return nil, err
	}

	dataT.Status = status.String
	if dataT.Status == "" {
		dataT.Status = constants.AuthMaterialStatus_UnAuth
	} else {
		dataT.AuthMaterialNo = authMaterialNo.String
		dataT.LicenseImgNo = licenseImgNo.String
		dataT.AuthName = authName.String
		dataT.AuthNumber = authNumber.String
		dataT.CreateTime = createTime.String
		dataT.AuthTime = authTime.String
		dataT.AccountUid = accountUid.String
		dataT.StartDate = startDate.String
		dataT.EndDate = endDate.String
		dataT.TermType = termType.String
		dataT.SimplifyName = simplifyName.String

		if dataT.Status == constants.AuthMaterialStatus_Passed { //只有通过的认证材料才能修改
			//修改简称的审核状态
			dataT.UpdateStatus, err = AuthMaterialDaoInst.GetAuthMaterialUpdateInfoStatus(dataT.AuthMaterialNo, dataT.AccountUid, constants.AccountType_EnterpriseBusiness)
			if err != nil && err != sql.ErrNoRows {
				ss_log.Error("查询修改认证材料信息状态失败，err=[%v]", err)
			}
		}
	}

	return dataT, nil
}

//确认要通过的企业认证名称是唯一的
func (AuthMaterialDao) CheckAuthNameUnique(authName string) bool {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	//sqlStr := "select count(1) from auth_material_enterprise where status = $1 and auth_name = $2 "
	sqlStr := "select count(1) from business where full_name = $1 and is_delete = '0' "
	var cnt sql.NullString
	if err := ss_sql.QueryRow(dbHandler, sqlStr, []*sql.NullString{&cnt}, authName); err != nil {
		ss_log.Error("err=[%v]", err)
		return false
	}

	return cnt.String == "0"
}

//确认要通过的企业注册号/机构组织代码是唯一的
func (AuthMaterialDao) CheckAuthNumberUnique(authNumber string) bool {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	//sqlStr := "select count(1) from auth_material_enterprise where status = $1 and auth_number = $2 "
	sqlStr := "select count(1) from business where auth_number = $1 and is_delete = '0' "
	var cnt sql.NullString
	if err := ss_sql.QueryRow(dbHandler, sqlStr, []*sql.NullString{&cnt}, authNumber); err != nil {
		ss_log.Error("err=[%v]", err)
		return false
	}

	return cnt.String == "0"
}

//确认要通过的简称是唯一的
func (AuthMaterialDao) CheckSimplifyNameUnique(simplifyName string) bool {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	sqlStr := "select count(1) from business where simplify_name = $1 and is_delete = '0' "
	var cnt sql.NullString
	if err := ss_sql.QueryRow(dbHandler, sqlStr, []*sql.NullString{&cnt}, simplifyName); err != nil {
		ss_log.Error("err=[%v]", err)
		return false
	}

	return cnt.String == "0"
}

//获取修改商家认证信息的审核状态（如果没提交过则为空）
func (AuthMaterialDao) GetAuthMaterialUpdateInfoStatus(authMaterialNo, accountUid, accountType string) (status string, err error) {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	sqlStr := "select status " +
		" from auth_material_update_info " +
		" where auth_material_no = $1 and account_uid = $2 and account_type = $3 " +
		" order by create_time desc " +
		" limit 1 "
	var statusT sql.NullString
	if err := ss_sql.QueryRow(dbHandler, sqlStr, []*sql.NullString{&statusT}, authMaterialNo, accountUid, accountType); err != nil {
		ss_log.Error("err=[%v]", err)
		return "", err
	}

	return statusT.String, nil
}

//确认要修改认证信息的申请是唯一的（没有审核中的）
func (AuthMaterialDao) CheckAuthMaterialPendingStatusUnique(accountNo, accountType string) bool {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	sqlStr := "select count(1) from auth_material_update_info where status = $1 and account_uid = $2 and account_type = $3 "
	var cnt sql.NullString
	if err := ss_sql.QueryRow(dbHandler, sqlStr, []*sql.NullString{&cnt}, constants.AuthMaterialStatus_Pending, accountNo, accountType); err != nil {
		ss_log.Error("err=[%v]", err)
		return false
	}

	return cnt.String == "0"
}

//添加修改个人、企业商家的认证信息（后台要审核）
func (AuthMaterialDao) AddAuthMaterialUpdateInfo(authMaterialNo, accountUid, accountType, simplifyName, oldSimplifyName string) error {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	id := strext.GetDailyId()
	sqlStr := "insert into auth_material_update_info(id, auth_material_no, account_uid, account_type, simplify_name, old_simplify_name, status, create_time) " +
		" values($1,$2,$3,$4,$5,$6,$7,current_timestamp)"
	if err := ss_sql.Exec(dbHandler, sqlStr, id, authMaterialNo, accountUid, accountType, simplifyName, oldSimplifyName, constants.AuthMaterialStatus_Pending); err != nil {
		ss_log.Error("err=[%v]", err)
		return err
	}

	return nil
}

func (AuthMaterialDao) GetAuthMaterialBusinessUpdateCnt(whereList []*model.WhereSqlCond) (total string, errR error) {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	whereModel := ss_sql.SsSqlFactoryInst.InitWhereList(whereList)

	sqlStr := "select count(1) " +
		" from auth_material_update_info am " +
		" left join account acc on am.account_uid = acc.uid "
	var cnt sql.NullString
	if err := ss_sql.QueryRow(dbHandler, sqlStr, []*sql.NullString{&cnt}, whereModel.Args...); err != nil {
		ss_log.Error("err=[%v]", err)
		return "0", err
	}

	return cnt.String, nil
}

func (AuthMaterialDao) GetAuthMaterialBusinessUpdateList(whereList []*model.WhereSqlCond, page, pageSize int32) (datas []*go_micro_srv_cust.AuthMaterialBusinessUpdateData, errR error) {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	whereModel := ss_sql.SsSqlFactoryInst.InitWhereList(whereList)

	ss_sql.SsSqlFactoryInst.AppendWhereExtra(whereModel, " order by am.create_time desc, case am.status when "+constants.AuthMaterialStatus_Pending+" then 1 end")
	ss_sql.SsSqlFactoryInst.AppendWhereLimitI32(whereModel, pageSize, page)

	sqlStr := "select  am.id, am.auth_material_no, am.create_time, am.account_uid, am.status," +
		" am.check_time, am.simplify_name, am.account_type, am.notes, am.old_simplify_name, " +
		" acc.account " +
		" from auth_material_update_info am " +
		" left join account acc on am.account_uid = acc.uid "

	rows, stmt, err := ss_sql.Query(dbHandler, sqlStr+whereModel.WhereStr, whereModel.Args...)
	if stmt != nil {
		defer stmt.Close()
	}
	if rows != nil {
		defer rows.Close()
	}
	if err != nil {
		ss_log.Error("err=[%v]", err)
		return nil, err
	}

	var checkTime sql.NullString

	for rows.Next() {
		data := &go_micro_srv_cust.AuthMaterialBusinessUpdateData{}

		err = rows.Scan(
			&data.Id,
			&data.AuthMaterialNo,
			&data.CreateTime,
			&data.AccountUid,
			&data.Status,

			&checkTime,
			&data.SimplifyName,
			&data.AccountType,
			&data.Notes,
			&data.OldSimplifyName,
			&data.Account,
		)
		if err != nil {
			ss_log.Error("err=[%v],Id=[%v]", err, data.Id)
			return nil, err
		}

		data.CheckTime = checkTime.String
		datas = append(datas, data)
	}

	return datas, nil
}

func (AuthMaterialDao) GetAuthMaterialBusinessUpdateDetail(whereList []*model.WhereSqlCond) (data *go_micro_srv_cust.AuthMaterialBusinessUpdateData, errR error) {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	whereModel := ss_sql.SsSqlFactoryInst.InitWhereList(whereList)

	sqlStr := "select  am.id, am.auth_material_no, am.create_time, am.account_uid, am.status," +
		" am.check_time, am.simplify_name, am.account_type, am.notes, am.old_simplify_name, " +
		" acc.account " +
		" from auth_material_update_info am " +
		" left join account acc on am.account_uid = acc.uid "
	rows, stmt, err := ss_sql.QueryRowN(dbHandler, sqlStr+whereModel.WhereStr, whereModel.Args...)
	if stmt != nil {
		defer stmt.Close()
	}
	if err != nil {
		ss_log.Error("err=[%v]", err)
		return nil, err
	}

	var checkTime sql.NullString
	dataT := &go_micro_srv_cust.AuthMaterialBusinessUpdateData{}
	err = rows.Scan(
		&dataT.Id,
		&dataT.AuthMaterialNo,
		&dataT.CreateTime,
		&dataT.AccountUid,
		&dataT.Status,

		&checkTime,
		&dataT.SimplifyName,
		&dataT.AccountType,
		&dataT.Notes,
		&dataT.OldSimplifyName,
		&dataT.Account,
	)
	if err != nil {
		ss_log.Error("err=[%v],Id=[%v]", err, dataT.Id)
		return nil, err
	}

	dataT.CheckTime = checkTime.String

	return dataT, nil
}

func (AuthMaterialDao) ModifyAuthMaterialBusinessUpdateStatus(tx *sql.Tx, id, status, notes string) error {
	sqlStr := " update auth_material_update_info set status = $3, notes = $4, check_time = current_timestamp where id = $1 and status = $2  "
	if err := ss_sql.ExecTx(tx, sqlStr, id, constants.AuthMaterialStatus_Pending, status, notes); err != nil {
		ss_log.Error("err=[%v]", err)
		return err
	}
	return nil
}

//修改企业商家认证材料的简称
func (AuthMaterialDao) ModifyAuthMaterialEnterpriseSimplifyName(tx *sql.Tx, authMaterialNo, simplifyName, oldStatus string) (errR error) {
	sqlStr := "update auth_material_enterprise set simplify_name = $2 where auth_material_no = $1 and status = $3 "
	err := ss_sql.ExecTx(tx, sqlStr, authMaterialNo, simplifyName, oldStatus)
	if err != nil {
		ss_log.Error("err=[%v]", err)
		return err
	}
	return nil
}

//审核商家认证材料后的修改账号商家认证状态
func (AuthMaterialDao) ModifyAccountBusinessAuthStatus(tx *sql.Tx, uid, status, oldStatus string) (errR error) {
	sqlStr := "update account set business_auth_status = $2 where uid = $1  and business_auth_status = $3 and is_delete = '0' "
	err := ss_sql.ExecTx(tx, sqlStr, uid, status, oldStatus)
	if err != nil {
		ss_log.Error("err=[%v]", err)
		return err
	}
	return nil
}

//上传商家认证材料后的修改账号商家认证状态
func (AuthMaterialDao) ModifyBusinessAccountAuthStatus(tx *sql.Tx, uid, status string) (errR error) {
	sqlStr := "update account set business_auth_status = $2 where uid = $1 and is_delete = '0' "
	err := ss_sql.ExecTx(tx, sqlStr, uid, status)
	if err != nil {
		ss_log.Error("err=[%v]", err)
		return err
	}
	return nil
}
