package file

import (
	"fmt"
	"os"
	"strings"

	"github.com/gjantsch/huffman/pkg/htree"
)

// Header Structure:
// - 0: byte with number of codes
// Entries:
// - 0: byte with simbol
// - 1: code len
// - n: encoded byte
func WriteHeader(f *os.File, m map[byte]htree.Code) error {

	numEntries := len(m)
	if numEntries == 0 {
		return fmt.Errorf("no entries on the header")
	}
	if numEntries > 255 {
		return fmt.Errorf("nun entries > 255")
	}

	_, err := f.Write([]byte{byte(numEntries)})
	if err != nil {
		return fmt.Errorf("error writing header number of entries: %v", err)
	}

	for _, code := range m {
		entry := append([]byte{code.Symbol, byte(len(code.Code))}, code.ToBytes()...)
		_, err := f.Write(entry)
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
	// validate num entries
	if numEntries <= 0 {
		return nil, fmt.Errorf("no entries")
	}

	// code array to be returned
	codes := make(map[string]htree.Code, numEntries)

	// read header
	for i := 0; i < int(numEntries); i++ {

		entry := []byte{0x0, 0x0}
		_, err := f.Read(entry)
		if err != nil {
			return nil, fmt.Errorf("error reading codes: %v", err)
		}

		symbol := entry[0]
		codeLength := int(entry[1])
		realCodeLen := htree.RealCodeLength(int(entry[1]))
		codeLenBytes := htree.QtdBytesFromCodeLength(realCodeLen)

		codeBytes := make([]byte, codeLenBytes)
		_, err = f.Read(codeBytes)
		if err != nil {
			return nil, fmt.Errorf("error reading codes: %v", err)
		}

		code := ""
		for _, b := range codeBytes {
			binary := fmt.Sprintf("%b", b)
			if len(binary) < 8 {
				binary = strings.Repeat("0", 8-len(binary)) + binary
			}
			code += binary
		}
		code = code[0:codeLength]
		codes[code] = htree.Code{Symbol: symbol, Code: code}
	}

	return codes, nil
}
