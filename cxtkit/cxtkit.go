package cxtkit

import (
	"context"
	"fmt"
)

var contextList = []*contextData{}

type contextData struct {
	ctx    context.Context
	cancel func()
}

func GetCtx(num int) context.Context {
	fmt.Println(num)
	if num < 0 || num > len(contextList) {
		panic("num out of range")
	}
	return contextList[num].ctx
}

func InitContext(txNum int) func() {
	ct := context.Background()
	i := 0
	for i < txNum {
		i++
		Ctx, Cancel := context.WithCancel(ct)
		contextList = append(contextList, &contextData{
			ctx:    Ctx,
			cancel: Cancel,
		})
	}
	return cancel
}
func cancel() {
	for _, k := range contextList {
		k.cancel()
	}
}
