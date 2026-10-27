package htree

import (
	"fmt"
	"strings"
)

const (
	Spacing   = "|  "
	Connector = "+--"
)

// Render the tree with ASCII codes
// [id] <symbol> (count)
// +--L: [id] <symbol> (count)
// |  +--L: [id] <symbol> (count)
// |  |  +--L: [id] <symbol> (count)
// |  |  +--R: [id] <symbol> (count)
// |  +--R: [id] <symbol> (count)
// +--R: [id] <symbol> (count)
func Render(nodeArray []*Node) {
	end := len(nodeArray) - 1
	if end < 0 {
		return
	}
	RenderNode(nodeArray[end], 0)
}

// Render Single Node
func RenderNode(n *Node, depth int) {
	if n == nil {
		return
	}

	// root node
	if depth == 0 {
		fmt.Printf("[%03d] %q (%d)\n", n.Id, n.Symbol, n.Count)
	}

	if n.Left != nil {
		fmt.Printf("%s%sL: [%03d] %s(%d)\n", strings.Repeat(Spacing, depth), Connector, n.Left.Id, PrettySymbol(n.Left.Symbol), n.Left.Count)
		RenderNode(n.Left, depth+1)
	}

	if n.Right != nil {
		fmt.Printf("%s%sR: [%03d] %s(%d)\n", strings.Repeat(Spacing, depth), Connector, n.Right.Id, PrettySymbol(n.Right.Symbol), n.Right.Count)
		RenderNode(n.Right, depth+1)
	}
}

func PrettySymbol(s byte) string {
	if s == 0 {
		return ""
	}
	return fmt.Sprintf("%q ", s)
}
