package progressbar

import (
	"crypto/rand"
	"github.com/cheggaaa/pb/v3"
	"io"
	"io/ioutil"
	"os"
	"testing"
)

func TestDemo(t *testing.T) {
	count := 100000

	// 创建进度条并开始
	//bar := pb.StartNew(count)

	tmpl := `{{ red "With funcs:" }} {{ bar . "<" "-" (cycle . "↖" "↗" "↘" "↙" ) "." ">"}} {{speed . | rndcolor }} {{percent .}} {{string . "my_green_string" | green}} {{string . "my_blue_string" | blue}}`

	// 开始基于我们的模板的进度条
	bar := pb.ProgressBarTemplate(tmpl).Start(count)

	// 设置字符串元素的值
	bar.Set("my_green_string", "green").Set("my_blue_string", "blue")
	//bar.SetTemplate("")
	bar.SetWriter(os.Stdout)
	for i := 0; i < count; i++ {
		bar.Increment()
	}

	// 结束进度条
	bar.Finish()
}

func TestRunPb(t *testing.T) {
	var limit int64 = 1024 * 1024 * 500

	// 我们将把500MiB从/dev/rand复制到/dev/null
	reader := io.LimitReader(rand.Reader, limit)
	writer := ioutil.Discard

	// 开始进度条
	bar := pb.Full.Start64(limit)

	// 创建代理读取器
	barReader := bar.NewProxyReader(reader)

	// 从代理读取器复制
	io.Copy(writer, barReader)

	// 结束进度条
	bar.Finish()

}
