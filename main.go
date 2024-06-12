package main

import (
	"fmt"
	"os"
	"strings"
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

	sortByTotalLines(talliedDir)	
	fmt.Println(strings.Repeat("-", 30))
	fmt.Printf("%-20s %s\n", "Language", "Total Lines")
	fmt.Println(strings.Repeat("-", 30))
	for _, lang := range talliedDir {
		fmt.Printf("%-20s %d\n", lang.Name, lang.TotalCount)
	}
	fmt.Println(strings.Repeat("-", 30))
	fmt.Printf("%-20s %d\n", "Total", SumLines(talliedDir))
	fmt.Println(strings.Repeat("-", 30))
	
}

