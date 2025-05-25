package sftpclient

import (
	"fmt"
	"io"
	"os"
	"time"

	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"

	"GoSFTP/prompts" // To use ServerDetails type
)

// ClientSession holds the active SFTP and SSH clients
type ClientSession struct {
	sftpClient *sftp.Client
	sshClient  *ssh.Client
	serverID   int // Store server ID for logging
}

// Connect establishes a real SSH and SFTP connection.
func Connect(details prompts.ServerDetails) (*ClientSession, error) {
	fmt.Printf("[Server %d] Attempting to connect to %s:%d...\n", details.ID, details.Host, details.Port)

	// Configure SSH client
	config := &ssh.ClientConfig{
		User: details.Username,
		Auth: []ssh.AuthMethod{
			ssh.Password(details.Password), // Use password authentication
			// For production: Consider ssh.PublicKeys(signer) for key-based authentication
		},
		// WARNING: InsecureIgnoreHostKey is highly insecure for production!
		// It bypasses host key verification, making you vulnerable to Man-in-the-Middle attacks.
		// For production, use ssh.FixedHostKey(hostKey) or ssh.KnownHosts(knownHostsFilePath)
		// to verify the server's identity.
		HostKeyCallback: ssh.InsecureIgnoreHostKey(), // DANGER: DO NOT USE IN PRODUCTION
		Timeout:         10 * time.Second,            // Connection timeout
	}

	addr := fmt.Sprintf("%s:%d", details.Host, details.Port)
	sshClient, err := ssh.Dial("tcp", addr, config)
	if err != nil {
		return nil, fmt.Errorf("failed to dial SSH server %s: %w", addr, err)
	}

	sftpClient, err := sftp.NewClient(sshClient)
	if err != nil {
		sshClient.Close() // Close SSH client if SFTP client creation fails
		return nil, fmt.Errorf("failed to create SFTP client: %w", err)
	}

	fmt.Printf("[Server %d] Successfully connected to %s:%d.\n", details.ID, details.Host, details.Port)
	return &ClientSession{
		sftpClient: sftpClient,
		sshClient:  sshClient,
		serverID:   details.ID,
	}, nil
}

// Close closes the SFTP and SSH client connections.
func (cs *ClientSession) Close() {
	if cs.sftpClient != nil {
		cs.sftpClient.Close()
	}
	if cs.sshClient != nil {
		cs.sshClient.Close()
	}
	fmt.Printf("[Server %d] Connection closed.\n", cs.serverID)
}

// UploadFile performs a real SFTP upload using the provided client session.
func UploadFile(cs *ClientSession, details prompts.ServerDetails) error {
	fmt.Printf("[Server %d] Uploading '%s' to '%s'...\n", cs.serverID, details.LocalPath, details.RemotePath)

	// Open the local file for reading
	localFile, err := os.Open(details.LocalPath)
	if err != nil {
		return fmt.Errorf("failed to open local file '%s': %w", details.LocalPath, err)
	}
	defer localFile.Close() // Ensure local file is closed

	// Create the remote file for writing
	remoteFile, err := cs.sftpClient.Create(details.RemotePath)
	if err != nil {
		return fmt.Errorf("failed to create remote file '%s': %w", details.RemotePath, err)
	}
	defer remoteFile.Close() // Ensure remote file is closed

	// Copy data from local file to remote file
	bytesWritten, err := io.Copy(remoteFile, localFile)
	if err != nil {
		return fmt.Errorf("failed to copy data to remote file: %w", err)
	}

	fmt.Printf("[Server %d] Upload of '%s' to '%s' COMPLETE (%d bytes).\n", cs.serverID, details.LocalPath, details.RemotePath, bytesWritten)
	return nil
}

// DownloadFile performs a real SFTP download using the provided client session.
func DownloadFile(cs *ClientSession, details prompts.ServerDetails) error {
	fmt.Printf("[Server %d] Downloading '%s' to '%s'...\n", cs.serverID, details.RemotePath, details.LocalPath)

	// Open the remote file for reading
	remoteFile, err := cs.sftpClient.Open(details.RemotePath)
	if err != nil {
		return fmt.Errorf("failed to open remote file '%s': %w", details.RemotePath, err)
	}
	defer remoteFile.Close() // Ensure remote file is closed

	// Create the local file for writing
	localFile, err := os.Create(details.LocalPath)
	if err != nil {
		return fmt.Errorf("failed to create local file '%s': %w", details.LocalPath, err)
	}
	defer localFile.Close() // Ensure local file is closed

	// Copy data from remote file to local file
	bytesRead, err := io.Copy(localFile, remoteFile)
	if err != nil {
		return fmt.Errorf("failed to copy data from remote file: %w", err)
	}

	fmt.Printf("[Server %d] Download of '%s' to '%s' COMPLETE (%d bytes).\n", cs.serverID, details.RemotePath, details.LocalPath, bytesRead)
	return nil
}
