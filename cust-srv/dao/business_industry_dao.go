package dao

import (
	"a.a/cu/db"
	"a.a/cu/ss_log"
	"a.a/cu/strext"
	"a.a/mp-server/common/constants"
	"a.a/mp-server/common/model"
	go_micro_srv_cust "a.a/mp-server/common/proto/cust"
	"a.a/mp-server/common/ss_count"
	"a.a/mp-server/common/ss_sql"
	"database/sql"
	"fmt"
)

type BusinessIndustryDao struct {
	Code   string
	NameCh string
	NameEn string
	NameKm string
	Level  string

	UpCode string
}

var BusinessIndustryDaoInst BusinessIndustryDao

func (BusinessIndustryDao) GetBusinessIndustryCnt(whereList []*model.WhereSqlCond) (total string) {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	whereModel := ss_sql.SsSqlFactoryInst.InitWhereList(whereList)

	var totalT sql.NullString
	sqlStr := "SELECT count(1) " +
		"FROM business_industry " + whereModel.WhereStr
	if err2 := ss_sql.QueryRow(dbHandler, sqlStr, []*sql.NullString{&totalT}, whereModel.Args...); err2 != nil {
		ss_log.Error("err=[%v]", err2)
		return ""
	}

	return totalT.String
}

func (BusinessIndustryDao) GetBusinessIndustryDatas(whereList []*model.WhereSqlCond, page, pageSize string) (datasT []*go_micro_srv_cust.MainIndustryData, err error) {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	whereModel := ss_sql.SsSqlFactoryInst.InitWhereList(whereList)

	ss_sql.SsSqlFactoryInst.AppendWhereExtra(whereModel, " ORDER BY code asc, level asc ")
	if page != "" && pageSize != "" {
		ss_sql.SsSqlFactoryInst.AppendWhereLimit(whereModel, strext.ToInt(pageSize), strext.ToInt(page))
	}

	sqlStr := "SELECT code, name_ch, name_en, name_km, level," +
		" up_code, create_time, modify_time " +
		"FROM business_industry " + whereModel.WhereStr
	rows, stmt, err2 := ss_sql.Query(dbHandler, sqlStr, whereModel.Args...)
	if stmt != nil {
		defer stmt.Close()
	}
	defer rows.Close()
	if err2 != nil {
		ss_log.Error("err2=[%v]", err2)
		return nil, err2
	}

	var datas []*go_micro_srv_cust.MainIndustryData
	for rows.Next() {
		data := go_micro_srv_cust.MainIndustryData{}
		var modifyTime sql.NullString
		err2 = rows.Scan(
			&data.Code,
			&data.NameCh,
			&data.NameEn,
			&data.NameKm,
			&data.Level,

			&data.UpCode,
			&data.CreateTime,

			&modifyTime,
		)

		if err2 != nil {
			ss_log.Error("err=[%v]", err2)
			continue
		}
		data.ModifyTime = modifyTime.String

		datas = append(datas, &data)
	}

	return datas, nil
}

func (BusinessIndustryDao) GetBusinessIndustryDetail(code string) (*go_micro_srv_cust.MainIndustryData, error) {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	whereModel := ss_sql.SsSqlFactoryInst.InitWhereList([]*model.WhereSqlCond{
		{Key: "is_delete", Val: "0", EqType: "="},
		{Key: "code", Val: code, EqType: "="},
	})

	sqlStr := "SELECT code, name_ch, name_en, name_km, level," +
		" up_code, create_time, modify_time " +
		"FROM business_industry " + whereModel.WhereStr
	rows, stmt, err2 := ss_sql.QueryRowN(dbHandler, sqlStr, whereModel.Args...)
	if stmt != nil {
		defer stmt.Close()
	}
	if err2 != nil {
		ss_log.Error("err2=[%v]", err2)
		return nil, err2
	}

	dataT := &go_micro_srv_cust.MainIndustryData{}
	var codeT, nameCh, nameEn, nameKm, level,
		upCode, createTime, modifyTime sql.NullString
	err2 = rows.Scan(
		&codeT,
		&nameCh,
		&nameEn,
		&nameKm,
		&level,

		&upCode,
		&createTime,
		&modifyTime,
	)

	if err2 != nil {
		ss_log.Error("err=[%v]", err2)
		return nil, err2
	}
	dataT.Code = codeT.String
	dataT.NameCh = nameCh.String
	dataT.NameEn = nameEn.String
	dataT.NameKm = nameKm.String
	dataT.Level = level.String

	dataT.UpCode = upCode.String
	dataT.CreateTime = createTime.String
	dataT.ModifyTime = modifyTime.String

	return dataT, nil
}

func (BusinessIndustryDao) getMaxCodeByLevelAndUpCode(level, upCode string) (maxCode string) {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	//sqlStr := "select code from business_industry where level = $1 and up_code = $2 and is_delete = '0' order by code desc limit 1 "
	sqlStr := "select code from business_industry where level = $1 and up_code = $2 order by code desc limit 1 "
	var code sql.NullString
	if err := ss_sql.QueryRow(dbHandler, sqlStr, []*sql.NullString{&code}, level, upCode); err != nil {
		switch level {
		case "1":
			return "000"
		case "2":
			return upCode + "000"
		}
	}
	return code.String
}

//level 不允许为空
func (b BusinessIndustryDao) AddBusinessIndustry(data BusinessIndustryDao) (err error) {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	code := ss_count.Add(b.getMaxCodeByLevelAndUpCode(data.Level, data.UpCode), "1")

	switch data.Level {
	case "1":
		code = fmt.Sprintf("%03d", strext.ToInt(code))
		code = fmt.Sprintf("%s", code)
	case "2":
		code = fmt.Sprintf("%06d", strext.ToInt(code))
		code = fmt.Sprintf("%s", code)
	}

	sqlStr := "insert into business_industry(code, name_ch, name_en, name_km, level, up_code, create_time, modify_time) " +
		" values($1,$2,$3,$4,$5,$6,current_timestamp,current_timestamp) "
	if err2 := ss_sql.Exec(dbHandler, sqlStr, code, data.NameCh, data.NameEn, data.NameKm, data.Level, data.UpCode); err2 != nil {
		ss_log.Error("err=[%v]", err2)
		return err2
	}

	return nil
}

func (b BusinessIndustryDao) UpdateBusinessIndustry(data BusinessIndustryDao) (err error) {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	sqlStr := "update business_industry set name_ch  = $2, name_en = $3, name_km = $4, level = $5, up_code = $6, modify_time = current_timestamp " +
		" where code = $1 "
	if err2 := ss_sql.Exec(dbHandler, sqlStr, data.Code, data.NameCh, data.NameEn, data.NameKm, data.Level, data.UpCode); err2 != nil {
		ss_log.Error("err=[%v]", err2)
		return err2
	}

	return nil
}

func (b BusinessIndustryDao) DelBusinessIndustry(code string) (err error) {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	sqlStr := "update business_industry set is_delete = '1' where code = $1 and is_delete = '0' "
	if err2 := ss_sql.Exec(dbHandler, sqlStr, code); err2 != nil {
		ss_log.Error("err=[%v]", err2)
		return err2
	}

	return nil
}

func (b BusinessIndustryDao) CountByUpCode(code string) (total string) {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	sqlStr := "select count(1) from business_industry where up_code = $1 and is_delete = '0' "
	var cnt sql.NullString
	if err2 := ss_sql.QueryRow(dbHandler, sqlStr, []*sql.NullString{&cnt}, code); err2 != nil {
		ss_log.Error("err=[%v]", err2)
		return "0"
	}

	return cnt.String
}

//处理数据，使数据符合ElementUI组件Cascader级联选择器的数据格式(本方法现只生成两级)
func (*BusinessIndustryDao) TreatMainIndustryData(datas []*go_micro_srv_cust.MainIndustryData, lang string) (datas2 []*go_micro_srv_cust.MainIndustryCascaderData) {
	m := make(map[string]*go_micro_srv_cust.MainIndustryData)
	var datasT []*go_micro_srv_cust.MainIndustryCascaderData

	for _, data := range datas {
		if _, ok := m[data.Code]; !ok { //没添加过该节点
			name := "" //行业名称
			switch lang {
			case constants.LangZhCN:
				name = data.NameCh
			case constants.LangEnUS:
				name = data.NameEn
			case constants.LangKmKH:
				name = data.NameKm
			default:
				return nil
			}

			// map 中不存在,需要创建新的切片,添加到map
			m[data.Code] = data

			if data.UpCode != "" && data.Level == "2" { //是二级子节点
				if _, ok := m[data.UpCode]; ok {
					for kk, vv := range datasT {
						if vv.Value == data.UpCode { //查询父节点，在父节点下添加它
							vv.Children = append(vv.Children, &go_micro_srv_cust.MainIndustryCascaderData{
								Value:    data.Code,
								Label:    name,
								Children: nil,
							})
							datasT[kk] = vv
						}
					}
				} else {
					ss_log.Error("找不到父节点,code=[%v]", data.Code)
				}

				continue
			}

			//不是子节点，直接添加
			datasT = append(datasT, &go_micro_srv_cust.MainIndustryCascaderData{
				Value:    data.Code,
				Label:    name,
				Children: nil,
			})
		}
	}

	return datasT
}
