package prompts

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

// ServerDetails struct holds information for connecting to a server.
// It is defined here as it's the data structure populated by this package.
type ServerDetails struct {
	ID         int    // Unique identifier for the server
	Host       string // Server hostname or IP address
	Port       int    // SSH/SFTP port
	Username   string // Username for SSH/SFTP
	Password   string // Password for SSH/SFTP (for simplicity, in real app use keys)
	Mode       string // "upload" or "download"
	LocalPath  string // Path to local file/directory
	RemotePath string // Path to remote file/directory
}

// GetServerInputs prompts the user for server configuration details interactively.
// It continues to prompt until the user decides to stop adding more servers.
func GetServerInputs() []ServerDetails {
	var serverList []ServerDetails
	reader := bufio.NewReader(os.Stdin) // Create a new reader for standard input
	serverID := 1                       // Initialize server ID counter

	fmt.Println("\n--- Enter Server Details ---")
	fmt.Println("Follow the prompts to configure each SFTP server.")
	fmt.Println("Type 'no' when asked to add another server to finish.")

	for { // Loop indefinitely until the user chooses to stop
		fmt.Printf("\n--- Server #%d Configuration ---\n", serverID)

		// Prompt for Host
		host := readInput(reader, "Enter Host (e.g., sftp.example.com): ")
		if host == "" {
			fmt.Println("Host cannot be empty. Please provide a valid host.")
			continue // Restart loop for current server if input is invalid
		}

		// Prompt for Port
		port := 0
		for {
			portStr := readInput(reader, "Enter Port (e.g., 22): ")
			var err error
			port, err = strconv.Atoi(portStr)
			if err != nil || port <= 0 {
				fmt.Println("Invalid port. Please enter a positive integer.")
			} else {
				break // Valid port entered, exit inner loop
			}
		}

		// Prompt for Username
		username := readInput(reader, "Enter Username: ")
		if username == "" {
			fmt.Println("Username cannot be empty. Please provide a valid username.")
			continue
		}

		// Prompt for Password (Note: For production, use SSH keys, not passwords)
		password := readInput(reader, "Enter Password (leave blank if using SSH key): ")
		// Password can be empty if SSH key is used, so no validation here

		// Prompt for Mode (Upload/Download)
		mode := ""
		for {
			mode = readInput(reader, "Choose Mode (upload/download): ")
			mode = strings.ToLower(mode) // Convert to lowercase for case-insensitive check
			if mode == "upload" || mode == "download" {
				break // Valid mode entered, exit inner loop
			} else {
				fmt.Println("Invalid mode. Please enter 'upload' or 'download'.")
			}
		}

		// Prompt for Local Path
		localPath := readInput(reader, "Enter Local Path (e.g., /path/to/local/file.txt): ")
		if localPath == "" {
			fmt.Println("Local path cannot be empty. Please provide a valid path.")
			continue
		}

		// Prompt for Remote Path
		remotePath := readInput(reader, "Enter Remote Path (e.g., /path/to/remote/file.txt): ")
		if remotePath == "" {
			fmt.Println("Remote path cannot be empty. Please provide a valid path.")
			continue
		}

		// Add the collected server details to the list
		serverList = append(serverList, ServerDetails{
			ID:         serverID,
			Host:       host,
			Port:       port,
			Username:   username,
			Password:   password,
			Mode:       mode,
			LocalPath:  localPath,
			RemotePath: remotePath,
		})

		serverID++ // Increment server ID for the next potential server

		// Ask if the user wants to add another server
		fmt.Print("Add another server? (yes/no): ")
		addAnother, _ := reader.ReadString('\n')
		addAnother = strings.ToLower(strings.TrimSpace(addAnother))
		if addAnother != "yes" && addAnother != "y" {
			break // User chose not to add more servers, exit the main loop
		}
	}

	fmt.Println("\n--- Finished collecting server inputs ---")
	return serverList
}

// readInput is a helper function to read a line of input from the user and trim whitespace.
func readInput(reader *bufio.Reader, prompt string) string {
	fmt.Print(prompt)
	input, _ := reader.ReadString('\n')
	return strings.TrimSpace(input)
}
