package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/gjantsch/huffman/pkg/file"
	"github.com/gjantsch/huffman/pkg/htree"
)

func main() {
	var input = flag.String("i", "", "Input file.")
	var output = flag.String("o", "", "Output file.")
	var decompress = flag.String("d", "", "File to decompress.")
	//var verbose = flag.Bool("v", true, "Verbose mode.")
	flag.Parse()

	inputFile := *input
	outpuFile := *output
	decompressFile := *decompress
	verboseMode := false

	if outpuFile == "" {
		outpuFile = outpuFile + ".uf"
	}

	if inputFile != "" {
		Compress(inputFile, outpuFile, verboseMode)
	}

	if decompressFile != "" {
		Decompress(decompressFile)
	}
}

func Compress(inputFile, outpuFile string, verboseMode bool) {
	fmt.Printf("Reading file %s\n", inputFile)
	// totalSymbols, symbols, content, err := file.Read(inputFile)
	parsedFile, err := file.Read(inputFile)
	if err != nil {
		fmt.Printf("error reading source file: %v\n", err)
		return
	}
	fmt.Printf("Done reading source file.\n")
	fmt.Println("Building symbol tree")
	nodes := htree.BuildSortedNodesArray(parsedFile.ByteCount)
	htree.JoinNodesWithSameValue(nodes)
	htree.JoinNext(nodes)

	if len(nodes) == 0 {
		fmt.Printf("Content generated no codes.\n")
		return
	}

	codesMap := htree.BuildCodes(nodes[len(nodes)-1])

	if verboseMode {
		fmt.Printf("Content size: %d\n", parsedFile.FileSize)
		fmt.Printf("Number of symbols: %d\n", len(parsedFile.ByteCount))

		// dump tree
		htree.Render(nodes)

		// dump codes
		codesArray := htree.SortCodes(codesMap)
		maxlen := 0
		for _, v := range codesArray {
			codesLength := len(v.Code)
			if codesLength > maxlen {
				maxlen = codesLength
			}
			fmt.Printf("%8q %s\n", v.Symbol, v.Code)
		}
		fmt.Printf("Max code length: %d\n", maxlen)
	}

	fmt.Println("Creating compressed file.")
	f, err := os.OpenFile(outpuFile, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Printf("failed to open file: %v", err)
		return
	}
	defer f.Close()

	err = file.WriteHeader(f, codesMap)
	if err != nil {
		fmt.Printf("failed to write header: %v", err)
		return
	}

	err = file.WriteContent(f, parsedFile.Content, codesMap)
	if err != nil {
		fmt.Println("Error writting content", err)
	}
	f.Close()
}

func Decompress(inputFile string) {
	f, err := os.OpenFile(inputFile, os.O_RDONLY, 0644)
	if err != nil {
		log.Fatalf("failed to open compressed file: %v", err)
	}
	defer f.Close()

	readHeader, err := file.ReadHeader(f)
	if err != nil {
		fmt.Printf("failed to read header: %v", err)
	}

	decoded, err := file.ReadContent(f, readHeader)
	if err != nil {
		fmt.Printf("fail to decode content: %v", err)
	}

	fmt.Println(decoded)
	f.Close()
}
