package dao

import (
	"a.a/cu/db"
	"a.a/cu/ss_log"
	"a.a/mp-server/common/constants"
	"a.a/mp-server/common/ss_sql"
	"database/sql"
)

type DictimagesDao struct{}

var DictimagesDaoInst DictimagesDao

func (DictimagesDao) GetImgIds(servicerNo string) (ids string, err error) {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	sqlStr := `select 
				si.img_ids 
			from  servicer_img si
			WHERE si.servicer_no = $1 AND si.is_delete = 0`
	var imgIds sql.NullString
	errT := ss_sql.QueryRow(dbHandler, sqlStr, []*sql.NullString{&imgIds}, servicerNo)
	if errT != nil {
		ss_log.Error("err=[%v]", errT)
		return "", errT
	}

	return imgIds.String, nil
}
func (DictimagesDao) GetKmImgPath(langKm string) (string, error) {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	sqlStr := `select di.image_url from lang l LEFT JOIN dict_images di ON l.lang_km = di.image_id 
	WHERE l.lang_km = $1`
	var image sql.NullString
	errT := ss_sql.QueryRow(dbHandler, sqlStr, []*sql.NullString{&image}, langKm)
	if errT != nil {
		ss_log.Error("err=[%v]", errT)
		return "", errT
	}

	return image.String, nil
}
func (DictimagesDao) GetEnImgPath(langEn string) (string, error) {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	sqlStr := `select di.image_url from lang l LEFT JOIN dict_images di ON l.lang_en = di.image_id 
	WHERE l.lang_en = $1`
	var image sql.NullString
	errT := ss_sql.QueryRow(dbHandler, sqlStr, []*sql.NullString{&image}, langEn)
	if errT != nil {
		ss_log.Error("err=[%v]", errT)
		return "", errT
	}

	return image.String, nil
}
func (DictimagesDao) GetChImgPath(langCh string) (string, error) {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	sqlStr := `select di.image_url from lang l LEFT JOIN dict_images di ON l.lang_ch = di.image_id 
	WHERE l.lang_ch = $1`
	var image sql.NullString
	errT := ss_sql.QueryRow(dbHandler, sqlStr, []*sql.NullString{&image}, langCh)
	if errT != nil {
		ss_log.Error("err=[%v]", errT)
		return "", errT
	}

	return image.String, nil
}

func (DictimagesDao) DeleteTx(tx *sql.Tx, imgPath string) error {
	return ss_sql.ExecTx(tx, `update dict_images set is_delete=1,modify_time = current_timestamp where image_url=$1`,
		imgPath)

}
func (DictimagesDao) Delete(imgPath string) error {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	return ss_sql.Exec(dbHandler, `update dict_images set is_delete=1,modify_time = current_timestamp where image_url=$1`,
		imgPath)
}

//删除图片失败的错误信息
func (DictimagesDao) AddDelFaildLog(notes string) error {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	return ss_sql.Exec(dbHandler, `insert into log_del_img_err(create_time,notes) values(current_timestamp,$1)`,
		notes)
}
