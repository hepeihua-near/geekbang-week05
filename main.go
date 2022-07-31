package main

import (
	"fmt"
	"geekbang/week05/counter"
	"sync/atomic"
	"time"
)

var (
	permit int64
	reject int64
)

func main() {
	counter := counter.NewWindowSliderCounter(10, 100, 20)
	go test01(counter)
	time.Sleep(10 * time.Second)
	fmt.Printf("permit: %v reject: %v", permit, reject)
}

func test01(count *counter.WindowSliderCounter) {
	for {
		if count.Check() {
			atomic.AddInt64(&permit, 1)
			fmt.Printf("%v permit \n", time.Now().UnixNano())
		} else {
			atomic.AddInt64(&reject, 1)
			fmt.Printf("%v reject \n", time.Now().UnixNano())
		}
		time.Sleep(48 * time.Millisecond)
	}
}
