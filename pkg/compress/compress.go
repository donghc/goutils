package compress

import (
	"bytes"
	"compress/gzip"
	"compress/zlib"
	"io"
)

//DoZlibCompress 进行zlib压缩
func DoZlibCompress(src []byte) []byte {
	var in bytes.Buffer
	w := zlib.NewWriter(&in)
	w.Write(src)
	w.Close()
	return in.Bytes()
}

//DoZlibUnCompress 进行zlib解压缩
func DoZlibUnCompress(compressSrc []byte) []byte {
	b := bytes.NewReader(compressSrc)
	var out bytes.Buffer
	r, _ := zlib.NewReader(b)
	io.Copy(&out, r)
	return out.Bytes()
}

//DoGZipCompress 进行gzip压缩
func DoGZipCompress(src []byte) []byte {
	var in bytes.Buffer
	w := gzip.NewWriter(&in)
	w.Write(src)
	w.Close()
	return in.Bytes()
}

//DoGZipUnCompress 进行gzip解压缩
func DoGZipUnCompress(compressSrc []byte) []byte {
	b := bytes.NewReader(compressSrc)
	var out bytes.Buffer
	r, _ := gzip.NewReader(b)
	io.Copy(&out, r)
	return out.Bytes()
}
