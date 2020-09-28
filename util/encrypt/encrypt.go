package encrypt

import (
	"crypto/md5"
	"crypto/sha256"
	"fmt"
)

func Md5(str string)  (encryptStr string) {
	data := []byte(str)
	has := md5.Sum(data)
	encryptStr = fmt.Sprintf("%x", has)
	return
}

func Sha256(str string)(sha256Str string){
	strByte := []byte(str)
	has :=sha256.Sum256(strByte)
	sha256Str =fmt.Sprintf("%x",has)
	return
}