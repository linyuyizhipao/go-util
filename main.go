package main

import (
	"crypto/sha1"
	"fmt"
)

func main() {
	aa :=SHA1(SHA1([]byte("aadfsdfs奋斗奋斗dgfg")))
	fmt.Println(len(aa))
}


func SHA1(data []byte) []byte {
	h := sha1.New()
	h.Write(data)
	return h.Sum(nil)
}