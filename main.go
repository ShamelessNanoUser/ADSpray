package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	read "github.com/ShamelessNanoUser/ADSpray/io"
	logger "github.com/ShamelessNanoUser/ADSpray/logger"
	spray "github.com/ShamelessNanoUser/ADSpray/spray"

	"github.com/spf13/pflag"
)

func askForConfirmation(line string) {
	var confirmation string
	fmt.Print(line + " (y/n) ")
	fmt.Scan(&confirmation)

	if strings.ToLower(confirmation) == "y" {
		return
	}

	// If the value of confirmation is not y, stop execution
	fmt.Println("Exiting...")
	os.Exit(0)
}

func main() {
	// Define argument variables
	var help bool
	var domain string
	var ipAddress string
	var usernameFile string
	var passwordFile string
	var attempts int
	var timeout int
	var jitter int
	var logFile string

	// Set variables
	pflag.BoolVarP(&help, "help", "h", false, "Argument to show help")
	pflag.StringVarP(&domain, "domain", "d", "", "FQDN of target domain")
	pflag.StringVarP(&ipAddress, "dc", "t", "", "IP address of target domain controller")
	pflag.StringVarP(&usernameFile, "users", "u", "", "Path to text file containing usernames to spray against")
	pflag.StringVarP(&passwordFile, "passwords", "p", "", "Path to text file containing passwords to spray")
	pflag.IntVarP(&attempts, "attempts", "a", 0, "Amount of passwords to try during one spray cycle")
	pflag.IntVarP(&timeout, "sleep", "s", 0, "Sleep time in minutes between spray cycles")
	pflag.IntVarP(&jitter, "jitter", "j", 50, "Sleep time in milliseconds between connection attempts, also known as jitter")
	pflag.StringVarP(&logFile, "output", "o", "spray_log.txt", "Output file for spray log")

	pflag.Parse()

	// If help argument is passed, show help
	if help {
		pflag.Usage()
		os.Exit(0)
	}

	// Check if all arguments are set
	pflag.VisitAll(func(f *pflag.Flag) {
		if f.Name != "help" && f.Name != "h" {
			if f.Value.String() == "" || f.Value.String() == "0" || f.Value.String() == "false" {
				fmt.Println("-"+f.Shorthand, "--"+f.Name+" was not set\n")
				pflag.Usage()
				os.Exit(0)
			}
		}
	})

	// Ask for confirmations before starting the password spray
	askForConfirmation("Did you check if there is a Fine Grained Password Policy in place?")
	askForConfirmation("ADSpray is now gonna start spraying " + strconv.Itoa(attempts) + " passwords per " + strconv.Itoa(timeout) + " minutes, continue?")

	// Initialize log file
	logger.InitiateLogFile(logFile)
	logger.AddLogEntry(logFile, "Starting password spray against "+domain)

	// Read files to arrays
	usernameArray := read.ReadFile(usernameFile)
	passwordArray := read.ReadFile(passwordFile)

	// Start spraying
	spray.StartSMBSpray(domain, ipAddress, usernameArray, passwordArray, logFile, attempts, timeout, jitter)
}
