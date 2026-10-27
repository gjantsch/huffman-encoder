package file

import (
	"encoding/binary"
	"fmt"
	"os"

	"github.com/gjantsch/huffman/pkg/htree"
)

func WriteContent(f *os.File, content []byte, codes map[byte]htree.Code) error {
	// save the content length is a lazy way
	// to handle the exceeding bytes at the
	// end of file, saves some logic
	contentLengthBytes := make([]byte, 8)
	binary.BigEndian.PutUint64(contentLengthBytes, uint64(len(content)))
	f.Write(contentLengthBytes)
	// write the compressed content
	buffer := ""
	for _, v := range content {
		buffer += codes[v].Code
		// every 1024 or more, write to file and empty the buffer
		if len(buffer) > 1024 && (len(buffer)%8) == 0 {
			generated, padding, err := htree.BinaryStringToByteArray(buffer)
			if err != nil {
				return fmt.Errorf("error converting: %v", err)
			}
			_, err = f.Write(generated)
			if err != nil {
				return fmt.Errorf("failed to write content: %v", err)
			}
			if padding != 0 {
				return fmt.Errorf("unexpected padding: %d", padding)
			}
			buffer = ""
		}
	}
	// write any left over
	if buffer != "" {
		generated, padding, err := htree.BinaryStringToByteArray(buffer)
		if err != nil {
			return fmt.Errorf("error converting: %v", err)
		}

		if padding > 0 {
			generated = append(generated, []byte{byte(padding)}...)
		}
		_, err = f.Write(generated)
		if err != nil {
			return fmt.Errorf("error writting: %v", err)
		}
	}
	fmt.Println()
	return nil
}

func ReadContent(f *os.File, codes map[string]htree.Code) (string, error) {

	contentLenghtBytes := make([]byte, 8)
	_, err := f.Read(contentLenghtBytes)
	if err != nil {
		return "", fmt.Errorf("failed to read content length: %v", err)
	}

	contentLength := binary.BigEndian.Uint64(contentLenghtBytes)
	content := ""
	minlen := 255
	maxlen := 0
	for k := range codes {
		l := len(k)
		if minlen > l {
			minlen = l
		}
		if maxlen < l {
			maxlen = l
		}
	}

	code := ""
	for {
		buffer := make([]byte, 1024)
		bytesRead, err := f.Read(buffer)
		if err != nil {
			break
		}
		for _, r := range buffer[:bytesRead] {
			code += fmt.Sprintf("%08b", r)
		}
	}

	s := 0
	for p := minlen; p < len(code) && len(content) < int(contentLength); p++ {
		val, ok := codes[code[s:p]]
		if ok {
			s = p
			content += string(val.Symbol)
		}
		if (p - s) > maxlen {
			return "", fmt.Errorf("corrupted content")
		}
	}

	return content, nil
}
