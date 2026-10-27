package htree

import (
	"slices"
)

type Node struct {
	Id     int
	Symbol byte
	Count  int
	Left   *Node
	Right  *Node
}

var nodesCount int

func NewNode(symbol byte, count int, left, right *Node) *Node {
	nodesCount++

	return &Node{
		Id:     nodesCount,
		Symbol: symbol,
		Count:  count,
		Left:   left,
		Right:  right,
	}
}

// build a node array from a map of [bytes]Count
func BuildSortedNodesArray(m map[byte]int) []*Node {
	nodeArray := make([]*Node, len(m))

	i := 0
	for k, v := range m {
		nodeArray[i] = NewNode(k, v, nil, nil)
		i++
	}

	slices.SortFunc(nodeArray, func(a, b *Node) int {
		return a.Count - b.Count
	})

	return nodeArray
}

// Lookup for the first non nil element of the array
func FirstNonNilElement(nodeArray []*Node) int {
	end := len(nodeArray) - 1
	if end < 0 {
		return 0
	}
	c := 0
	for nodeArray[c] == nil && c < end {
		c++
	}
	return c
}

// look up a sorted array and merge nodes with
// same Count value
func JoinNodesWithSameValue(nodeArray []*Node) {
	end := len(nodeArray) - 1
	if end < 0 {
		return
	}

	joined := true
	for joined {
		c := FirstNonNilElement(nodeArray)
		joined = false
		for c < end {
			// if c<0 there is c and a next item n
			n := c + 1
			if nodeArray[n].Count == nodeArray[c].Count {
				JoinNodes(nodeArray, c, n)
				joined = true
			}
			c = n
		}
	}
}

// Join two consecutive nodes on a new node,
// remove old items, place the new node on the array
// remove gaps and reorder
func JoinNodes(nodeArray []*Node, i, j int) {
	end := len(nodeArray) - 1
	if end < 0 {
		return
	}

	if i == j {
		return
	}

	// make sure i is the first (outer)
	if i > j {
		t := i
		i = j
		j = t
	}

	newNode := NewNode(0, nodeArray[i].Count+nodeArray[j].Count, nodeArray[i], nodeArray[j])
	// the higher position receives the new node
	nodeArray[j] = newNode
	// the lower position is cleaned
	nodeArray[i] = nil

	// maybe we created a nil gap like [nil, 1, nil, 2, 3, 4]
	// need to bring previous items to front
	for i > 0 && nodeArray[i-1] != nil {
		nodeArray[i] = nodeArray[i-1]
		nodeArray[i-1] = nil
		i--
	}

	// now need to re-position the new node in the proper order
	q := j + 1
	for j < end && nodeArray[j].Count > nodeArray[q].Count {
		// need swap
		temp := nodeArray[j]
		nodeArray[j] = nodeArray[q]
		nodeArray[q] = temp
		j++
		q++
	}
}

// JoinNext (greedy Huffman): minimizes average code length
// Always merges the two globally smallest nodes. For a natural
// text like oz.txt with very skewed frequencies (e.g., 'e'
// appears ~70,000 times, while 'X' or 'Z' appear maybe a dozen),
// rare symbols get pushed deep into the tree — giving them very
// long codes (18 bits). But common symbols get short codes, which
// minimizes total compressed size.
func JoinNext(nodeArray []*Node) {
	end := len(nodeArray) - 1
	if end < 0 {
		return
	}

	// greedy approach instead of JoinPairs
	c := FirstNonNilElement(nodeArray)
	for c != end {
		JoinNodes(nodeArray, c, c+1)
		c++
	}
}

// JoinPairs (tournament bracket): minimizes max code depth
// Each pass pairs up ALL consecutive elements. With ~75 distinct
// symbols in oz.txt, you need ~log₂(75) ≈ 7 passes, which caps the
// tree depth around 14. It's not optimal for compression (common
// symbols may get unnecessarily long codes), but it bounds the maximum
// code length.
func JoinPairs(nodeArray []*Node) {
	end := len(nodeArray) - 1
	if end < 0 {
		return
	}

	c := FirstNonNilElement(nodeArray)
	for c != end {
		for c < end {
			// and they form a node
			JoinNodes(nodeArray, c, c+1)
			c = c + 2
		}

		// find first non null item
		c = 0
		for nodeArray[c] == nil {
			c++
		}
	}
}
