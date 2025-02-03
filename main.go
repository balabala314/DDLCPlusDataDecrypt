package main

import (
	"fmt"
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

var fileOutput *os.File

func xorString(data []byte, key byte) []byte {
	result := make([]byte, len(data))
	for i := range data {
		result[i] = data[i] ^ key
	}
	return result
}

func main() {
	var file_info os.FileInfo
	var file_size int64
	var chunk []byte
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
	_, err := os.Stat(outputFile)
	if err == nil {
		var ans string
		fmt.Printf("File %s already exist, run anyway?(y,n)", outputFile)
		fmt.Scan(&ans)
		if ans == "y" || ans == "Y" {
			err = os.Remove(outputFile)
			if err != nil {
				fmt.Println("Failed to operate!")
				os.Exit(1)
			}
		} else {
			os.Exit(1)
		}
	}
	file, err := os.Open(inputFile)
	file_info, err = file.Stat()
	file_size = file_info.Size()

	if err != nil {
		fmt.Printf(fileOpenErr, err)
		os.Exit(1)
	}
	fmt.Printf("- Version %s\n", version)
	fmt.Printf("- Loading encrypted file %s...\n", inputFile)
	fmt.Printf("- File size: %d Bytes.\n", file_size)
	fmt.Printf("- Cache size: %d Bytes.\n", cacheCapacity)

	fileOutput, err = os.Create(outputFile)
	if err != nil {
		fmt.Printf("Failed to create file %s\n", outputFile)
		os.Exit(1)
	}
	defer fileOutput.Close()

	chunk = make([]byte, cacheCapacity)
	var i int64
	for i = 0; i < file_size; i += cacheCapacity {
		end := i + cacheCapacity
		if end > file_size {
			end = file_size
			chunk = make([]byte, end-i)
		}
		file.ReadAt(chunk, int64(i))
		decryptedChunk := xorString(chunk, key)
		_, err = fileOutput.Write(decryptedChunk)
		fmt.Println("- Decrypting cached data...")
	}
	fmt.Printf(decryptSuccess, outputFile)
}
