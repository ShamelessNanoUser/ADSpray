package read

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func ReadFile(filename string) []string {
	// Initialize array
	var lines []string

	// Open file
	file, err := os.Open(filename)

	// Check if the file could be opened, else exit
	if err != nil {
		fmt.Println("Error opening file", err)
		file.Close()
		os.Exit(1)
	}

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		// If line != empty -> add to array
		if line != "" {
			lines = append(lines, line)
		}
	}

	return lines
}
