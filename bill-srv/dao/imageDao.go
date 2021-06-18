package dao

import (
	"database/sql"

	"a.a/cu/ss_log"

	"a.a/cu/db"
	"a.a/cu/strext"
	"a.a/mp-server/common/constants"
	"a.a/mp-server/common/ss_err"
	"a.a/mp-server/common/ss_sql"
)

type ImageDao struct {
}

var (
	ImageDaoInstance ImageDao
)

func (r *ImageDao) InsertImage(accontNo, imageURL string) error {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	imageIDT := strext.NewUUID()
	err := ss_sql.Exec(dbHandler, `insert into dict_images(image_id,image_url,status,account_no,create_time,modify_time)values($1,$2,$3,$4,current_timestamp,current_timestamp)`,
		imageIDT, imageURL, "1", accontNo)
	if err != nil {
		ss_log.Error("err=%v", err)
		return err
	}
	return nil
}

func (ImageDao) GetImageUrlById(id string) (string, string) {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	var imageUrl sql.NullString
	sqlStr := "select image_url from dict_images where image_id =$1 and is_delete = '0' "
	err := ss_sql.QueryRow(dbHandler, sqlStr, []*sql.NullString{&imageUrl}, id)
	if err != nil {
		return "", ss_err.ERR_SYS_DB_GET
	}
	return imageUrl.String, ss_err.ERR_SUCCESS
}
