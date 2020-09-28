package encrypt

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/sha1"
	"encoding/base64"
	"errors"
	"fmt"
	"math/rand"
	"time"
)

const (
	secretKey = "fsdgdfhfgfggweas"   //给第三方客户的密钥，不同的第三方用户我们可以根据业务去为他分配一个唯一key
	either = 16  // golang目前提供 AES-128
	encryptLength = 128
)

type OpenApiEncrypt struct {

}


func (o *OpenApiEncrypt)Encrypt(content string)(encodeStr string,err error){

	if content == "" {
		err = errors.New("content不能为空")
		return
	}

	var key string
	text := []byte(content)
	key= o.AESSHA1PRNG(secretKey, encryptLength)
	iv := o.getIv()
	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		return
	}
	blockSize := block.BlockSize()
	originData := o.pad(text, blockSize)
	blockMode := cipher.NewCBCEncrypter(block, iv)
	crypted := make([]byte, len(originData))
	blockMode.CryptBlocks(crypted, originData)
	encodeStr = string(iv) + base64.StdEncoding.EncodeToString(crypted)

	return
}

//用于邮箱手机号码解密
func (o *OpenApiEncrypt) Decrypt(text string) (content string, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = errors.New("content if r := recover(); r != nil { Decrypt.解密失败")
		}
	}()

	if len(text) < either {
		content = text
		err = errors.New("Decrypt len(text) < either 解密字符串长度不合法")
		return
	}
	iv := text[0:either]
	encryptStr := text[either:]
	key:= o.AESSHA1PRNG(secretKey, encryptLength)
	decodeData, err := base64.StdEncoding.DecodeString(encryptStr)
	if err != nil {
		err = errors.New(fmt.Sprintf("decodeData%s",err.Error()))
		return
	}
	block, _ := aes.NewCipher([]byte(key))
	blockMode := cipher.NewCBCDecrypter(block, []byte(iv))
	originData := make([]byte, len(decodeData))
	blockMode.CryptBlocks(originData, decodeData)
	return string(o.unpad(originData)), nil
}


// secretKey 长度可能不为16，这个函数是一个映射16位操作的行为
func (o *OpenApiEncrypt) AESSHA1PRNG(content string, encryptLength int) (encryptStr string) {
	keyBytes := []byte(content)
	hashs := o.SHA1(o.SHA1(keyBytes))
	maxLen := len(hashs)
	realLen := encryptLength / 8
	if realLen > maxLen {
		lenN := realLen-maxLen
		for i:=0;i<lenN;i++{
			hashs = append(hashs,'a')
		}
	}
	encryptStr = string(hashs[0:realLen])
	return
}

//获取iv随机变量
func(o *OpenApiEncrypt) getIv() []byte {
	str := "0123456789abcdefghijklmnopqrstuvwxyz"
	bytes := []byte(str)
	result := []byte{}
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < either; i++ {
		result = append(result, bytes[r.Intn(len(bytes))])
	}
	return result
}

//加密内容长度不够自动补全
func(o *OpenApiEncrypt) pad(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}

func (o *OpenApiEncrypt)unpad(ciphertext []byte) []byte {
	length := len(ciphertext)
	//去掉最后一次的padding
	unpadding := int(ciphertext[length-1])
	return ciphertext[:(length - unpadding)]
}

func(o *OpenApiEncrypt) SHA1(data []byte) []byte {
	h := sha1.New()
	h.Write(data)
	return h.Sum(nil)
}