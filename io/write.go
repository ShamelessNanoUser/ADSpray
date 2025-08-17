package read

import (
	"fmt"
	"os"
)

func CreateLogFile(filename string) {
	// Create file
	file, err := os.Create(filename)

	// Failed to created file
	if err != nil {
		fmt.Println("Failed to create file, check permissions of current directory")
		os.Exit(1)
	}

	// Close the files
	file.Close()
}
