package evaluator

func NewLoginPasswordWrongItem(uid string) Item {
	return &LoginPasswordWrongItem{uid: uid}
}

// 登录密码错评估项
type LoginPasswordWrongItem struct {
	uid string
}

// 评估项名称
func (l *LoginPasswordWrongItem) Name() string {
	return ItemNameLoginPasswordWrong
}

// 对登录密码错误进行评估
func (l *LoginPasswordWrongItem) Evaluate() (int, error) {
	// todo
	return 0, nil
}

/*
本函数为开发调试时的示例代码
func (l *LoginPasswordWrongItem) Evaluate() (int, error) {
	sd, _ := time.ParseDuration("-24h")
	nTime := time.Now().Add(sd * 3)
	lastTime := ss_time.ForPostgres(nTime)

	list, err := dao.UserLoginLogDaoInstance.GetLastByTime(l.uid, lastTime)
	if err != nil {
		return 0, err
	}

	wrongNum := 0
	for _, v := range list {
		if v.Result == dao.LoginResultPassWrong {
			wrongNum++
		}
	}

	riskScore := 0

	switch {
	case wrongNum < 3:
		riskScore += 0
	case wrongNum >= 3 && wrongNum < 8:
		riskScore += 5
	case wrongNum > 8:
		riskScore += 10
	}

	return riskScore, nil
}
*/
