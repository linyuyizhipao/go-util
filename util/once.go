package util

import (
	"fmt"
	"sync"
	"sync/atomic"
	"time"
)

var icons map[string]string
var once = &sync.Once{}

func loadIcons() {
	icons = map[string]string{
		"spades.png":   "spades111.png",
		"hearts.png":   "spades2222.png",
		"diamonds.png": "spades333.png",
		"clubs.png":    "spades6666.png",
	}
}

// NOTE: not concurrency-safe!
func Icon(name string) string {
	if icons == nil {
		icons = map[string]string{}
		time.Sleep(time.Nanosecond * 10)
		loadIcons()
	}
	return icons[name]
}

func CompareAdd(){
	var a int32 = 12
	var b int32 = 1
	atomic.AddInt32(&a,b)
	fmt.Println(a,1)

	atomic.SwapInt32(&a,19)
	fmt.Println(a,2)
}