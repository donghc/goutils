package browser

import (
	"fmt"
	"testing"
)

func TestOpen(t *testing.T) {
	err := Open("www.baidu.com")

	fmt.Println(err)
}
