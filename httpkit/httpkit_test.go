package httpkit

import (
	"sync"
	"testing"
)

func TestGet(t *testing.T) {
	wg := sync.WaitGroup{}
	wg.Add(2)
	go func() {
		defer wg.Done()
		for i := 0; i < 10000; i++ {
			url := "http://192.168.100.200:3030/peer/select/tracker"
			Get(url, "getTracker", map[string]string{
				"appid": "250",
			})
		}
	}()
	go func() {
		defer wg.Done()
		for i := 0; i < 10000; i++ {
			url := "http://192.168.100.200:3030/peer/select/tracker"
			Get(url, "getTracker", map[string]string{
				"appid": "251",
			})
		}
	}()
	wg.Wait()

}
