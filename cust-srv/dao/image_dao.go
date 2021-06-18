package dao

import (
	"a.a/cu/db"
	"a.a/cu/ss_log"
	"a.a/cu/strext"
	"a.a/mp-server/common/cache"
	"a.a/mp-server/common/constants"
	"a.a/mp-server/common/ss_sql"
	"database/sql"
)

var (
	ImageDaoInstance ImageDao
)

type ImageDao struct {
	ImageId     string
	ImageUrl    string
	ContentType string
	IsEncrypt   int32
	Ext         string
}

func (r *ImageDao) InsertImage(accountNo, imageName, contentType, ext string, isEncrypt int) (string, error) {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	imageIDT := strext.GetDailyId()

	sql := "insert into dict_images(image_id, image_url, status, account_no, content_type, is_encrypt, ext, create_time, modify_time)"
	sql += "values($1,$2,$3,$4, $5, $6, $7, current_timestamp, current_timestamp)"

	err := ss_sql.Exec(dbHandler, sql,
		imageIDT, imageName, "1", accountNo, contentType, isEncrypt, ext)
	if err != nil {
		ss_log.Error("err=[%v]", err)
		return "", err
	}
	return imageIDT, nil
}

func (r *ImageDao) GetImageUrlById(id string) (ImageDao, error) {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	imgDao := ImageDao{}
	var imageId, imageUrl, contentType, isEncrypt, ext sql.NullString

	sqlStr := "select image_id, image_url, content_type, is_encrypt, ext from dict_images where image_id =$1 and is_delete = '0' "
	qErr := ss_sql.QueryRow(dbHandler, sqlStr, []*sql.NullString{&imageId, &imageUrl, &contentType, &isEncrypt, &ext}, id)
	if qErr != nil {
		return imgDao, qErr
	}

	imgDao.ImageId = imageId.String
	imgDao.ImageUrl = imageUrl.String
	imgDao.ContentType = contentType.String
	imgDao.IsEncrypt = strext.ToInt32(isEncrypt.String)
	imgDao.Ext = ext.String

	return imgDao, qErr
}

//查询实名认证的两张图片id是否在数据库中
func (r *ImageDao) CheckImageById1Id2(id1, id2 string) (total string, err error) {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	var totalT sql.NullString
	sqlStr := "select count(1) from dict_images where image_id in ($1,$2) and is_delete = '0' "
	qErr := ss_sql.QueryRow(dbHandler, sqlStr, []*sql.NullString{&totalT}, id1, id2)
	if qErr != nil {
		ss_log.Error("err=[%v]", qErr)
		return "", qErr
	}

	return totalT.String, qErr
}

func (r *ImageDao) GetImgUrlsByImgIds(imgIds []string) (imgUrls []string) {
	// 获取不需要授权的图片路径
	_, imageBaseUrl, _ := cache.ApiDaoInstance.GetGlobalParam("image_base_url")

	var imgUrlsT []string
	if len(imgIds) > 0 {
		for _, imgId := range imgIds {
			if imgDao, err2 := ImageDaoInstance.GetImageUrlById(imgId); err2 != nil {
				ss_log.Error("查询图片记录失败,ImageId:%s, err:%v", imgId, err2)
				imgUrlsT = append(imgUrlsT, "")
			} else {
				imgUrlsT = append(imgUrlsT, imageBaseUrl+"/"+imgDao.ImageUrl)
			}
		}
	}
	return imgUrlsT
}
