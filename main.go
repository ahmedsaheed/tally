package main

import (
	"fmt"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: tally <directory>")
		return
	}

	root := os.Args[1]

	stat, err := os.Stat(root)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	if !stat.IsDir() {
		fmt.Println("Error: not a directory")
		return
	}
	talliedDir, err := TallyDirectory(root)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	BuildTable(talliedDir)
}
