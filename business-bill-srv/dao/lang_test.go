package dao

import "testing"

func TestLangDao_GetLangTextByKey(t *testing.T) {
	key := "转账成功推送消息模板"
	lang := ""
	langText := LangDaoInstance.GetLangTextByKey(key, lang)
	t.Logf("结果：%v", langText)
}
