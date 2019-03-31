package huffman

import (
	"fmt"
	"testing"
)

func TestCompression(t *testing.T) {
	uncompressed := []byte("0123456789024680")
	compressed := Encode(uncompressed)
	decoded := Decode(compressed)

	fmt.Printf("uncompressed = %s\n", string(uncompressed))
	fmt.Printf("compressed = %s\n", string(compressed))
	fmt.Printf("decoded = %s\n", string(decoded))

	if len(uncompressed) >= len(compressed) {
		t.Errorf("no reduction in size: %d >= %d\n", len(uncompressed), len(compressed))
	}
}
