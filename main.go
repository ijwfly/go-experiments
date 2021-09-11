package main

import (
	"dirSizeScanner/dirdrill"
	"fmt"
	"os"
)

func main() {
	if len(os.Args) != 2 {
		panic("usage dirSizeScanner [path]")
	}
	path := os.Args[1]

	structure := dirdrill.GetDirStructure(path)
	fmt.Printf("%d bytes\n", structure.GetSize())
	fmt.Printf("%f KB\n", float64(structure.GetSize()) / 1024)
	fmt.Printf("%f MB\n", float64(structure.GetSize()) / 1024 / 1024)
	fmt.Printf("%d files\n", structure.GetFilesCount())
}
