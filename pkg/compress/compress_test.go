package compress

import (
	"fmt"
	"testing"
)

func TestDoZlibCompress(t *testing.T) {
	s := "hello, gasdgdsahgfjhkl"
	zip := DoZlibCompress([]byte(s))
	fmt.Println(len(zip))
	fmt.Println(len(s))
	fmt.Println(string(DoZlibUnCompress(zip)))
}

func TestDoGzipCompress(t *testing.T) {
	s := "hello, gasdgdsahgfjhkl"
	zip := DoGZipCompress([]byte(s))
	fmt.Println(len(zip))
	fmt.Println(len(s))
	fmt.Println(string(DoGZipUnCompress(zip)))
}
