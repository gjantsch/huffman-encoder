package file

import (
	"fmt"
	"os"
	"slices"

	"github.com/gjantsch/huffman/pkg/htree"
)

// Header structure:
// - 1 byte: number of entries
// Per entry, sorted by (code length, symbol value):
// - 1 byte: symbol
// - 1 byte: code length in bits
// Canonical codes are reconstructed from lengths on read.
func WriteHeader(f *os.File, m map[byte]htree.Code) error {
	numEntries := len(m)
	if numEntries == 0 {
		return fmt.Errorf("no entries on the header")
	}
	if numEntries > 255 {
		return fmt.Errorf("num entries > 255")
	}

	_, err := f.Write([]byte{byte(numEntries)})
	if err != nil {
		return fmt.Errorf("error writing header number of entries: %v", err)
	}

	ordered := htree.SortCodes(m)
	slices.SortStableFunc(ordered, func(a, b htree.Code) int {
		if len(a.Code) != len(b.Code) {
			return len(a.Code) - len(b.Code)
		}
		return int(a.Symbol) - int(b.Symbol)
	})

	for _, code := range ordered {
		_, err := f.Write([]byte{code.Symbol, byte(len(code.Code))})
		if err != nil {
			return fmt.Errorf("error writing header symbols: %v", err)
		}
	}
	return nil
}

func ReadHeader(f *os.File) (map[string]htree.Code, error) {
	data := make([]byte, 1)
	_, err := f.Read(data)
	if err != nil {
		return nil, fmt.Errorf("error reading header struct: %v", err)
	}

	numEntries := int(data[0])
	if numEntries <= 0 {
		return nil, fmt.Errorf("no entries")
	}

	type entry struct {
		symbol byte
		length int
	}
	entries := make([]entry, numEntries)
	for i := range entries {
		pair := []byte{0, 0}
		_, err := f.Read(pair)
		if err != nil {
			return nil, fmt.Errorf("error reading codes: %v", err)
		}
		entries[i] = entry{symbol: pair[0], length: int(pair[1])}
	}

	// Reconstruct canonical codes from sorted (length, symbol) entries.
	// Uses the same assignment algorithm as CanonicalCodes in htree.
	codes := make(map[string]htree.Code, numEntries)
	code := 0
	prevLen := 0
	for _, e := range entries {
		code <<= (e.length - prevLen)
		bits := fmt.Sprintf("%0*b", e.length, code)
		codes[bits] = htree.Code{Symbol: e.symbol, Code: bits}
		code++
		prevLen = e.length
	}

	return codes, nil
}
