package test

import (
	"a.a/cu/ss_log"
	"a.a/cu/strext"
	"a.a/mp-server/common/cache"
	"a.a/mp-server/common/constants"
	go_micro_srv_cust "a.a/mp-server/common/proto/cust"
	"a.a/mp-server/common/ss_count"
	"a.a/mp-server/cust-srv/handler"
	"fmt"
	jsoniter "github.com/json-iterator/go"
	"sort"
	"strings"
	"testing"
	"time"
)

func TestLatLng(t *testing.T) {
	posSn := "ax-001"
	// ======岭南创意园=======
	lat2 := "23.0795382400"  // 维度
	lng2 := "113.3406764300" // 经度
	// ======岭南创意园=======
	// ======创投小镇=======
	//lat3 := "23.0714778800"  // 维度
	//lng3 := "113.3125147400" // 经度
	// ======创投小镇=======
	redisKey := "pos_sn_" + posSn
	var lat, lng, scope string
	var err error
	//if value, _ := ss_data.GetPosSnFromCache(redisKey, cache.RedisCli, constants.DefPoolName); value == "" { // 查询数据库并设置进redis
	if value, _ := cache.RedisClient.Get(redisKey).Result(); value == "" { // 查询数据库并设置进redis

		lat, lng, scope, err = handler.GetPosLatLng(posSn)
		if err != nil {
			ss_log.Error("%s", err.Error())
			return
		}
		redisValue := fmt.Sprintf("%s,%s,%s", lat, lng, scope)
		// 设进redis
		//if err := ss_data.SetPosSnToCache(redisKey, redisValue, cache.RedisCli, constants.DefPoolName); err != nil {
		if err := cache.RedisClient.Set(redisKey, redisValue, constants.PosNoKeySecV2); err != nil {
			ss_log.Error("经纬度存进redis失败,posNo--->%s,lat--->%s,lng--->%s,scope--->%s", posSn, lat, lng, scope)
		}
	} else {
		split := strings.Split(value, ",")
		lat = split[0]
		lng = split[1]
		scope = split[2]
	}
	// 计算距离
	distance := ss_count.CountCircleDistance(strext.ToFloat64(lat), strext.ToFloat64(lng), strext.ToFloat64(lat2), strext.ToFloat64(lng2))
	if distance > strext.ToFloat64(scope) {
		ss_log.Error("pos机超出使用范围,计算的范围是--->%v,限定的范围是--->%v", distance, scope)
		return
	}
	ss_log.Error("通过----------->数据库距离为--->%s,计算的距离为--->%v", scope, distance)
}

func TestMap(t *testing.T) {
	var m map[string]interface{}
	m = make(map[string]interface{})
	m["a"] = "a"
	fmt.Println(m)
	m["a"] = 1
	fmt.Println(m)
}

type sort1 struct {
	A string
	B float64
}

func TestQuickSort(t *testing.T) {
	s := sorts{}
	s1 := &sort1{A: "a", B: 1.3}
	s = append(s, s1)
	s2 := &sort1{A: "a", B: 1.1}
	s = append(s, s2)
	s3 := &sort1{A: "a", B: 1.2}
	s = append(s, s3)
	s4 := &sort1{A: "a", B: 1.6}
	s = append(s, s4)
	s5 := &sort1{A: "a", B: 1.5}
	s = append(s, s5)

	toString, _ := jsoniter.MarshalToString(s)
	fmt.Println("排序前--->", toString)
	//quickSort1(s, 0, len(s)-1)
	sort.Sort(sorts{})
	toString1, _ := jsoniter.MarshalToString(s)
	fmt.Println("排序后--->", toString1)

}

type sorts []*sort1

//Len()
func (s sorts) Len() int {
	return len(s)
}

//Less(): 成绩将有低到高排序
func (s sorts) Less(i, j int) bool {
	return s[i].B < s[j].B
}

//Swap()
func (s sorts) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

// quickSort 使用快速排序算法，排序整型数组
func quickSort1(values []*sort1, left int, right int) {
	if left < right {
		// 设置分水岭
		temp := values[left]

		// 设置哨兵
		i, j := left, right
		for {
			// 从右向左找，找到第一个比分水岭小的数
			for values[j].B >= temp.B && i < j {
				j--
			}

			// 从左向右找，找到第一个比分水岭大的数
			for values[i].B <= temp.B && i < j {
				i++
			}

			// 如果哨兵相遇，则退出循环
			if i >= j {
				break
			}

			// 交换左右两侧的值
			values[i], values[j] = values[j], values[i]
		}

		// 将分水岭移到哨兵相遇点
		values[left] = values[i]
		values[i] = temp

		// 递归，左右两侧分别排序
		quickSort1(values, left, i-1)
		quickSort1(values, i+1, right)
	}
}
func TestQest(t *testing.T) {
	// 测试代码
	arr := []int{9, 8, 7, 6, 5, 1, 2, 11, 4, 0}
	fmt.Println(arr)
	//quickSort1(arr, 0, len(arr)-1)
	fmt.Println(arr)
}

// 快速排序
func quickSort(values []*go_micro_srv_cust.NearbyServicerData, left int, right int) {
	if left < right {
		// 设置分水岭
		temp := values[left]
		// 设置哨兵
		i, j := left, right
		for {
			// 从右向左找，找到第一个比分水岭小的数
			for strext.ToFloat64(values[j].Distance) >= strext.ToFloat64(temp.Distance) && i < j {
				j--
			}
			// 从左向右找，找到第一个比分水岭大的数
			for strext.ToFloat64(values[i].Distance) <= strext.ToFloat64(temp.Distance) && i < j {
				i++
			}
			// 如果哨兵相遇，则退出循环
			if i >= j {
				break
			}
			// 交换左右两侧的值
			values[i], values[j] = values[j], values[i]
		}
		// 将分水岭移到哨兵相遇点
		values[left] = values[i]
		values[i] = temp
		// 递归，左右两侧分别排序
		quickSort(values, left, i-1)
		quickSort(values, i+1, right)
	}
}

func TestCountSrvCoordinate(t *testing.T) {
	//handler.CountSrvCoordinate(165.33, 33.666)
	nTime := time.Now()
	logDay := nTime.Format("2006-01-02")
	fmt.Println(logDay)
}
