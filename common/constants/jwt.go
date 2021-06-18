package constants

import "time"

const (
	// jwt的有效期
	JwtLifeTime = 3 * 24 * time.Hour // 3天

	// 客户端刷新token的时间间隔 (单位: 分钟)
	RefreshTokenInterval int32 = 60 * 6 // 6小时
)
