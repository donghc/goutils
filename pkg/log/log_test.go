package log

import (
	"fmt"
	"github.com/donghc/goutils/pkg/alerthook"
	"github.com/donghc/goutils/pkg/log/writer"
	"go.uber.org/zap"
	"testing"

	"github.com/hashicorp/go-retryablehttp"
)

func TestNewCustomLogger_Hook(t *testing.T) {
	var (
		options []zap.Option
	)
	//var count int64 = 1
	retryClient := retryablehttp.NewClient()
	retryClient.RetryMax = 3
	hook := alerthook.NewDingDingHook(
		retryClient.StandardClient(),
		"报警",
		"SEC3fbb876618e19ba2150b2869e18246667287450b41ea5c2bfc2e1193dce9d860", []string{}, "123", "baidu")

	options = append(options, zap.Hooks(hook.ZapHook()))
	ws := writer.BuildWriteSyncer([]string{"stdout"})

	logger := NewCustomLogger("info", ws, options)
	str := "代码错误"
	fmt.Println()
	logger.Panicf(str)
}
