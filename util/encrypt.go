package util

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/sha1"
	"encoding/base64"
	"errors"
	"math/rand"
	"time"
)

const secretKey = "fdfdsfkjhfjsdjfljalksjfkljaslkdf"
const encryptLength = 128
const either = 16

//华为要求我们返回的账号，密码得如此加密才认为我们正确  aes-cbc-128
//字符串加密规则： 1.secretKey是约定好的，但长度是不符合go的aes加密算法中的16，32，64这三种长度要求的，所以借助 AESSHA1PRNG 拿到16位，AESSHA1PRNG本身就是一个算法，java的一个函数相对于，但是其实我们可以直接去前16位
//2. 加密算法选择16的那种后，被加密的字符串长度必须是16的整数倍，不是的话 pad 函数来补齐
//3. 最后再把16位约定好的随机iv向量准备好
//4. 上正菜	encodeStr = string(iv) + base64.StdEncoding.EncodeToString(crypted)
func Encrypt(content string) (encodeStr string, err error) {
	var key string
	text := []byte(content)
	if content == "" {
		err = errors.New("content不能为空")
		return
	}
	//1.选择16加密方式，需要我们使用约定好的16位随机字符串，加解密过程都使用一样的加密字符串,AESSHA1PRNG帮助我们从长字符串中按一定规则拿到固定的16位字符串
	key, err = AESSHA1PRNG(secretKey, encryptLength)
	//获取16随机字符串，充当iv向量
	iv := getIv()
	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		return
	}
	blockSize := block.BlockSize()

	//2.待加密的字符串的长度必须是16的整倍数，不够则pad补齐
	originData := pad(text, blockSize)
	blockMode := cipher.NewCBCEncrypter(block, iv)
	crypted := make([]byte, len(originData))
	blockMode.CryptBlocks(crypted, originData)
	encodeStr = string(iv) + base64.StdEncoding.EncodeToString(crypted)
	return
}


// 模拟 Java 接口 generatorKey()
// 目前 encryptLength 仅支持 128bit 长度
// 因为 SHA1() 返回的长度固定为 20byte 160bit
// 所以 encryptLength 超过这个长度，就无法生成了
// 因为不知道 java 中 AES SHA1PRNG 的生成逻辑
//https://support.huaweicloud.com/accessg-marketplace/zh-cn_topic_0070649063.html
//对接华为云的key加密有个随机值过程，采用此方法随机
func AESSHA1PRNG(content string, encryptLength int) (encryptStr string, err error) {
	keyBytes := []byte(content)
	hashs := SHA1(SHA1(keyBytes))
	maxLen := len(hashs)
	realLen := encryptLength / 8
	if realLen > maxLen {
		err = errors.New("realLen 不能比 maxLen 还要大")
		return
	}
	encryptStr = string(hashs[0:realLen])
	return
}

func SHA1(data []byte) []byte {
	h := sha1.New()
	h.Write(data)
	return h.Sum(nil)
}

//加密内容长度不够自动补全
func pad(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}

func unpad(ciphertext []byte) []byte {
	length := len(ciphertext)
	//去掉最后一次的padding
	unpadding := int(ciphertext[length-1])
	return ciphertext[:(length - unpadding)]
}

//用于邮箱手机号码解密
func Decrypt(text string) (content string, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = errors.New("解密失败")
		}
	}()

	if len(text) < either {
		content = text
		err = errors.New("解密字符串长度不合法")
		return
	}

	iv := text[0:either]
	encryptStr := text[either:]
	key, err := AESSHA1PRNG(secretKey, encryptLength)
	if err != nil {
		return
	}
	decodeData, err := base64.StdEncoding.DecodeString(encryptStr)
	if err != nil {
		return
	}

	//生成密码数据块cipher.Block
	block, _ := aes.NewCipher([]byte(key))
	//解密模式
	blockMode := cipher.NewCBCDecrypter(block, []byte(iv))
	//输出到[]byte数组
	originData := make([]byte, len(decodeData))
	blockMode.CryptBlocks(originData, decodeData)
	//去除填充,并返回
	return string(unpad(originData)), nil
}

//获取iv随机变量
func getIv() []byte {
	str := "0123456789abcdefghijklmnopqrstuvwxyz"
	bytes := []byte(str)
	result := []byte{}
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < either; i++ {
		result = append(result, bytes[r.Intn(len(bytes))])
	}
	return result
}
