package middleware

import (
	"a.a/cu/container"
	"a.a/cu/ss_img"
	"a.a/cu/ss_log"
	"a.a/cu/strext"
	"a.a/mp-server/api-webbusiness/common"
	"a.a/mp-server/api-webbusiness/util"
	"a.a/mp-server/common/ss_err"
	"github.com/gin-gonic/gin"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"path"
)

type GetParamsMw struct {
}

var GetParamsMwInst GetParamsMw

/**
 * 获取post过来的json
 */
func (r GetParamsMw) FetchPostJsonBodyParams(notJson []string) gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.GetBool(common.INNER_IS_STOP) {
			ss_log.Info("skip")
			return
		}
		//===============================================
		//	只处理post
		if http.MethodPost != c.Request.Method && http.MethodDelete != c.Request.Method && http.MethodPatch != c.Request.Method {
			ss_log.Info("not post")
			return
		}
		ss_log.Info("post")

		traceNo := c.GetString(common.INNER_TRACE_NO)
		if container.GetKey(c.Request.RequestURI, notJson) < 0 { // 不在里面
			buf, err := ioutil.ReadAll(c.Request.Body)
			defer c.Request.Body.Close()
			if err != nil {
				ss_log.Error("请求包体为空|err=[%v]", err)
				c.Set(common.INNER_IS_STOP, true)
				c.Set(common.INNER_ERR_REASON, ss_err.ERR_SYS_EMPTY_BODY)
				c.Abort()
				return
			}

			p := strext.Json2Map(buf)
			if p == nil {
				ss_log.Error("%v|请求包体不是json|buf=[%v]", traceNo, string(buf))
				c.Set(common.INNER_IS_STOP, true)
				c.Set(common.INNER_ERR_REASON, ss_err.ERR_SYS_BODY_NOT_JSON)
				c.Abort()
				return
			}

			ss_log.Info("%v|----------------------------POST的参数", traceNo)
			for k, v := range p {
				ss_log.Info("%v|[%v]=>[%v]", traceNo, k, v)
			}
			ss_log.Info("%v|----------------------------", traceNo)
			c.Set(common.INNER_PARAM_MAP, p)
		} else {
			// 这里有可能是上传，不需要转json
			//
			uploadPath, filename, fileId, filenameWithSuffix := r.upload(c)
			defer c.Request.Body.Close()
			c.Set(common.INNER_PARAM_MAP, map[string]interface{}{
				"upload_path":          uploadPath,
				"filename":             filename,
				"file_id":              fileId,
				"filename_with_suffix": filenameWithSuffix,
			})
		}

		c.Set(common.INNER_FMT, "json")
		ss_log.Info("%v|recv[%v]=>[%v]", traceNo, c.Request.Method, strext.ToStringNoPoint(c.Request.RequestURI))
		return
	}
}

/**
 * 读取get的参数
 */
func (GetParamsMw) FetchGetParams() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.GetBool(common.INNER_IS_STOP) {
			ss_log.Info("skip")
			return
		}
		//===============================================
		//	只处理get
		if http.MethodGet != c.Request.Method {
			return
		}
		// 读取跟踪号
		traceNo := c.GetString(common.INNER_TRACE_NO)
		// 读取get的参数
		queryForm, _ := url.ParseQuery(c.Request.URL.RawQuery)
		p := map[string]interface{}{}
		ss_log.Info("%v|----------------------------GET的参数", traceNo)
		for k, v := range queryForm {
			ss_log.Info("%v|[%v]=>[%v]", traceNo, k, v[0])
			p[k] = v[0]
		}
		ss_log.Info("%v|----------------------------", traceNo)
		c.Set(common.INNER_PARAM_MAP, p)
		c.Set(common.INNER_IS_ENCODED, false)
		ss_log.Info("recv[%v]=>[%v]", c.Request.Method, strext.ToStringNoPoint(c.Request.RequestURI))
		return
	}
}

func (GetParamsMw) upload(c *gin.Context) (uploadPathT, filenameT, fileIdT, filenameWithSuffixT string) {
	//得到上传的文件
	file, header, err := c.Request.FormFile("file")
	if err != nil {
		ss_log.Error("获取文件信息出错，err=[%v]", err)
		util.ChknRet(c, ss_err.ERR_UPLOAD)
		return
	}
	//文件的名称
	filename := header.Filename
	filenameWithSuffix := path.Ext(header.Filename)
	ss_log.Info("%v,%v", err, filename)
	//创建文件
	pathStr := os.TempDir()
	ss_log.Info("临时目录:[%v]", pathStr)

	fileId := strext.GetDailyId()
	fileName2 := fileId + filenameWithSuffix
	uploadPath := path.Join(pathStr, fileName2)
	ss_log.Info("uploadPath=[%v]\n", uploadPath)

	out, err := os.Create(uploadPath)
	if err != nil {
		ss_log.Error("创建文件出错,err=[%v]\n", err)
		return "", "", "", ""
	}
	defer out.Close()
	_, err = io.Copy(out, file)
	if err != nil {
		ss_log.Error("写入文件出错,err=[%v]\n", err)
		return "", "", "", ""
	} else {
		fileB, err := ioutil.ReadFile(uploadPath)
		ss_log.Error("err=[%v]", err)
		ft := ss_img.SsImgInst.GetFileTypeFromMagic(fileB)
		ss_log.Error("ft=[%v]", ft)
		// 判定filetype
		switch ft {
		case "zip":
		// app
		case "":
		default:
			// !!! 有问题的文件，改后缀
			ss_log.Error("有问题的文件，改后缀")
			//关闭前面打开的文件,不然改不了名
			out.Close()
			errRename := os.Rename(uploadPath, uploadPath+".bad")
			if errRename != nil {
				ss_log.Error("errRename=[%v]", errRename)
			}
			return "", "", "", ""
		}
		return uploadPath, fileName2, fileId, filenameWithSuffix
	}
}
