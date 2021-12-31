package str

import (
	"fmt"
	"testing"
)

func TestBytesOf(t *testing.T) {
	fmt.Println(StringOf(BytesOf("helllo")))
}

func BenchmarkName(b *testing.B) {

	for i := 0; i < b.N; i++ {

	}
}
