package main

import (
	"compression/huffman"
	"fmt"
)

func main() {
	uncompressed := []byte("0123456789024680")
	compressed := huffman.Encode(uncompressed)
	decoded := huffman.Decode(compressed)

	fmt.Printf("uncompressed = %s\n", string(uncompressed))
	fmt.Printf("compressed = %s\n", string(compressed))
	fmt.Printf("decoded = %s\n", string(decoded))
}
