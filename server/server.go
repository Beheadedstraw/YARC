package main

import (
	"bufio"
	"fmt"
	"net"
	"strconv"
	"strings"
	"time"
)

// Command structure for map operations
type command struct {
	action string
	key    string
	value  string
	db     int
	result chan string
	date   string
}

func main() {
	// Shared map managed by a goroutine
	commands := make(chan command)

	// Start map manager goroutine
	go mapManager(commands, 10, true)

	// Start the socket server
	address := "localhost:8080"
	listener, err := net.Listen("tcp", address)
	if err != nil {
		fmt.Println("Error starting server:", err)
		return
	}
	defer listener.Close()

	fmt.Println("Server listening on", address)

	for {
		// Accept a new connection
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting connection:", err)
			continue
		}
		// Handle the connection in a new goroutine
		go handleConnection(conn, commands)
	}
}

// Goroutine to manage the shared map
func mapManager(commands chan command, max_db int, debugging bool) {
	es3s_channel := []map[string]map[string][]byte{}
	for dbs := 1; dbs <= max_db; dbs++ {
		es3s_channel = append(es3s_channel, make(map[string]map[string][]byte))

	}

	fmt.Println("Max number of DB's: ", len(es3s_channel))

	for cmd := range commands {
		switch cmd.action {
		case "GET":
			// Handle GET request
			if debugging {
				fmt.Println("Getting:", cmd.key, "in DB", cmd.db)
			}

			if value, ok := es3s_channel[cmd.db][cmd.key]; ok {
				cmd.result <- string(value["data"]) // Don't unmarshal the raw JSON
			} else {
				cmd.result <- "(nil)"
			}

		case "SET":
			// Handle SET request
			es3s_channel[cmd.db][cmd.key] = make(map[string][]byte)
			es3s_channel[cmd.db][cmd.key]["data"] = []byte(cmd.value)
			es3s_channel[cmd.db][cmd.key]["datetime"] = []byte(cmd.date)

			if debugging {
				fmt.Println("Setting:", cmd.key, "in DB", cmd.db, " to ", cmd.value)
			}

			cmd.result <- "OK"

		case "DEL":
			// Handle DELETE request
			if _, ok := es3s_channel[cmd.db][cmd.key]; ok {
				delete(es3s_channel[cmd.db], cmd.key)
				cmd.result <- "OK"
			} else {
				cmd.result <- "(key not found)"
			}

		case "EXIT":
			// Handle DELETE request
			cmd.result <- "CLOSECONN"

		default:
			cmd.result <- "(unknown command)"
		}
	}
}

// Handle individual client connections
func handleConnection(conn net.Conn, commands chan command) {
	defer conn.Close()
	can_continue := true
	fmt.Println("Client connected:", conn.RemoteAddr())

	reader := bufio.NewScanner(conn)
	for reader.Scan() {
		input := strings.TrimSpace(reader.Text())
		parts := strings.Fields(input)
		db_int := int(0)

		if len(parts) == 0 {
			conn.Write([]byte("Invalid command\n"))
			continue
		}
		db_int, err := strconv.Atoi(parts[0])
		if err == nil {
			//
		} else {
			can_continue = false
		}
		action := strings.ToUpper(parts[1])
		var key, value string
		if len(parts) > 2 {
			key = strings.TrimSpace(parts[2])
		}
		if len(parts) > 3 {
			value = strings.Join(parts[3:], " ")
		}

		date := time.Now().String()

		// Handle commands
		if can_continue {
			result := make(chan string)
			commands <- command{action: action, key: key, date: date, value: value, db: db_int, result: result}
			response := <-result
			if response == "CLOSECONN" {
				conn.Write([]byte(response + "\n"))
				conn.Close()
			}
			conn.Write([]byte(response + "\n"))
		}
	}

	fmt.Println("Client disconnected:", conn.RemoteAddr())
}
