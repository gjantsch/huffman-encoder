package htree

import (
	"fmt"
	"slices"
	"strconv"
	"strings"
)

type Code struct {
	Symbol byte
	Code   string
}

// BinaryStringToByteArray converts a string of '0'/'1' characters to bytes,
// padding the last byte with trailing zeros if needed.
// Returns the bytes, the number of padding bits added, and any error.
func BinaryStringToByteArray(binaryAsString string) ([]byte, int, error) {
	inputLength := len(binaryAsString)
	bytesToGenerate := inputLength / 8
	if inputLength%8 > 0 {
		bytesToGenerate++
	}
	generatedBytes := make([]byte, bytesToGenerate)
	padding := 0

	for i := 0; i < bytesToGenerate; i++ {
		start := i * 8
		end := start + 8
		if end > inputLength {
			end = inputLength
		}
		eightBits := binaryAsString[start:end]
		if len(eightBits) < 8 {
			padding = 8 - len(eightBits)
			eightBits += strings.Repeat("0", padding)
		}
		val, err := strconv.ParseUint(eightBits, 2, 8)
		if err != nil {
			return []byte{}, 0, fmt.Errorf("error parsing binary string: %v", err)
		}
		generatedBytes[i] = byte(val)
	}
	return generatedBytes, padding, nil
}

var codes map[byte]Code

func Traverse(n *Node, code string) {
	if n == nil {
		return
	}

	if n.Left == nil && n.Right == nil {
		codes[n.Symbol] = Code{Symbol: n.Symbol, Code: code}
	}
	Traverse(n.Left, code+"0")
	Traverse(n.Right, code+"1")
}

func BuildCodes(n *Node) map[byte]Code {
	codes = make(map[byte]Code)
	Traverse(n, "")
	return codes
}

func SortCodes(m map[byte]Code) []Code {
	ordered := make([]Code, len(m))
	i := 0
	for _, v := range m {
		ordered[i] = v
		i++
	}

	slices.SortFunc(ordered, func(a, b Code) int {
		return len(a.Code) - len(b.Code)
	})
	return ordered
}

// CanonicalCodes converts arbitrary Huffman codes to canonical form.
// Canonical codes preserve code lengths but assign bit patterns
// deterministically by sorting symbols by (length, symbol value).
// This allows the header to store only lengths, not actual bit patterns.
func CanonicalCodes(m map[byte]Code) map[byte]Code {
	ordered := SortCodes(m)
	slices.SortStableFunc(ordered, func(a, b Code) int {
		if len(a.Code) != len(b.Code) {
			return len(a.Code) - len(b.Code)
		}
		return int(a.Symbol) - int(b.Symbol)
	})

	canonical := make(map[byte]Code, len(m))
	code := 0
	prevLen := 0
	for _, c := range ordered {
		curLen := len(c.Code)
		code <<= (curLen - prevLen)
		bits := fmt.Sprintf("%0*b", curLen, code)
		canonical[c.Symbol] = Code{Symbol: c.Symbol, Code: bits}
		code++
		prevLen = curLen
	}
	return canonical
}
