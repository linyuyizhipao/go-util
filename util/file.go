package util

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"sync"
)
// 判断文件是否存在
func IsExist(path string) bool {
	_, err := os.Stat(path)
	return err == nil || os.IsExist(err)
}

//往文件里面写内容
func WriteString(path string,content string)(err error){
	f ,err := os.Create(path)
	if err != nil {
		return
	}
	defer f.Close()

	_,err = f.WriteString(content)
	if err != nil {
		return
	}
	return
}
//将字节内容写入指定文件
func WriteByte(path string,byteContent []byte)(err error){
	f,err := os.Create(path)
	if err!=nil{
		return
	}
	defer f.Close()
	_,err = f.Write(byteContent)
	if err!=nil{
		return
	}
	return
}
//一行行内容的写入
func WriteLine(path string,contents []string)(err error){
	f,err := os.Create(path)
	if err!=nil{
		return
	}
	defer f.Close()

	for c := range contents {
		_,err = fmt.Fprintln(f,c)
		if err!=nil{
			return
		}
	}
	return
}
//追加写
func appendWrite(path string,content string)(err error){
	f,err := os.OpenFile(path,os.O_APPEND|os.O_WRONLY, 0644)
	if err !=nil {
		return
	}
	defer f.Close()

	_,err = fmt.Fprintln(f,content)
	return
}

//并发写入内容到文件
func MoreWrite()(err error){
	data := make(chan string)
	path := "/Applications/project/go-mongo/utils/data.txt"
	errChan := make(chan error)
	f , err :=os.Create(path)
	if err!=nil{
		return
	}
	defer f.Close()

	//并发写的内容收集到chan
	go func() {
		defer func() {
			close(data)
		}()

		wg := new(sync.WaitGroup)
		for i:=0;i<300;i++{
			wg.Add(1)

			go func() {
				data <- "并发写的内容"
				wg.Done()
			}()
		}
		wg.Wait()
	}()

	go func() {
		for content := range data {
			if _,wErr := fmt.Fprintln(f,content);wErr!=nil{
				errChan <- wErr
				return
			}
		}
		errChan <- nil
	}()

	err = <-errChan
	return
}

//读取文件内容到内存
func Read(){
	path := "/Applications/project/go-mongo/utils/data.txt"

	content,err :=ioutil.ReadFile(path)
	fmt.Println(content,"一次性读到内存中")

	f,err :=os.Open(path)
	if err!=nil{
		return
	}
	defer f.Close()

	//contentSize := make([]byte,10)
	//r := bufio.NewReader(f)
	//for{
	//	_,err = r.Read(contentSize)
	//	if err != nil {
	//		return
	//	}
	//
	//	fmt.Println(string(contentSize),"按一定byte大小读出来的内容")
	//}

	//一行行的读
	rn := bufio.NewScanner(f)
	for rn.Scan(){
		fmt.Println("as","一行行的数据内容")
	}
	if rn.Err() != nil {

	}
}


