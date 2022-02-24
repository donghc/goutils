package writer

import (
	"github.com/donghc/goutils/pkg/pubsub"
	"github.com/wjiec/gdsn"
	"io"
	"os"
	"strings"
)

//BuildWriteSyncer 创建writer
func BuildWriteSyncer(out []string) []io.Writer {
	if len(out) == 0 {
		return nil
	}
	var ws []io.Writer
	for i := range out {
		str := strings.ToLower(out[i])
		if strings.HasPrefix(str, "kafka") {
			syncProducer, _ := pubsub.CreateKafkaPublisher(str)
			d, _ := gdsn.Parse(str)
			topic := d.Query().Get("topic")
			ws = append(ws, NewKafkaWriter(topic, syncProducer))
		} else if strings.HasPrefix(str, "stdout") {
			ws = append(ws, os.Stdout)
		} else if strings.HasPrefix(str, "stderr") {
			ws = append(ws, os.Stderr)
		} else {
			f, _ := os.Open(str)
			ws = append(ws, f)
		}
	}
	return ws
}
