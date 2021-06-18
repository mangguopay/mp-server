package main

import (
	"a.a/cu/ss_mail"
	"a.a/cu/strext"
	"bufio"
	"fmt"
	jsoniter "github.com/json-iterator/go"
	"github.com/tealeg/xlsx"
	"io"
	"io/ioutil"
	"os"
	"sort"
	"time"
)

type FileContentJsonStruct struct {
	Row          string
	ToAccount    string
	Amount       string
	CurrencyType string
}

func init() {
	fmt.Println("init start")
}

func main() {

	z, zErr := time.LoadLocation("Asia/Phnom_Penh")
	if zErr != nil {
		fmt.Printf("err = %v", zErr)
	}

	time1 := time.Now()
	fmt.Printf("time1 : %v \n", time1) //打印结果：

	timeIn1 := time1.In(z)
	fmt.Printf("timeIn1 : %v \n", timeIn1) //打印结果：

	// 设置time包中的默认时区
	time.Local = z

	time2 := time.Now()
	fmt.Printf("time2 : %v \n", time2) //打印结果：

	timeIn2 := time2.In(z)
	fmt.Printf("timeIn2 : %v \n", timeIn2) //打印结果：

	fmt.Println("----------------------------") //打印结果：

	timeString := "2020-11-11 00:00:00"
	newTime1, _ := time.ParseInLocation("2006-01-02 15:04:05", timeString, time.Local)
	fmt.Printf("newTime1 : %v \n", newTime1) //打印结果：
	newTime2, _ := time.Parse("2006-01-02 15:04:05", timeString)
	newTime2 = newTime2.In(time.Local)
	fmt.Printf("newTime2 : %v \n", newTime2) //打印结果：

	//timesss := time.Now()
	//intA := timesss.Unix()
	//fmt.Printf("timesss ----- > %v \n timeNowUnix  ---->  %v \n", timesss ,intA)
	//
	////settleDate := time.Unix(strext.ToInt64(aaa), 0).Format(ss_time.DateTimeDashFormat)
	//settleDate := time.Unix(intA, 0).Format(ss_time.DateTimeDashFormat)
	//fmt.Printf("settleDate ----->  %v \n ---------------------- \n", settleDate)
	//
	//time2, err := time.Parse(ss_time.DateTimeDashFormat, settleDate)
	//if err != nil {
	//	ss_log.Error("转换时间为时间戳出错,err[%v]", err)
	//}
	//
	//time2B := time2.In(z)
	//fmt.Printf("time2 ----->  %v \n unix   ------->  %v \n ", time2B, time2B.Unix())
	//time2BB := time.Unix(time2B.Unix(), 0).Format(ss_time.DateTimeDashFormat)
	//fmt.Printf("time22 ----->  %v \n", time2BB)

	//生成一个时间
	//oldtime := time.Now()
	//timeUnix := oldtime.Unix() //已知的时间戳
	//fmt.Printf("timeUnix : %v ------ \n", timeUnix) //打印结果：
	//formatTimeStr := time.Unix(timeUnix, 0).Format("2006-01-02 15:04:05")
	//fmt.Println(formatTimeStr) //打印结果：
	//
	////formatTime, err := time.Parse("2006-01-02 15:04:05", formatTimeStr)
	//formatTime, err := time.ParseInLocation("2006-01-02 15:04:05", formatTimeStr, time.Local)
	//if err == nil {
	//	fmt.Println(formatTime) //打印结果：
	//}
	//
	//fmt.Printf("formatTimeUnix : %v \n", formatTime.Unix()) //打印结果：

	//batchTest()
	//fmt.Printf("%v",strext.ToInt("1234"))
	//fmt.Printf(fmt.Sprintf("a is %v\n", &a.Row))

	//balance=3875600, ss_count.Add(req.Amount, fee)=500,

	//a := ss_count.Sub("10", "")
	//fmt.Printf("a  %v", a)

	//fmt.Printf("account %d ..", strext.ToInt("0855"))

	//s := []string{"1", "1", "3", "3", "2"}
	//
	////aaa(s)
	////
	//fmt.Println(strings.Join(s, ","))

	//s := []string{"1", "1", "3", "3", "2"}
	//ss := deduplication(s)
	//fmt.Println(ss)

	//s := "    我的 hello 1 world       "
	//
	////s =strings.TrimLeft(s," ")
	////s =strings.TrimRight(s," ")
	//s = strings.Trim(s," ")
	//s = strings.ToTitle(s)
	//
	//sArr := strings.Split(s," ")
	//s = strings.Join(sArr,"&&")
	//
	//
	//fmt.Println(s)

	//stringTest()

	//多种读取文件方式的时间比较
	//readTime()

	//生成json
	//createJson()

	//读xlsx文件
	//readXlsx()

	//发送邮件
	//sendMail()
}

func batchTest() {
	var datas []*FileContentJsonStruct
	for i := 0; i < 10; i++ {
		datas = append(datas, &FileContentJsonStruct{Row: strext.ToString(i)})
	}

	var s []string
	var s2 []string
	for k, v := range datas {
		fmt.Printf("k %v, len %v,  row %v \n", k, len(datas), v.Row)
		s = append(s, v.Row)
		if len(s) == 3 || k == len(datas)-1 {
			fmt.Printf("!!!!\n")
			s2 = append(s2, s...)
			s = s[0:0]
		}
	}

	fmt.Printf("s2 %v", s2)
}

func stringTest() {
	//s := "err_payment_pwd_1_asdeasda"
	//arr1 := strings.Split(s, "err_payment_pwd_")
	//fmt.Println("arr1:", arr1[1])
	//arr2 := strings.Split(arr1[1], "_")
	//fmt.Println("-----", len(arr2))
	//
	//ss := "0888123456"
	//ss = strings.TrimPrefix(ss, "0")
	//fmt.Println(ss)

	//s := "WEIXIN_MODERNPAY_"
	//switch {
	//case strings.HasPrefix(s, constants.TradeTypeModernpayPrefix):
	//	fmt.Printf("-----------1")
	//case strings.HasPrefix(s, constants.TradeTypeWeiXinPrefix):
	//	fmt.Printf("-----------2")
	//case strings.HasPrefix(s, constants.TradeTypeAlipayPrefix):
	//	fmt.Printf("-----------3")
	//default:
	//	fmt.Printf("产品[%v]交易类型未知[%v]", "req.SceneNo", s)
	//}
}

func readTime() {
	file := "run.log"
	start := time.Now()

	read0(file)
	t0 := time.Now()
	fmt.Printf("Cost time %v\n", t0.Sub(start))

	read1(file)
	t1 := time.Now()
	fmt.Printf("Cost time %v\n", t1.Sub(t0))

	read2(file)
	t2 := time.Now()
	fmt.Printf("Cost time %v\n", t2.Sub(t1))

	read3(file)
	t3 := time.Now()
	fmt.Printf("Cost time %v\n", t3.Sub(t2))

}

func read0(path string) string {
	f, err := ioutil.ReadFile(path)
	if err != nil {
		fmt.Printf("%s\n", err)
		panic(err)
	}
	return string(f)
}

func read1(path string) string {
	fi, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	defer fi.Close()

	chunks := make([]byte, 1024, 1024)
	buf := make([]byte, 1024)
	for {
		n, err := fi.Read(buf)
		if err != nil && err != io.EOF {
			panic(err)
		}
		if 0 == n {
			break
		}
		chunks = append(chunks, buf[:n]...)
	}
	return string(chunks)
}

func read2(path string) string {
	fi, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	defer fi.Close()
	r := bufio.NewReader(fi)

	chunks := make([]byte, 1024, 1024)

	buf := make([]byte, 1024)
	for {
		n, err := r.Read(buf)
		if err != nil && err != io.EOF {
			panic(err)
		}
		if 0 == n {
			break
		}
		chunks = append(chunks, buf[:n]...)
	}
	return string(chunks)
}

func read3(path string) string {
	fi, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	defer fi.Close()
	fd, err := ioutil.ReadAll(fi)
	return string(fd)
}

func createJson() {

	//生成文件的内容json字符串
	var temps []*FileContentJsonStruct

	temp := &FileContentJsonStruct{
		Row:          "1",
		ToAccount:    "a",
		Amount:       "a",
		CurrencyType: "a",
	}
	temp2 := &FileContentJsonStruct{
		Row:          "2",
		ToAccount:    "b",
		Amount:       "b",
		CurrencyType: "b",
	}

	temps = append(temps, temp)
	temps = append(temps, temp2)

	//文件内容json字符串
	fileContent := strext.ToJson(temps)

	var jsonToProto []*FileContentJsonStruct
	err := jsoniter.Unmarshal([]byte(fileContent), &jsonToProto)
	//err := json.Unmarshal([]byte(fileContent), &jsonToProto)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(fileContent)
	fmt.Println("------------")
	for _, v := range jsonToProto {
		fmt.Println(v.ToAccount)
	}
}

func sendMail() {
	//邮箱账号
	user := "aa44225533123@163.com"
	//user := "qq425312361@gmail.com"
	//注意，此处为授权码、不是密码
	password := "ccguydvywnlrwbxr"
	//smtp地址及端口
	host := "smtp.gmail.com:587"
	//接收者，内容可重复，邮箱之间用；隔开
	to := "425312361@qq.com"
	//邮件主题
	subject := "测试通过golang发送邮件"
	//邮件内容
	text := "你好！"
	body := `
	<html>
	<body>
	<h3>
	"测试通过golang发送邮件"` + text + `
	</h3>
	</body>
	</html>
	`
	fmt.Println("send email")

	err := ss_mail.SsMailInst.SendMail("11223344", user, password, host, to, subject, body, "html")

	//执行逻辑函数
	//err := sendMail("",user, password, host, to, subject, body, "html")
	if err != nil {
		fmt.Println("发送邮件失败!")
		fmt.Println(err)
	} else {
		fmt.Println("发送邮件成功!")
	}
}

//读xlsx文件
func readXlsx() {
	filepath := "./222.xls"
	file1, err1 := os.Open(filepath)
	if err1 != nil {
		fmt.Printf("open failed: %s\n", err1)
	}

	bb, err2 := ioutil.ReadAll(file1)
	if err2 != nil {
		fmt.Printf("open failed: %s\n", err2)
	}

	//xlsx.ReadZip()
	xlFile, err := xlsx.OpenBinary(bb)
	//xlFile, err := xlsx.OpenFile(filepath)
	if err != nil {
		fmt.Printf("open failed: %s\n", err)
	}

	for i, sheet := range xlFile.Sheets { //工作表
		fmt.Printf("i: %v  Sheet Name: %s\n", i, sheet.Name)
		for j, row := range sheet.Rows { //行
			fmt.Printf("--------------------------j %v ---", j)
			if len(row.Cells) == 0 { //去掉空行
				continue
			}
			fmt.Printf(" ------cell %v\n", len(row.Cells[0].String()))

			for k, cell := range row.Cells { //列

				//fmt.Printf("k %v, cell %v\n", k, cell)

				text := cell.String()
				fmt.Printf("k %v, text %s\n", k, text)
			}
		}
	}

}

func deduplication(s []string) (ss []string) {
	m := map[string]string{}

	for _, v := range s {
		l := len(m)
		m[v] = "0"
		if len(m) != l { // 加入map后，map长度变化，则元素不重复
			ss = append(ss, v)
		}
	}

	//正序
	sort.Strings(ss)

	//倒序
	sort.Sort(sort.Reverse(sort.StringSlice(ss)))

	return ss
}

func aaa(s []string) {
	fmt.Printf("aaa  s %p", s)
	s[0] = "222"
	fmt.Printf("aaa  s %p", s)
	s = append(s, "5") //append返回的是新的
	fmt.Printf("aaa  s %p", s)
}
