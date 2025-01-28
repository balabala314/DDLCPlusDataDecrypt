package main

import (
	"fmt"
	"io/ioutil"
	"os"
)

const (
	key            = 40
	cacheCapacity  = 100 * 1024 * 1024 // 100MB
	version        = "1.0"
	usageMessage   = "- Usage: %s <asset file *.cy> [output file]\n"
	supportedExt   = ".cy"
	fileOpenErr    = "! Unable to open file: %v\n"
	decryptSuccess = "- Done! Decrypted file is output to %s.\n"
)

func xorString(data []byte, key byte) []byte {
	result := make([]byte, len(data))
	for i := range data {
		result[i] = data[i] ^ key
	}
	return result
}

func main() {
	args := os.Args
	if len(args) < 2 {
		fmt.Printf(usageMessage, args[0])
		os.Exit(1)
	}

	inputFile := args[1]
	outputFile := ""

	if len(args) >= 3 {
		outputFile = args[2]
	} else {
		if len(inputFile) > 2 && inputFile[len(inputFile)-3:] == supportedExt {
			outputFile = inputFile[:len(inputFile)-2] + "assets"
		} else {
			fmt.Printf("! Only supports files with .cy extension.\n")
			os.Exit(1)
		}
	}

	data, err := ioutil.ReadFile(inputFile)
	if err != nil {
		fmt.Printf(fileOpenErr, err)
		os.Exit(1)
	}
	fmt.Printf("- Version %s\n", version)
	fmt.Printf("- Loading encrypted file %s...\n", inputFile)
	fmt.Printf("- File size: %d Bytes.\n", len(data))
	fmt.Printf("- Cache size: %d Bytes.\n", cacheCapacity)

	var decryptedData []byte
	for i := 0; i < len(data); i += cacheCapacity {
		end := i + cacheCapacity
		if end > len(data) {
			end = len(data)
		}
		chunk := data[i:end]
		decryptedChunk := xorString(chunk, key)
		decryptedData = append(decryptedData, decryptedChunk...)
		fmt.Println("- Decrypting cached data...")
	}

	err = ioutil.WriteFile(outputFile, decryptedData, 0644)
	if err != nil {
		fmt.Printf(fileOpenErr, err)
		os.Exit(1)
	}

	fmt.Printf(decryptSuccess, outputFile)
}



