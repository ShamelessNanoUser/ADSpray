package logger

import (
	write "github.com/ShamelessNanoUser/ADSpray/io"
	"fmt"
	"os"
	"strings"
	"time"
)

func InitiateLogFile(filename string) {
	// Check if file exists
	_, err := os.Stat(filename)

	// If file does not exist create file
	if os.IsNotExist(err) {
		write.CreateLogFile(filename)
		return
	}

	// If file exists ask if user wants to overwrite
	var overwrite string
	fmt.Print("File ", filename, " exists, overwrite? (y/n) ")
	fmt.Scan(&overwrite)

	if strings.ToLower(overwrite) == "y" {
		write.CreateLogFile(filename)
		return
	}

	// If the value of overwrite is not y, stop execution
	fmt.Println("Exiting...")
	os.Exit(0)
}

func AddLogEntry(filename string, logLine string) {

	// Create date header for start of each log entry
	currentTime := time.Now().Format("02/01/2006 15:04:05")
	dateHeader := "[" + currentTime + "]"
	logEntry := dateHeader + " - " + logLine + "\n"

	// Open file for write operations
	file, err := os.OpenFile(filename, os.O_APPEND|os.O_WRONLY, 0644)

	if err != nil {
		fmt.Println("Failed to append to log file")
		return
	}

	// Write log entry and close file
	file.Write([]byte(logEntry))
	file.Close()
}
