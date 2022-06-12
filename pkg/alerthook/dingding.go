package alerthook

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/sirupsen/logrus"
	"go.uber.org/zap/zapcore"
)

type dingdingHook struct {
	client    *http.Client
	subject   string //keywords 钉钉机器人需要关键字
	key       string
	apiURL    string
	mentioned []string
	ip        string
	hostname  string
}

//https://open.dingtalk.com/document/app#/serverapi2/qf2nxq

type dingdingMsg struct {
	MsgType string             `json:"msgtype"` // text/markdown/image/news
	Text    dingdingMsgContent `json:"text"`
	At      at                 `json:"at"`
}

type dingdingMsgContent struct {
	Content string `json:"content"`
}

type at struct {
	AtMobiles []string `json:"atMobiles,omitempty"`
	AtUserIds []string `json:"atUserIds,omitempty"`
	IsAtAll   bool     `json:"isAtAll"`
}

func NewDingDingHook(client *http.Client, subject string, key string, mentioned []string, ip, hostname string) *dingdingHook {
	if subject == "" {
		subject = os.Args[0]
	}
	if client == nil {
		client = &http.Client{}
	}

	return &dingdingHook{
		subject:   subject,
		key:       key,
		client:    &http.Client{},
		apiURL:    "https://oapi.dingtalk.com/robot/send?access_token=ba55d89e0d0bd88d927fff578f8a89c2dffe61ec98194105bce7bce04938a629",
		mentioned: mentioned,
		ip:        ip,
		hostname:  hostname,
	}
}

func (hook *dingdingHook) Levels() []logrus.Level {
	return []logrus.Level{logrus.PanicLevel, logrus.ErrorLevel}
}

func (hook *dingdingHook) Fire(entry *logrus.Entry) error {
	return hook.sendLogrus(entry)
}

func (hook *dingdingHook) sendLogrus(entry *logrus.Entry) error {
	if hook.key == "" {
		return nil
	}
	client := hook.client

	var msg workWxMsg
	msg.MsgType = "text"
	txt := &msg.Text
	txt.Content = hook.subject + "\n" + entry.Message
	if len(hook.mentioned) > 0 {
		txt.MobileList = hook.mentioned
	} else {
		txt.MentionedList = []string{"@all"}
	}
	body, _ := json.Marshal(msg)

	req, err := http.NewRequest("POST",
		hook.apiURL,
		bytes.NewBuffer(body))

	resp, err := client.Do(req)
	if err != nil {
		log.Println("post workWx:", err)
	}
	defer resp.Body.Close()
	return nil
}

func (hook *dingdingHook) sendZap(entry zapcore.Entry) error {
	if hook.key == "" {
		return nil
	}
	var flag bool
	for _, level := range hook.Levels() {
		if level.String() == entry.Level.String() {
			flag = true
			break
		}
	}
	if !flag {
		return nil
	}
	client := hook.client

	if len(entry.Message) > 2048 {
		entry.Message = entry.Message[0:2048]
	}
	content := &content{
		Entry:    entry,
		Subject:  hook.subject,
		IP:       hook.ip,
		HostName: hook.hostname,
	}
	contentByte, _ := json.Marshal(content)

	var at at
	if len(hook.mentioned) == 0 {
		at.IsAtAll = true
	} else {
		at.AtMobiles = hook.mentioned
		at.IsAtAll = false
	}
	var msg dingdingMsg
	msg.MsgType = "text"
	msg.At = at

	txt := &msg.Text
	txt.Content = string(contentByte)

	body, _ := json.Marshal(msg)

	req, err := http.NewRequest("POST",
		hook.Sign(),
		bytes.NewBuffer(body))

	req.Header.Set("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		log.Println("post sms:", err)
	}
	defer resp.Body.Close()
	return nil
}

func (hook *dingdingHook) LogrusHook() logrus.Hook {
	return hook
}

func (hook *dingdingHook) ZapHook() func(zapcore.Entry) error {
	return hook.sendZap
}

func hmacSha256(stringToSign string, secret string) string {
	h := hmac.New(sha256.New, []byte(secret))
	h.Write([]byte(stringToSign))
	return base64.StdEncoding.EncodeToString(h.Sum(nil))
}

// Sign 获取加密后的url地址
func (hook *dingdingHook) Sign() string {
	secret := hook.key
	webhook := hook.apiURL
	timestamp := time.Now().UnixNano() / 1e6
	stringToSign := fmt.Sprintf("%d\n%s", timestamp, secret)
	sign := hmacSha256(stringToSign, secret)
	url := fmt.Sprintf("%s&timestamp=%d&sign=%s", webhook, timestamp, sign)
	return url
}

// SendDingDingMessage 发送消息
//func (hook *dingdingHook) SendDingDingMessage(contentData string) bool {
//	var atMap = make(map[string]string)
//	atMap["isAtAll"] = "true"
//	atMap["msgtype"] = "text"
//	content, data := make(map[string]string), make(map[string]interface{})
//	content["content"] = hook.subject + contentData
//	data["msgtype"] = "text"
//	data["text"] = content
//	data["at"] = atMap
//	b, _ := json.Marshal(data)
//
//	defer resp.Body.Close()
//	//body, _ := ioutil.ReadAll(resp.Body)
//	//fmt.Println(string(body))
//	return true
//}

//check 检查消息发送频率是否超过限制：每个机器人每分钟最多发送20条。如果超过20条，会限流10分钟。
