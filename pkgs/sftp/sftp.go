package sftp

import (
	"fmt"
	"io"
	"kavigo/pkgs/globvars"
	"log"
	"os"

	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
)

var sftpClient *sftp.Client
var conn *ssh.Client

func CreateRemoteConn() error {

	host := globvars.Host
	port := globvars.Port
	user := globvars.User
	password := globvars.Password
	key, err := os.ReadFile(globvars.SshKey)
	if err != nil {
		log.Println("Unable to read private key (encrypted keys are not supported, yet) ")
		return err
	}

	// Create a signer for private key.
	signer, err := ssh.ParsePrivateKey(key)
	if err != nil {
		log.Println("Unable to parse private key \nTrying password authentication...")
		return err
	}

	// Create a config
	config := &ssh.ClientConfig{
		User: user,
		Auth: []ssh.AuthMethod{
			ssh.Password(password),
			ssh.PublicKeys(signer),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(), //ssh.FixedHostKey(hostKey),
	}

	conn, err = ssh.Dial("tcp", fmt.Sprintf("%s:%d", host, port), config)
	if err != nil {
		fmt.Println("Failed to connect to ssh server")
		return err
	}
	//defer conn.Close()

	// Open SFTP session
	sftpClient, err = sftp.NewClient(conn)
	if err != nil {
		fmt.Println("Failed to open SFTP session ")
		return err
	}
	fmt.Println("connection successful")
	//defer sftpClient.Close()

	return nil

}

func CopyFilesToRemote(origin, remoteToBe, remoteDirToMk string) {

	err := sftpClient.MkdirAll(remoteDirToMk)
	if err != nil {
		fmt.Println("Cannot make directory")
		log.Fatal(err)
	}

	localFile, err := os.Open(origin)
	if err != nil {
		fmt.Println("Faield to open local file:", err)
		return
	}
	defer localFile.Close()

	remoteFile, err := sftpClient.Create(remoteToBe)
	if err != nil {
		fmt.Println("Failed to create remote file:", err)
		return
	}
	defer remoteFile.Close()

	_, err = io.Copy(remoteFile, localFile)
	if err != nil {
		fmt.Println("Failed to upload file:", err)
		return
	}
}

func CloseRemoteConn() {

	// Close SFTP session
	sftpClient.Close()

	// Close SSH connection
	conn.Close()

}
