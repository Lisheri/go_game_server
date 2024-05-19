package utils

import (
	"bytes"
	"compress/gzip"
	"crypto/md5"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"io"
	"io/ioutil"

	"github.com/forgoer/openssl"
)

// 加密
func AesCBCEncrypt(src, key, iv []byte, padding string) ([]byte, error) {
	data, err := openssl.AesCBCEncrypt(src, key, iv, padding)
	if err != nil {
		return nil, err
	}
	return []byte(hex.EncodeToString(data)), nil
}

// 解密
func AesCBCDecrypt(src, key, iv []byte, padding string) ([]byte, error) {
	data, err := hex.DecodeString(string(src))
	if err != nil {
		return nil, err
	}
	return openssl.AesCBCDecrypt(data, key, iv, padding)
}

func Md5(text string) string {
	hashMd5 := md5.New()
	io.WriteString(hashMd5, text)
	return fmt.Sprintf("%x", hashMd5.Sum(nil))
}

func Zip(data []byte) ([]byte, error) {
	var b bytes.Buffer
	/// 9级gzip压缩
	gz, _ := gzip.NewWriterLevel(&b, 9)
	if _, err := gz.Write([]byte(data)); err != nil {
		return nil, err
	}
	if err := gz.Flush(); err != nil {
		return nil, err
	}
	if err := gz.Close(); err != nil {
		return nil, err
	}

	return b.Bytes(), nil
}

// 解压缩
func UnZip(data []byte) ([]byte, error) {
	b := new(bytes.Buffer)
	binary.Write(b, binary.LittleEndian, data)
	r, err := gzip.NewReader(b)
	if err != nil {
		return nil, err
	}
	defer r.Close()

	unzipData, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, err
	}
	// 解压缩后的文件
	return unzipData, nil
}

func Password(pwd string, pwdCode string) string {
	return Md5(pwd + pwdCode)
}
