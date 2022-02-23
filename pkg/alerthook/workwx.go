package alerthook

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/sirupsen/logrus"
	"go.uber.org/zap/zapcore"
)

// https://qyapi.weixin.qq.com/cgi-bin/webhook/send?key=a544c793-4800-4e12-b739-66fc47250dfa

type workWxHook struct {
	client    *http.Client
	subject   string
	key       string
	apiURL    string
	mentioned []string
	ip        string
	hostname  string
}

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

type content struct {
	zapcore.Entry
	Subject  string
	IP       string
	HostName string
}

func NewWorkWxHook(client *http.Client, subject string, key string, mentioned []string, ip, hostname string) *workWxHook {
	if subject == "" {
		subject = os.Args[0]
	}
	if client == nil {
		client = &http.Client{}
	}

	return &workWxHook{
		subject:   subject,
		key:       key,
		client:    &http.Client{},
		apiURL:    "https://qyapi.weixin.qq.com/cgi-bin/webhook/send?key=" + key,
		mentioned: mentioned,
		ip:        ip,
		hostname:  hostname,
	}
}

func (self *workWxHook) Levels() []logrus.Level {
	return []logrus.Level{logrus.PanicLevel}
}

func (self *workWxHook) Fire(entry *logrus.Entry) error {
	return self.sendLogrus(entry)
}

func (self *workWxHook) sendLogrus(entry *logrus.Entry) error {
	if self.key == "" {
		return nil
	}
	client := self.client

	var msg workWxMsg
	msg.MsgType = "text"
	txt := &msg.Text
	txt.Content = self.subject + "\n" + entry.Message
	if len(self.mentioned) > 0 {
		txt.MobileList = self.mentioned
	} else {
		txt.MentionedList = []string{"@all"}
	}
	body, _ := json.Marshal(msg)

	req, err := http.NewRequest("POST",
		self.apiURL,
		bytes.NewBuffer(body))

	resp, err := client.Do(req)
	if err != nil {
		log.Println("post workWx:", err)
	}
	defer resp.Body.Close()
	return nil
}

func (self *workWxHook) sendZap(entry zapcore.Entry) error {
	if self.key == "" {
		return nil
	}
	var flag bool
	for _, level := range self.Levels() {
		if level.String() == entry.Level.String() {
			flag = true
			break
		}
	}
	if !flag {
		return nil
	}
	client := self.client

	ct := &content{
		Entry:    entry,
		Subject:  self.subject,
		IP:       self.ip,
		HostName: self.hostname,
	}
	marshal, _ := json.Marshal(ct)

	var msg workWxMsg
	msg.MsgType = "text"
	txt := &msg.Text
	txt.Content = string(marshal)
	if len(self.mentioned) == 0 {
		txt.MentionedList = []string{"@all"}
	} else {
		txt.MobileList = self.mentioned
	}
	body, _ := json.Marshal(msg)

	req, err := http.NewRequest("POST",
		self.apiURL,
		bytes.NewBuffer(body))

	resp, err := client.Do(req)
	if err != nil {
		log.Println("post sms:", err)
	}
	defer resp.Body.Close()
	return nil
}

func (self *workWxHook) LogrusHook() logrus.Hook {
	return self
}

func (self *workWxHook) ZapHook() func(zapcore.Entry) error {
	return self.sendZap
}
