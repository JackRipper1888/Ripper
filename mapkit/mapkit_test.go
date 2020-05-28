package mapkit

import "testing"

func TestConcurrentMap(t *testing.T) {
	t.Run("", func(t *testing.T) {

	})
}

//go:generata
func BenchmarkConcurrentMap_Count(b *testing.B) {
	b.StartTimer()
	for i := 0; i < b.N; i++ {

	}
	b.StopTimer()
}
