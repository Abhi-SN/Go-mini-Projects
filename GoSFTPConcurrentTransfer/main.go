package main

import (
	"fmt"
	"sync"

	"GoSFTP/prompts"       // Import the prompts package
	"GoSFTP/serverhandler" // Import the serverhandler package
)

func main() {
	fmt.Println("Starting SFTP Transfer Application (Real SFTP Operations)...")

	// Get server details from the prompts package, interactively from the user
	serverList := prompts.GetServerInputs()

	// If no server details are provided, exit
	if len(serverList) == 0 {
		fmt.Println("No server details provided. Exiting application.")
		return
	}

	// Use a WaitGroup to wait for all goroutines to complete
	var wg sync.WaitGroup

	// Loop over the server list and spawn a goroutine for each server
	for _, server := range serverList {
		wg.Add(1) // Increment the counter for each goroutine
		// Pass the server details and the WaitGroup to the handleServer function
		// from the serverhandler package.
		go serverhandler.HandleServer(server, &wg)
	}

	// Wait for all goroutines to finish
	wg.Wait()

	fmt.Println("All SFTP transfers complete.")
	fmt.Println("Application Shutting Down.")
}
