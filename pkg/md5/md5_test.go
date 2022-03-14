package md5

import (
	"testing"
)

func TestStringMd5(t *testing.T) {
	md5 := StringMd5("12345")
	t.Log(md5)

}
