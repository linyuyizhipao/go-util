
package main

import (
	"fmt"
	"golang.org/x/time/rate"
	"time"
)

var gameScene = rate.NewLimiter(1, 10)

func main() {

	for i:=0;i<100;i++{
		k :=i
		if isGameSceneAllow(){
			fmt.Println("我是被接受的请求",time.Now().Unix(),k)
		}
	}


	//9秒钟sleep，忽略代码执行时间，那么将会产生9个
	time.Sleep(time.Second * 9)


	//以下打印9个，则证明限流起作用了
	for i:=0;i<100;i++{
		k :=i
		if isGameSceneAllow(){
			fmt.Println("我是被接受的请求2222",time.Now().Unix(),k)
		}
	}

	time.Sleep(time.Second * 9)

}

func isGameSceneAllow()(b bool){
	if gameScene.Allow() == false {
		return
	}
	b =true
	return
}