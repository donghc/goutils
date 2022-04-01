package main

import (
	"bytes"
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"path/filepath"
	"strings"
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

func main() {
	m := make(map[string]string)
	fmt.Println(len(m))
	var n map[string]string = nil
	fmt.Println(len(n))

	d := filepath.Join("/data", "sadeq", "fdsg")
	all := strings.ReplaceAll(d, "\\", "/")
	fmt.Println(all)
	fmt.Println(time.Now().Unix())

	flag.Int64Var(&max, "max", 10, "set max count , default 10")
	flag.Parse()
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
