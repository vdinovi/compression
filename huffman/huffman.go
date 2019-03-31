package huffman

import (
	"container/heap"
	"encoding/json"
	"fmt"
	"strings"
)

type Symbol struct {
	glyph    byte
	count    uint32
	index    int // for heap
	encoding byte
	numBits  uint8
}
type SymbolPQ []*Symbol

func (pq SymbolPQ) Len() int           { return len(pq) }
func (pq SymbolPQ) Less(i, j int) bool { return pq[i].count > pq[j].count }
func (pq SymbolPQ) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
	pq[i].index = i
	pq[j].index = j
}
func (pq *SymbolPQ) Push(x interface{}) {
	n := len(*pq)
	symbol := x.(*Symbol)
	symbol.index = n
	*pq = append(*pq, symbol)
}
func (pq *SymbolPQ) Pop() interface{} {
	old := *pq
	n := len(old)
	symbol := old[n-1]
	symbol.index = -1 // for safety
	*pq = old[0 : n-1]
	return symbol
}
func (pq *SymbolPQ) update(symbol *Symbol, glyph byte, count uint32) {
	symbol.glyph = glyph
	symbol.count = count
	heap.Fix(pq, symbol.index)
}
func (pq *SymbolPQ) String() string {
	pqCopy := make(SymbolPQ, (*pq).Len())
	copy(pqCopy, (*pq))
	strs := make([]string, pqCopy.Len())
	index := 0
	for pqCopy.Len() > 0 {
		sym := heap.Pop(&pqCopy)
		strs[index] = fmt.Sprintf("(%T) %+v", sym, sym)
		index++
	}
	return strings.Join(strs, "\n")
}

type Node struct {
	symbol *Symbol
	parent *Node
	left   *Node
	right  *Node
}

func (node *Node) addLeft(symbol *Symbol) *Node {
	(*node).left = &Node{symbol: symbol, parent: node, left: nil, right: nil}
	return node.left
}
func (node *Node) addRight(symbol *Symbol) *Node {
	(*node).right = &Node{symbol: symbol, parent: node, left: nil, right: nil}
	return node.right
}
func (node *Node) subtreeString(accum *[]string, indent int) *[]string {
	sym := node.symbol
	*accum = append((*accum), fmt.Sprintf("%s%+v", strings.Repeat(" ", indent), sym))
	if node.left != nil {
		node.left.subtreeString(accum, indent+2)
	}
	if node.right != nil {
		node.right.subtreeString(accum, indent+2)
	}
	return accum
}
func (node *Node) String() string {
	accum := make([]string, 0)
	node.subtreeString(&accum, 0)
	return strings.Join(accum, "\n")
}

func byteMapString(dict *map[byte]uint32) string {
	jsonDict, err := json.MarshalIndent(*dict, "", "  ")
	if err != nil {
		fmt.Println("error:", err)
	}
	return string(jsonDict)
}

func encodingString(dict *map[byte]*StreamToken) string {
	strMap := make(map[byte]string)
	for k, v := range *dict {
		strMap[k] = v.String()
	}
	jsonDict, err := json.MarshalIndent(strMap, "", "  ")
	if err != nil {
		fmt.Println("error:", err)
	}
	return string(jsonDict)
}

func makeByteMap(data []byte) *map[byte]uint32 {
	byteMap := make(map[byte]uint32)
	for _, b := range data {
		if val, ok := byteMap[b]; ok {
			byteMap[b] = val + 1
		} else {
			byteMap[b] = 0
		}
	}
	return &byteMap
}

func makeSymbolPQ(byteMap *map[byte]uint32) *SymbolPQ {
	pq := make(SymbolPQ, len(*byteMap))
	index := 0
	for glyph, count := range *byteMap {
		pq[index] = &Symbol{glyph: glyph, count: count, index: index, encoding: 0, numBits: 0}
		index++
	}
	heap.Init(&pq)
	return &pq
}

func makeSymbolTree(pq *SymbolPQ) *Node {
	queue := make([]*Node, 1) // makeshift queue
	rootSym := &Symbol{glyph: 0, count: 0, index: 0, encoding: 0, numBits: 0}
	root := &Node{symbol: rootSym, left: nil, right: nil}
	queue[0] = root

	for len(queue) > 0 && pq.Len() > 0 {
		node := queue[0]
		queue = queue[1:] // deque
		if pq.Len() > 0 {
			sym := heap.Pop(pq).(*Symbol)
			queue = append(queue, node.addLeft(sym))
		}
		if pq.Len() > 0 {
			sym := heap.Pop(pq).(*Symbol)
			queue = append(queue, node.addRight(sym))
		}
	}
	return root
}

type StreamToken struct {
	numBits byte
	data    byte
}

func (token StreamToken) String() string {
	return fmt.Sprintf("(%d, %d)", token.numBits, token.data)
}

func makeEncoding(root *Node) (byte, *map[byte]*StreamToken) {
	queue := make([]*Node, 1) // makeshift queue
	queue[0] = root
	encoding := make(map[byte]*StreamToken)
	var maxBits uint8
	maxBits = 0

	for len(queue) > 0 {
		node := queue[0]
		queue = queue[1:]

		if node.left != nil {
			node.left.symbol.encoding = node.symbol.encoding
			node.left.symbol.numBits = node.symbol.numBits + 1
			if node.left.symbol.numBits > byte(maxBits) {
				maxBits = uint8(node.left.symbol.numBits)
			}
			encoding[node.left.symbol.glyph] = &StreamToken{numBits: node.left.symbol.numBits, data: node.left.symbol.encoding}
			queue = append(queue, node.left)
		}
		if node.right != nil {
			node.right.symbol.encoding = node.symbol.encoding | (0x01 << node.symbol.numBits)
			node.right.symbol.numBits = node.symbol.numBits + 1
			if node.right.symbol.numBits > byte(maxBits) {
				maxBits = uint8(node.right.symbol.numBits)
			}
			encoding[node.right.symbol.glyph] = &StreamToken{numBits: node.right.symbol.numBits, data: node.right.symbol.encoding}
			queue = append(queue, node.right)
		}
	}
	return maxBits, &encoding
}

func Encode(data []byte) []byte {
	byteMap := makeByteMap(data)
	//fmt.Println(byteMapString(byteMap))

	symbolPQ := makeSymbolPQ(byteMap)
	//fmt.Println(symbolPQ.String())

	symbolTree := makeSymbolTree(symbolPQ)
	//fmt.Printf("%+v\n", symbolTree)

	maxBits, encoding := makeEncoding(symbolTree)
	fmt.Println(maxBits)
	fmt.Println(encodingString(encoding))

	// TODO create packed payload of encoding
	// format should be
	// <uint32 maxBits><numBits,charEncoding,...>
	// where numBits < maxBits

	return data
}

func Decode(data []byte) []byte {
	return data
}
