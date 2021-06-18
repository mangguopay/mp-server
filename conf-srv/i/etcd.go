package i

import (
	"a.a/cu/strext"
	"a.a/mp-server/common/constants"
	"fmt"
	"github.com/micro/go-micro/config"
	"github.com/micro/go-micro/store"
	"github.com/micro/go-micro/store/etcd"
	"strings"
)

// 将配置信息载入etcd中
func StoreToEtcd(etcdAddrs []string) error {
	s := etcd.NewStore(store.Nodes(etcdAddrs...))

	l := []*store.Record{}
	prefix := constants.ETCDPrefix

	r, _ := s.List()
	for _, v := range r {
		if strings.HasPrefix(v.Key, prefix) {
			s.Delete(v.Key)
		}
	}

	base := map[string]interface{}{}
	config.Get("mp-server", "base").Scan(&base)
	l = append(l,
		&store.Record{
			// 必须两层路径，否则接收方会拆解kv，少了一层
			Key:   fmt.Sprintf("%s/%s", prefix, "base/base"),
			Value: []byte(strext.ToJson(base)),
		},
	)
	// 数据库
	dbs := config.Get("mp-server", "database")
	a := map[string]map[string]interface{}{}
	_ = dbs.Scan(&a)
	for k, m := range a {
		b := &store.Record{
			Key:   fmt.Sprintf("%s/database/%s", prefix, k),
			Value: []byte(strext.ToStringNoPoint(m)),
		}
		l = append(l, b)
	}
	// cache
	caches := config.Get("mp-server", "cache")
	b := map[string]map[string]interface{}{}
	_ = caches.Scan(&b)
	for k, m := range b {
		host := &store.Record{
			Key:   fmt.Sprintf("%s/cache/%s", prefix, k),
			Value: []byte(strext.ToStringNoPoint(m)),
		}
		l = append(l, host)
	}
	// mq
	mqs := config.Get("mp-server", "mq")
	c := map[string]map[string]interface{}{}
	_ = mqs.Scan(&c)
	for k, m := range c {
		host := &store.Record{
			Key:   fmt.Sprintf("%s/mq/%s", prefix, k),
			Value: []byte(strext.ToStringNoPoint(m)),
		}
		l = append(l, host)
	}
	// message
	messages := config.Get("mp-server", "message")
	d := map[string]map[string]interface{}{}
	_ = messages.Scan(&d)
	for k, m := range d {
		host := &store.Record{
			Key:   fmt.Sprintf("%s/message/%s", prefix, k),
			Value: []byte(strext.ToStringNoPoint(m)),
		}
		l = append(l, host)
	}

	// aws-s3
	aws := config.Get("mp-server", "aws")
	e := map[string]map[string]interface{}{}
	_ = aws.Scan(&e)
	for k, m := range e {
		host := &store.Record{
			Key:   fmt.Sprintf("%s/aws/%s", prefix, k),
			Value: []byte(strext.ToStringNoPoint(m)),
		}
		l = append(l, host)
	}

	return s.Write(l...)
}
