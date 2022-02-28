package md5

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io"
)

func StringMd5(data string) string {
	h := md5.New()
	h.Write([]byte(data))
	return hex.EncodeToString(h.Sum(nil))
}

func StringMd5Sum(data string) string {
	w := md5.New()
	io.WriteString(w, data)
	return fmt.Sprintf("%x", w.Sum(nil))
}

func ByteMd5(data []byte) string {
	h := md5.New()
	h.Write(data)
	return hex.EncodeToString(h.Sum(nil))
}

func ByteMd5Sum(data []byte) string {
	has := md5.Sum(data)
	return fmt.Sprintf("%x", has)
}
