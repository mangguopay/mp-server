package read_file

import (
	"fmt"
	"github.com/360EntSecGroup-Skylar/excelize"
	"os"
	"strings"
)

func ReadExcelFmtPhp() (fileName string) {
	url := "D:\\Documents\\WeChat Files\\wxid_bt9xfpofhut722\\FileStorage\\File\\2020-11\\后台需要翻译的高棉语（筛选后）校对版1.0(1)(1).xlsx"
	file, err := excelize.OpenFile(url)
	if err != nil {
		fmt.Println(err)
		return ""
	}

	var s = `<?php
return [
`

	rows := file.GetRows(file.GetSheetName(1))
	for _, row := range rows {
		if row != nil {
			if len(row) > 2 && (strings.TrimSpace(row[0]) != "" && strings.TrimSpace(row[2]) != "") {
				s += fmt.Sprintf(" '%s'	=>  '%s', \n", row[0], row[2])
			}

			//for k, colCell := range row {
			//
			//
			//	if colCell != "" && colCell != " " {
			//		if k == 0 {
			//			if colCell == "" {
			//
			//			}
			//			s += fmt.Sprintf("'%s'	=>  ", colCell)
			//		} else if k == 1 {
			//
			//		} else if k == 2 {
			//			if colCell == "" {
			//				continue
			//			}
			//			s += fmt.Sprintf("'%s',\n ", colCell)
			//		}
			//	}
			//}
		}
	}
	s += `];`
	name := "./软件备用字幕翻译2020-11-20.php"
	fp, err := os.Create(name)
	if err != nil {
		fmt.Println(err)
		return ""
	}

	defer fp.Close()
	_, _ = fp.Write([]byte(s))

	return name
}
