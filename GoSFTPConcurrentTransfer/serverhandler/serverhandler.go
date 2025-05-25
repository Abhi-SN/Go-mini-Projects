package serverhandler

import (
	"fmt"
	"sync"

	"GoSFTP/prompts"    // Import prompts to use ServerDetails type
	"GoSFTP/sftpclient" // Import the sftpclient package for SFTP operations
)

// HandleServer processes a single server's SFTP operations.
// It takes a ServerDetails struct (from the prompts package) and a WaitGroup
// to signal its completion to the main goroutine.
// This function is designed to be run as a concurrent goroutine.
func HandleServer(server prompts.ServerDetails, wg *sync.WaitGroup) {
	// defer wg.Done() ensures that wg.Done() is called when this goroutine
	// finishes, regardless of whether it completes successfully or encounters an error.
	defer wg.Done()

	fmt.Printf("\n--- Starting processing for Server ID: %d (Host: %s) ---\n", server.ID, server.Host)

	// Step 1: Establish SSH/SFTP session
	// Call the Connect function from the sftpclient package.
	session, err := sftpclient.Connect(server)
	if err != nil {
		fmt.Printf("[Server %d] Error connecting: %v\n", server.ID, err)
		return // Exit this goroutine if connection fails
	}
	// Ensure the SFTP and SSH connections are closed when this function exits.
	defer session.Close()

	// Step 2: Perform upload or download based on the server's configured mode
	switch server.Mode {
	case "upload":
		// If mode is "upload", call the UploadFile function from sftpclient.
		err = sftpclient.UploadFile(session, server) // Pass the session object
		if err != nil {
			fmt.Printf("[Server %d] Upload failed: %v\n", server.ID, err)
		} else {
			fmt.Printf("[Server %d] Upload successful.\n", server.ID)
		}
	case "download":
		// If mode is "download", call the DownloadFile function from sftpclient.
		err = sftpclient.DownloadFile(session, server) // Pass the session object
		if err != nil {
			fmt.Printf("[Server %d] Download failed: %v\n", server.ID, err)
		} else {
			fmt.Printf("[Server %d] Download successful.\n", server.ID)
		}
	default:
		// Handle unexpected modes (though input validation in prompts should prevent this)
		fmt.Printf("[Server %d] Unknown mode '%s'. Skipping file transfer.\n", server.ID, server.Mode)
	}

	fmt.Printf("--- Finished processing for Server ID: %d (Host: %s) ---\n", server.ID, server.Host)
}
