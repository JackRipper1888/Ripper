package cxtkit

import (
	"fmt"
	"testing"
	"time"
)

func demo(num int) {
	for {
		select {
		case <-time.After(1 * time.Second):
		case <-GetCtx(num).Done():
			fmt.Println("demo stop ", num)
			return
		}
	}
}
func TestGetCtx(t *testing.T) {
	cancel := InitContext(3)
	go demo(0)
	go demo(1)
	go demo(2)
	time.Sleep(10 * time.Second)
	cancel()
	time.Sleep(3 * time.Second)
	fmt.Println("stop")
}
