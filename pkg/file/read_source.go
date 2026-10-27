package file

import (
	"fmt"
	"io"
	"os"
)

type ParsedFile struct {
	FileSize  int
	ByteCount map[byte]int
	Content   []byte
}

// Read the source file and return the parsed file struct and error
func Read(filename string) (ParsedFile, error) {

	p := ParsedFile{
		FileSize:  0,
		ByteCount: make(map[byte]int),
		Content:   make([]byte, 0),
	}

	fileDescriptor, err := os.Open(filename)
	if err != nil {
		return p, fmt.Errorf("error reading file: %v", err)
	}
	defer fileDescriptor.Close()

	for {
		buffer := make([]byte, 1024)
		bytesRead, err := fileDescriptor.Read(buffer)
		if err == io.EOF {
			break
		}
		if err != nil {
			return p, fmt.Errorf("error reading input file: %v", err)
		}
		p.Content = append(p.Content, buffer[:bytesRead]...)
		for _, r := range buffer[:bytesRead] {
			p.ByteCount[r]++
			p.FileSize++
		}
	}
	return p, nil
}
