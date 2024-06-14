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
	talliedDir, err := TallyDirectory(root)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	BuildTable(talliedDir)
}

