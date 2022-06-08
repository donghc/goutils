package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os/exec"
	"regexp"
	"time"
)

// https://work.weixin.qq.com/api/doc/90000/90136/91770
type workWxMsgContent struct {
	Content       string   `json:"content"`
	MentionedList []string `json:"mentioned_list"`
	MobileList    []string `json:"mentioned_mobile_list"`
	IP            string   `json:"ip"`
	HostName      string   `json:"hostname"`
}

type workWxMsg struct {
	MsgType string           `json:"msgtype"` // text/markdown/image/news
	Text    workWxMsgContent `json:"text"`
}

func wxNotify(subject, text string) {
	var msg workWxMsg
	msg.MsgType = "text"
	txt := &msg.Text

	txt.Content = subject + "\n" + text
	if len(mentioned) > 0 {
		txt.MobileList = mentioned
	} else {
		txt.MentionedList = []string{"@all"}
	}
	body, _ := json.Marshal(msg)
	req, err := http.NewRequest("POST",
		wxURI,
		bytes.NewBuffer(body))

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("post workWx:", err)
	}
	defer resp.Body.Close()
}

var (
	ctx       = context.Background()
	client    = &http.Client{}
	wxURI     = `https://qyapi.weixin.qq.com/cgi-bin/webhook/send?key=3e689b2c-fa1f-4f85-9462-1674dab4a0cf`
	mentioned []string
	max       int64
	count     int64 = 0
	notice          = 1
)
var (
	EngineFormat = map[string]string{
		"AhnLab":      "not scan",
		"ALYac":       "not scan",
		"Gridinsoft":  "not scan",
		"Filseclab":   "not scan",
		"Systweak":    "not scan",
		"TACHYON":     "not scan",
		"GDATA":       "not scan",
		"Panda":       "not scan",
		"Kingsoft":    "not scan",
		"Baidu":       "not scan",
		"Emsisoft":    "not scan",
		"Qihu360":     "not scan",
		"Comodo":      "not scan",
		"QuickHeal":   "not scan",
		"Xvirus":      "not scan",
		"JiangMin":    "not scan",
		"Antiy":       "not scan",
		"Sunbelt":     "not scan",
		"Avast":       "not scan",
		"AVG":         "not scan",
		"NANO":        "not scan",
		"Sangfor":     "not scan",
		"Rising":      "not scan",
		"Arcabit":     "not scan",
		"Avira":       "not scan",
		"ClamAV":      "not scan",
		"Cyren":       "not scan",
		"DrWeb":       "not scan",
		"IKARUS":      "not scan",
		"K7":          "not scan",
		"F-Prot":      "not scan",
		"Fortinet":    "not scan",
		"F-Secure":    "not scan",
		"McAfee":      "not scan",
		"TrendMicro":  "not scan",
		"VirusBuster": "not scan",
		"VBA32":       "not scan",
		"ESET":        "not scan",
	}
)

func GetVersion() (string, error) {
	buf := &bytes.Buffer{}
	param := "-v"
	//cmd := exec.Cmd{Path: c.ScriptPath, Args: append([]string{param}, input...), Stdout: buf}
	cmd := exec.Command("/opt/scancl/scancl", param)
	cmd.Stdout = buf
	if err := cmd.Run(); err != nil {
		return "", err
	}
	reg, err := regexp.Compile(`VDF Version:\s+([0-9.]+)`)
	if err != nil {
		return "", err
	}
	m := reg.FindAllStringSubmatch(buf.String(), -1)
	fmt.Println("antivir version：", m[0][1])
	return m[0][1], nil
}

type Samp struct {
	T bool   `json:"t"`
	M string `json:"m"`
}

func main() {
	//location, err2 := time.LoadLocation("Local")
	//fmt.Println(err2)
	//fmt.Println(location)
	//parse, err := time.ParseInLocation("2006-01-02 15:04:05", "2010-11-20 21:11:43.057 +0800 CST", location)
	//fmt.Println(err)
	//fmt.Println(parse.String())
	timeFormat := "2006-01-02 15:04:05 +0800 CST"
	time2, er := time.Parse(timeFormat, "2010-11-20 21:11:43.057 +0800 CST") //洛杉矶时间
	fmt.Println(er)
	fmt.Println(time2)
	fmt.Println(time2.Format(timeFormat))

	fmt.Println(time.Now())
	fmt.Println(time.Now().Format(timeFormat))

	//file, err := os.Open("D:\\software\\Apifox\\ffmpeg.dll")
	//defer file.Close()
	//
	//body := &bytes.Buffer{}
	//writer := multipart.NewWriter(body)
	//part, err := writer.CreateFormFile("file", "D:\\software\\Apifox\\ffmpeg.dll")
	//
	//_, err = io.Copy(part, file)
	//
	//err = writer.WriteField("KEY", "abcdefgh")
	//err = writer.Close()
	//
	////data := strings.NewReader(`KEY=abcdefgh`)
	////生成post请求
	//req, err := http.NewRequest("POST", "https://lsb102.threatbook-inc.cn/api/file/submit", body)
	//if err != nil {
	//	log.Fatal(err)
	//}
	//req.Header.Set("Content-Type", writer.FormDataContentType())
	//
	////req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	//req.Header.Set("authorization", "Basic bHNidGVzdDpsc2J0ZXN0")
	//req.Header.Set("authority", "lsb102.threatbook-inc.cn")
	//resp, err := client.Do(req)
	//if err != nil {
	//	log.Fatal(err)
	//}
	//bodyText, err := ioutil.ReadAll(resp.Body)
	//if err != nil {
	//	log.Fatal(err)
	//}
	//fmt.Printf("%s\n", bodyText)
}

func f1() (result int) {
	defer func() {
		result++
	}()

	return 10
}

func f2() (r int) {
	t := 5
	defer func() {
		r = t + 1
	}()
	return 2 * t
}
func main1() {
	a := []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}
	b := a[2:5]
	c := b[2:6:7]
	c = append(c, 100)
	c = append(c, 200)
	b[2] = 20
	fmt.Println(b)
	fmt.Println(c)
	fmt.Println(a)

	//version, err := GetVersion()
	//fmt.Println(err)
	//fmt.Println(version)

	//for {
	//	time.Sleep(time.Second * 3)
	//	cmd := exec.CommandContext(ctx, "tail", "-1", "/var/log/messages")
	//	output, _ := cmd.Output()
	//	r := string(output)
	//	log.Println(r)
	//	//r := "Mar 28 10:54:58 10-52-6-85 csinterface: [GIN] 2022/03/28 - 10:54:58 | 200 | 11.288055ms | 10.10.14.79 | POST \"/v1/file/sha256/batch\""
	//	if strings.Contains(strings.ToLower(r), "[gin]") &&
	//		strings.Contains(strings.ToLower(r), "csinterface") &&
	//		strings.Contains(strings.ToLower(r), "/v1/file/sha256/") {
	//		split := strings.Split(strings.ToLower(r), "[gin]")
	//		s := split[1]
	//		rs := strings.Split(s, "|")
	//		//fmt.Println(strings.TrimSpace(rs[0]))
	//		timeTemplate := "2006/01/02 - 15:04:05"
	//		stamp, _ := time.ParseInLocation(timeTemplate, strings.TrimSpace(rs[0]), time.Local)
	//		subSec := time.Now().Sub(stamp).Seconds()
	//		if subSec < 3 {
	//			count++
	//			log.Println("累计次数 ", count)
	//		}
	//		if count >= max {
	//			//开始报警
	//			log.Println("wx 报警开始")
	//			wxNotify("云查接口正在被调用", "")
	//			notice++
	//		}
	//		if notice > 3 {
	//			log.Println("wx 报警阈值达到 3次")
	//			time.Sleep(time.Minute * 30)
	//			notice = 1
	//		}
	//	} else {
	//		log.Println("not gin log")
	//	}
	//}
}
