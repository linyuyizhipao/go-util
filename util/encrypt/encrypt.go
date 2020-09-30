package encrypt

import (
	"crypto/hmac"
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


func Signature(key,data string)(str string)  {
	mac := hmac.New(sha256.New, []byte(key))
	mac.Write([]byte(data))
	str = string(mac.Sum(nil))
	return
}