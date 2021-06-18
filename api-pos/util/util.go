package util

import (
	"net/http"
	"strconv"

	"a.a/cu/strext"
	"github.com/gin-gonic/gin"
)

func TransitionToFloat64(json interface{}, key string) float64 {
	var str interface{}
	str = json.(map[string]interface{})[key]
	switch data := str.(type) {
	case int:
		return float64(data)
	case float64:
		return data
	case string:
		d, err := strconv.ParseFloat(data, 64)
		if err == nil {
			return d
		} else {
			return 0
		}
	default:
		return 0
	}
}

func TransitionToInt32(json interface{}, key string) int32 {
	var str interface{}
	str = json.(map[string]interface{})[key]
	switch data := str.(type) {
	case int:
		return int32(data)
	case float64:
		return int32(data)
	case string:
		d, err := strconv.Atoi(data)
		if err == nil {
			return int32(d)
		} else {
			return 0
		}
	default:
		return 0
	}
}
func TransitionToString(json interface{}, key string) string {
	str := json.(map[string]interface{})[key]
	switch data := str.(type) {
	case int:
		return strconv.Itoa(data)
	case float64:
		return strconv.FormatFloat(data, 'E', -1, 64)
	case string:
		return data
	default:
		return ""
	}
}

/**
 * 参数为空，返回信息
 */
func ResultJsonEmptyData(c *gin.Context, json interface{}, key string) bool {
	data := json.(map[string]interface{})[key]
	if data == nil || data == "" {
		c.Set("status", 400)
		c.Set("resp", gin.H{
			"status":  400,
			"retcode": 0,
			"msg":     "参数有误",
			"errmsg":  key + ":为空",
		})
		c.Abort()
		return true
	}
	return false
}

/**
 * 参数为空，返回信息
 */
func ResultEmptyData(c *gin.Context, data interface{}, dataName string) bool {
	if data == nil || data == "" {
		c.Set("status", http.StatusForbidden)
		c.Set("resp", gin.H{
			"status":  400,
			"retcode": 0,
			"msg":     "参数有误",
			"errmsg":  dataName + ":为空",
		})
		c.Abort()
		return true
	}
	return false
}

func ResultList(c *gin.Context, data interface{}, total int32) {
	c.Set("status", 200)
	c.Set("resp", gin.H{
		"status":  200,
		"retcode": 0,
		"msg":     "成功",
		"data": gin.H{
			"data":  data,
			"total": total,
		},
	})
	c.Next()
}

func ResultMap(c *gin.Context, data interface{}) {
	c.Set("status", 200)
	c.Set("resp", gin.H{
		"status":  200,
		"retcode": 0,
		"msg":     "成功",
		"data": gin.H{
			"data": data,
		},
	})
	c.Next()
}

func ResultUid(c *gin.Context, uid string) {
	c.Set("status", 200)
	c.Set("resp", gin.H{
		"status":  200,
		"retcode": 0,
		"data": gin.H{
			"uid": uid,
		},
	})
	c.Next()
}

func ResultInfo(c *gin.Context) {
	c.Set("status", 200)
	c.Set("resp", gin.H{
		"status":  200,
		"retcode": 0,
		"msg":     "成功",
	})
	c.Next()
}

func ResultPath(c *gin.Context, path string) {
	c.Set("status", 200)
	c.Set("resp", gin.H{
		"status":  200,
		"retcode": 0,
		"msg":     "成功",
		"data": gin.H{
			"path": path,
		},
	})
	c.Next()
}

func IsContainMenu(b interface{}, l string) bool {
	switch b.(type) {
	case []map[string]interface{}:
		for _, v := range b.([]map[string]interface{}) {
			for k2, v2 := range v {
				if strext.ToStringNoPoint(k2) == "children" {
					for _, v3 := range v2.([]interface{}) {
						if v3.(map[string]interface{})["path"] == l {
							return true
						} else if v3.(map[string]interface{})["children"] != nil {
							if IsContainMenu(v3.(map[string]interface{})["children"].([]interface{}), l) {
								return true
							}
						}
					}
				} else {
					if strext.ToStringNoPoint(k2) == "path" && strext.ToStringNoPoint(v2) == l {
						return true
					}
				}
			}
		}
	case []interface{}:
		for _, v := range b.([]interface{}) {
			for k2, v2 := range v.(map[string]interface{}) {
				if strext.ToStringNoPoint(k2) == "children" {
					for _, v3 := range v2.([]interface{}) {
						if v3.(map[string]interface{})["path"] == l {
							return true
						} else if v3.(map[string]interface{})["children"] != nil {
							if IsContainMenu(v3.(map[string]interface{})["children"].([]interface{}), l) {
								return true
							}
						}
					}
				} else {
					if strext.ToStringNoPoint(k2) == "path" && strext.ToStringNoPoint(v2) == l {
						return true
					}
				}
			}
		}
	}

	return false
}
