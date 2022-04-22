package log

import (
	"fmt"
	"github.com/donghc/goutils/pkg/alerthook"
	"github.com/donghc/goutils/pkg/log/writer"
	"go.uber.org/zap"
	"math/rand"
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
	hook := alerthook.NewWorkWxHook(
		retryClient.StandardClient(),
		"interface告警",
		"a544c793-4800-4e12-b739-66fc47250dfa", []string{"18510091324"}, "", "")

	options = append(options, zap.Hooks(hook.ZapHook()))
	ws := writer.BuildWriteSyncer([]string{"stdout"})

	logger := NewCustomLogger("info", ws, options)
	str := "0123456789abcdefghijklmnopqrstuvwxyz"
	bytes := []byte(str)
	var result []byte
	for i := 0; i < 5120; i++ {
		result = append(result, bytes[rand.Intn(len(bytes))])
	}
	fmt.Println()
	logger.Panicf(string(result))
}
