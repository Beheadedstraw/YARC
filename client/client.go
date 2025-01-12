package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
)

func main() {
	// Check if the user provided a command as an argument
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run client.go <command>")
		fmt.Println("Example: go run client.go 'SET mykey myvalue'")
		return
	}

	// Get the command from the command line arguments
	command := strings.Join(os.Args[1:], " ")
	fmt.Println(os.Args)

	// Connect to the server
	serverAddress := "localhost:8080" // Update if server uses a different address or port
	conn, err := net.Dial("tcp", serverAddress)
	if err != nil {
		fmt.Println("Error connecting to server:", err)
		return
	}
	defer conn.Close()

	// Send the command to the server
	fmt.Println("Sending: ", command)
	_, err = conn.Write([]byte(command + "\n"))
	if err != nil {
		fmt.Println("Error sending command:", err)
		return
	}

	// Read the server's response
	reader := bufio.NewReader(conn)
	response, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println("Error reading response:", err)
		return
	}

	// Print the server's response
	fmt.Print("Response from server: ", response)
}
