package processor

import (
	go_micro_srv_auth "a.a/mp-server/auth-srv/proto/auth"
	"a.a/mp-server/common/ss_sql"
	"github.com/gin-gonic/gin"
	"sort"
)

var (
	MenuProcInst MenuProc
)

type MenuProc struct {
}

type menuItem []*gin.H

func (c menuItem) Len() int {
	return len(c)
}

func (c menuItem) Swap(i, j int) {
	c[i], c[j] = c[j], c[i]
}

func (c menuItem) Less(i, j int) bool {
	return int((*c[i])["idx"].(int32)) < int((*c[j])["idx"].(int32))
}

type menuItem2 []gin.H

func (c menuItem2) Len() int {
	return len(c)
}

func (c menuItem2) Swap(i, j int) {
	c[i], c[j] = c[j], c[i]
}

func (c menuItem2) Less(i, j int) bool {
	return int(c[i]["idx"].(int32)) < int(c[j]["idx"].(int32))
}

// XXX 这个函数超级复杂，现在算是测完了，最终结果输出是个按每层idx有序的[*gin.H]树的列表
//     谁要是蛋疼动这里，首先要自己测完再提交，不然要烦死
func (*MenuProc) TransMenu(raws []*go_micro_srv_auth.RouteData) []*gin.H {
	var routes []*gin.H
	var root menuItem2
	children := map[string]menuItem2{}
	var tmp gin.H
	// 整理成单级索引表，所以是父节点uid
	for _, v := range raws {
		tmp = gin.H{
			"id":             v.UrlUid,
			"path":           v.Url,
			"name":           v.UrlName,
			"title":          v.Title,
			"icon":           v.Icon,
			"component_name": v.ComponentName,
			"idx":            v.Idx,
		}

		if "" != v.ComponentPath {
			tmp["component_path"] = v.ComponentPath
		}

		if "" != v.Redirect {
			tmp["redirect"] = gin.H{
				"name": v.Redirect,
			}
		}

		if 1 == v.IsHidden {
			tmp["hidden"] = true
		}

		if ss_sql.UUID == v.ParentUid {
			root = append(root, tmp)
		} else {
			if nil == children[v.ParentUid] {
				children[v.ParentUid] = append([]gin.H{}, tmp)
			} else {
				children[v.ParentUid] = append(children[v.ParentUid], tmp)
			}
		}
	}

	//ss_log.Error("children=%v\n", children)
	//ss_log.Error("root=%v\n", root)
	// 从上到下的树生成，注意每个地方都得做个排序
	// 循环最多10+1层，如果觉得还是太多的话，再改吧
	// 现在也就3层最多
	var tmpParent menuItem
	var tmpParent2 menuItem
	sort.Sort(root)
	for k, _ := range root {
		tmpParent = append([]*gin.H{}, &root[k])
		// 最多10+1级
		for i := 0; i < 10 && len(tmpParent) > 0; i++ {
			//ss_log.Error("i=%v\n", i)
			tmpParent2 = []*gin.H{}
			//ss_log.Error("tmpParent=%v\n", tmpParent)
			sort.Sort(tmpParent)
			for k2, _ := range tmpParent {
				//ss_log.Error("v2=%v\n", tmpParent[k2])

				sort.Sort(children[(*tmpParent[k2])["id"].(string)])
				for k3, _ := range children[(*tmpParent[k2])["id"].(string)] {
					//ss_log.Error("v3=%v|&v3=%v\n", v3, &v3)
					// 拼多级路径
					children[(*tmpParent[k2])["id"].(string)][k3]["path"] = (*tmpParent[k2])["path"].(string) + children[(*tmpParent[k2])["id"].(string)][k3]["path"].(string)
					if nil == (*tmpParent[k2])["children"] {
						(*tmpParent[k2])["children"] = append([]*gin.H{}, &children[(*tmpParent[k2])["id"].(string)][k3])
					} else {
						(*tmpParent[k2])["children"] = append((*tmpParent[k2])["children"].([]*gin.H), &children[(*tmpParent[k2])["id"].(string)][k3])
					}
					// 广度优先，所以要记录每层的数据，作为下层的父节点
					//ss_log.Error("tmpParent2|before|1=%v\n", tmpParent2, &children[(*tmpParent[k2])["id"].(string)][k3])
					tmpParent2 = append(tmpParent2, &children[(*tmpParent[k2])["id"].(string)][k3])
					//ss_log.Error("tmpParent2|before|2=%v\n", tmpParent2, &children[(*tmpParent[k2])["id"].(string)][k3])
				}

				//ss_log.Error("tmpParent2|after=%v\n", tmpParent2)
			}
			tmpParent = tmpParent2
			//ss_log.Error("len(tmpParent)|after=%v\n", len(tmpParent))
		}

		routes = append(routes, &root[k])
	}

	return routes
}
