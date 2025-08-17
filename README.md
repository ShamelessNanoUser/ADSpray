# ADSpray

ADSpray is a password spraying that can be used to spray passwords against Active Directory. Since it's written in Go, it should work on any operating system and any architecture. Operators can configure the amount of attempts and how long the software should wait to try another spray. In addition, the timeout between authentication attempts can be changed. A logon attempt is seen as valid when the user is able to log in, the account is disabled or when the status must change. All login attempts are logged to a text file. 

## Features

- **Password Spraying:** Attempts a specified number of passwords per user, per cycle.
- **SMB Protocol:** Operates over SMB for authentication attempts.
- **Configurable Delays:** Supports jitter (milliseconds) between attempts and sleep (minutes) between cycles to avoid detection and lockouts.
- **Logging:** Outputs results to a log file for review.

## Installation

```
git clone https://github.com/ShamelessNanoUser/ADSpray
cd ADSpray
go build 
```

## Usage

```sh
./ADSpray -h
Usage of ./ADSpray:
  -a, --attempts int       Amount of passwords to try during one spray cycle
  -t, --dc string          IP address of target domain controller
  -d, --domain string      FQDN of target domain
  -h, --help               Argument to show help
  -j, --jitter int         Sleep time in milliseconds between connection attempts, also known as jitter (default 50)
  -o, --output string      Output file for spray log (default "spray_log.txt")
  -p, --passwords string   Path to text file containing passwords to spray
  -s, --sleep int          Sleep time in minutes between spray cycles
  -u, --users string       Path to text file containing usernames to spray against
```

## Example

```
./ADSpray -d howdoiexitvim.local -t 192.168.1.0 -u users.txt -p passwords.txt -a 3 -s 15 
Did you check if there is a Fine Grained Password Policy in place? (y/n) y
ADSpray is now gonna start spraying 2 passwords per 10 minutes, continue? (y/n) y
File spray_log.txt exists, overwrite? (y/n) y
[+] Spraying the password Password01
[+] Spraying the password Password01!
[!] Found valid credentials: john:Password01!
[!] Found valid credentials: lucas:Password01! - NT_STATUS_DISABLED
[!] Found valid credentials: noah:Password01! - STATUS_PASSWORD_EXPIRED
[+] Spraying the password Password123
[+] Reached maximum attempts, sleeping for: 15:00

```
