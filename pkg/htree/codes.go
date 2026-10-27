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

// Get a random length string of 0s and 1s and convert to an array
// of bytes.
// Return the array of bytes and the number of 0 bits added to the
// end to complete a full 8 byte if needed.
func BinaryStringToByteArray(binaryAsString string) ([]byte, int, error) {
	// find how many bytes will be generated and create
	// the array
	inputLength := len(binaryAsString)
	bytesToGenerate := inputLength / 8
	if inputLength%8 > 0 {
		bytesToGenerate++
	}
	generatedBytes := make([]byte, bytesToGenerate)
	padding := 0

	// fill the array of bytes
	for i := 0; i < bytesToGenerate; i++ {
		start := i * 8
		end := start + 8
		if end > inputLength {
			end = inputLength
		}
		eightBits := binaryAsString[start:end]
		if len(eightBits) < 8 {
			// pad right with zeros
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

func RealCodeLength(codeLength int) int {
	// code length
	if codeLength == 0 {
		return 0
	}

	// we need the real final length that is a multiple of 8
	remainder := (codeLength - 1) % 8
	return codeLength + 8 - remainder - 1
}

func QtdBytesFromCodeLength(codeLength int) int {
	return codeLength / 8
}
func (c *Code) ToBytes() []byte {
	codeLength := len(c.Code)
	realLength := RealCodeLength(codeLength)

	// how many bytes needed to store?
	bytes := QtdBytesFromCodeLength(realLength)
	data := make([]byte, bytes)

	// need to pad the codes with zeros
	fullCode := c.Code + strings.Repeat("0", realLength-codeLength)

	// split the binary chunks
	for i := 0; i < bytes; i++ {
		start := i * 8
		end := start + 8
		chunk := fullCode[start:end]
		val, _ := strconv.ParseUint(chunk, 2, 8)
		data[i] = byte(val)
	}
	return data
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
