package encrypt_test

import (
	"fmt"
	"testing"
	"util/util/encrypt"
)

func TestOpenApiEncrypt_Encrypt(t *testing.T) {
	content :="ghdfgdfsdfsdffgdfggsdfgsfgsfg"
	e :=encrypt.OpenApiEncrypt{}
	encryptContent,err :=e.Encrypt(content)
	fmt.Println(encryptContent)
	if err!=nil{
		t.Log("err",err.Error())
		return
	}
	decryptContent,err :=e.Decrypt(encryptContent)
	if err!=nil{
		t.Log("err",err.Error())
		return
	}
	if content != decryptContent {
		t.Error("err","加解密不可逆")
		return
	}
	t.Log("success")
}