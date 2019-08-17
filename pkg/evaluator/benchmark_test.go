package evaluator

import (
	"testing"
)

func BenchmarkHello(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, _ = ParseToType(10, NumberType)
	}
}
