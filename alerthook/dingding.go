package alerthook

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

type DingDing struct {
	Security string `json:"security"`
	WebHook  string `json:"web_hook"`
	KeyWords string `json:"key_words"`
}

func hmacSha256(stringToSign string, secret string) string {
	h := hmac.New(sha256.New, []byte(secret))
	h.Write([]byte(stringToSign))
	return base64.StdEncoding.EncodeToString(h.Sum(nil))
}

// Sign 获取加密后的url地址
func (dd *DingDing) sign() string {
	timestamp := time.Now().UnixNano() / 1e6
	stringToSign := fmt.Sprintf("%d\n%s", timestamp, dd.Security)
	sign := hmacSha256(stringToSign, dd.Security)
	url := fmt.Sprintf("%s&timestamp=%d&sign=%s", dd.WebHook, timestamp, sign)
	return url
}

// SendDingDingMessage 发送消息
func (dd *DingDing) SendDingDingMessage(contentData string) bool {
	defer func() {
		if err := recover(); err != nil {
			log.Println("err : ", err)
		}
	}()
	if checkLimit() {
		log.Println("钉钉消息超过限制，不发送。")
		return false
	}
	var atMap = make(map[string]string)
	atMap["isAtAll"] = "true"
	atMap["msgtype"] = "text"
	content, data := make(map[string]string), make(map[string]interface{})
	content["content"] = dd.KeyWords + contentData
	data["msgtype"] = "text"
	data["text"] = content
	data["at"] = atMap
	b, _ := json.Marshal(data)

	resp, err := http.Post(dd.sign(), "application/json", bytes.NewBuffer(b))
	if err != nil {
		log.Println(err)
		return false
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {

		}
	}(resp.Body)
	return true
}

//check 检查消息发送频率是否超过限制：每个机器人每分钟最多发送20条。如果超过20条，会限流10分钟。
func checkLimit() bool {
	//发送前设置一个缓存一分钟的自增数
	//value := 10
	//if value > 20 {
	//	return true
	//}
	return false
}
