package dao

import (
	"a.a/cu/db"
	"a.a/mp-server/common/constants"
	"a.a/mp-server/common/ss_sql"
	"database/sql"
)

type BusinessScene struct {
}

var BusinessSceneDao BusinessScene

func (*BusinessScene) GetSceneIsEnabled(sceneNo string) (bool, error) {
	dbHandler := db.GetDB(constants.DB_CRM)
	defer db.PutDB(constants.DB_CRM, dbHandler)

	sqlStr := "SELECT is_delete FROM business_scene WHERE scene_no = $1 "
	var isEnabled sql.NullString
	err := ss_sql.QueryRow(dbHandler, sqlStr, []*sql.NullString{&isEnabled}, sceneNo)
	if err != nil {
		return false, err
	}

	if isEnabled.String == "0" {
		return true, nil
	} else if isEnabled.String == "1" {
		return false, nil
	} else {
		return false, nil
	}
}
