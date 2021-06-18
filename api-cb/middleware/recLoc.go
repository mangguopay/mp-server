package middleware

import (
	"a.a/cu/db"
	"a.a/cu/ss_log"
	constants2 "a.a/mp-server/common/constants"
	"a.a/mp-server/common/ss_sql"
	"github.com/gin-gonic/gin"
	"time"
)

type RecLocMw struct {
}

/**
 * 记录地理信息
 */
func (RecLocMw) SaveLoc() gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := c.ClientIP()
		ua := c.Request.UserAgent()
		t := c.GetInt64("RecLocMw.Clock.end")

		p := c.Query("p")
		q := c.Query("q")
		r := c.Query("r")

		dbHandler := db.GetDB(constants2.DB_CRM)
		defer db.PutDB(constants2.DB_CRM, dbHandler)

		err := ss_sql.Exec(dbHandler, `insert into score_html_stat(p,q,r,ip,ua,time,log_time)values($1,$2,$3,$4,$5,$6,current_timestamp)`, p, q, r, ip, ua, t)
		if err != nil {
			ss_log.Error("err=[%v]", err)
		}
	}
}

func (RecLocMw) BeginClock() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("RecLocMw.Clock", time.Now().UnixNano())
	}
}

func (RecLocMw) EndClock() gin.HandlerFunc {
	return func(c *gin.Context) {
		bt := c.GetInt64("RecLocMw.Clock")
		et := time.Now().UnixNano()
		t := et - bt
		c.Set("RecLocMw.Clock.end", t)
	}
}
