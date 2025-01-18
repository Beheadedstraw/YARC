package main

import (
	"bufio"
	"flag"
	"fmt"
	"net"
	"os"
	"time"
)

func main() {
	// Define the command-line switch for the passcode
	passcode := flag.String("passcode", "", "Passcode to authenticate with the server")
	serverAddr := flag.String("server", "localhost:8080", "Address of the server")
	serverCommand := flag.String("command", "", "Sends command to the server")
	flag.Parse()

	if *passcode == "" {
		fmt.Println("Error: Passcode is required.")
		os.Exit(1)
	}

	// Connect to the server
	conn, err := net.Dial("tcp", *serverAddr)
	if err != nil {
		fmt.Printf("Error connecting to server: %v\n", err)
		os.Exit(1)
	}
	defer conn.Close()

	// Send the passcode to the server
	fmt.Println("Auth String: ", []byte(*passcode))
	conn.Write([]byte("CSTRING " + *passcode + "\n"))

	// Read the server's response
	scanner := bufio.NewScanner(conn)
	if scanner.Scan() {
		response := scanner.Text()
		fmt.Println("Auth response:", response)
		if response != "AUTH SUCCESSFUL" {
			fmt.Println("Authentication failed. Exiting.")
			return
		}
	}

	// Continue sending commands after successful authentication
	//fmt.Println("Connected. Enter commands:")
	command := []byte(*serverCommand + "\n")
	fmt.Println("Command to send:", command)
	conn.Write(command)
	time.Sleep(1)
	if scanner.Scan() {
		fmt.Println("Server response:", scanner.Text())
	}
}
