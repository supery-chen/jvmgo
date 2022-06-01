package main

import (
	"fmt"
	"time"
)

func main() {
	startedAt := time.Now()
	defer fmt.Println(time.Since(startedAt))
	defer test(startedAt)

	fmt.Printf("start:%s\n", time.Now())
	time.Sleep(time.Second)
	fmt.Printf("end:%s\n", time.Now())
}

func test(t time.Time) {
	fmt.Println("call test")
	fmt.Printf("test:%s,%s\n", t, time.Now())
	ret := time.Now().Sub(t)
	fmt.Printf("test2:%s\n", ret)
}
