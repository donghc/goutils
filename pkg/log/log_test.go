package log

import (
	"github.com/donghc/goutils/pkg/log/writer"
	"go.uber.org/zap"
	"testing"

	"github.com/donghc/goutils/pkg/alerthook"
	"github.com/hashicorp/go-retryablehttp"
)

func TestNewCustomLogger_Hook(t *testing.T) {
	var (
		options []zap.Option
	)
	//var count int64 = 1
	retryClient := retryablehttp.NewClient()
	retryClient.RetryMax = 3
	hook := alerthook.NewWorkWxHook(
		retryClient.StandardClient(),
		"interface告警",
		"", []string{"18510091324"}, "", "")

	options = append(options, zap.Hooks(hook.ZapHook()))
	ws := writer.BuildWriteSyncer([]string{"stdout"})

	logger := NewCustomLogger("", ws, options)

	logger.Info("info hello world")
	logger.Errorf("err hello world")
	logger.Panicf("panic hello world")
}
