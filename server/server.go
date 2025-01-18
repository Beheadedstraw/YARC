package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"runtime"
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
	//set standard variables
	version := "0.1.1"
	max_db := int(10)
	cstring := "1234abcd"

	//open log file for... logging?
	file, err := os.OpenFile("yarc-server.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatalf("Failed to open log file: %s", err)
	}
	defer file.Close()
	commands := make(chan command)
	log.SetOutput(file)
	log.Println("\n------------------------------\nStarting YARC Server: ", version, "\n------------------------------")
	// Start map manager goroutine
	go mapManager(commands, max_db, true)

	// Start the socket server
	address := "localhost:8080"
	listener, err := net.Listen("tcp", address)

	if err != nil {
		fmt.Println("Error starting server:", err)
		return
	}
	defer listener.Close()

	log.Println("Server listening on", address)

	for {
		// Accept a new connection
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting connection:", err)
			continue
		}
		// Handle the connection in a new goroutine
		go handleConnection(conn, commands, cstring)
	}
}

// Goroutine to manage the shared map
func mapManager(commands chan command, max_db int, debugging bool) {

	//Initialize DB chache map
	yarc_channel := []map[string]map[string][]byte{}
	//yarc_sort := [][]string{}
	for dbs := 1; dbs <= max_db; dbs++ {
		yarc_channel = append(yarc_channel, make(map[string]map[string][]byte))
	}

	//Initialize sort map

	//fmt.Println("Max number of DB's: ", len(yarc_channel))
	log.Println("DBs Initialized: ", len(yarc_channel))
	var memStats runtime.MemStats
	for cmd := range commands {
		switch cmd.action {
		case "GET":
			// Handle GET request
			if debugging {
				log.Println("Getting:", cmd.key, "in DB", cmd.db)
			}

			if value, ok := yarc_channel[cmd.db][cmd.key]; ok {
				cmd.result <- string(value["data"]) // Don't unmarshal the raw JSON
			} else {
				cmd.result <- "(nil)"
			}

		case "SET":
			// Handle SET request
			yarc_channel[cmd.db][cmd.key] = make(map[string][]byte)
			yarc_channel[cmd.db][cmd.key]["data"] = []byte(cmd.value)
			yarc_channel[cmd.db][cmd.key]["datetime"] = []byte(cmd.date)

			if debugging {
				log.Println("Setting:", cmd.key, "in DB", cmd.db, " to ", cmd.value)
			}

			cmd.result <- "OK"

		case "DEL":
			// Handle DELETE request
			if _, ok := yarc_channel[cmd.db][cmd.key]; ok {
				delete(yarc_channel[cmd.db], cmd.key)
				cmd.result <- "OK"
			} else {
				cmd.result <- "(key not found)"
			}

		case "STATS":
			//get DB usage
			for db := 0; db < len(yarc_channel); db++ {
				//arraysize := int(unsafe.Sizeof(yarc_channel[db]))
				//fmt.Printf("\n - Size of DB %d is %d", db, arraysize)
			}

			//get general memory usage
			runtime.ReadMemStats(&memStats)
			totalAll := (float64(memStats.Alloc) / (1024 * 1024))
			gcAll := (float64(memStats.Sys) / (1024 * 1024))
			stats := fmt.Sprintf("TotalAllocated: %.2fMB, System Usaged: %.2fMB", totalAll, gcAll)
			cmd.result <- stats

		case "PURGE":
			// Handle PURGE request
			yarc_channel[cmd.db] = make(map[string]map[string][]byte)
			cmd.result <- "OK"

		case "EXIT":
			// Handle DELETE request
			cmd.result <- "CLOSECONN"

		default:
			cmd.result <- ("(unknown command " + string(cmd.db) + ")")
		}
	}
}

// Handle individual client connections
func handleConnection(conn net.Conn, commands chan command, cstring string) {
	defer conn.Close()
	isAuthed := false
	//can_continue := true
	//fmt.Println("Client connected:", conn.RemoteAddr())
	log.Println("Client connected:", conn.RemoteAddr())

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
			// Handle commands
			if isAuthed {
				action := ""
				if len(parts) > 1 {
					action = strings.ToUpper(parts[1])
				} else {
					conn.Write([]byte("Invalid Command String\n"))
				}

				var key, value string
				if len(parts) > 2 {
					key = strings.TrimSpace(parts[2])
				}
				if len(parts) > 3 {
					value = strings.Join(parts[3:], " ")
				}

				date := time.Now().String()
				result := make(chan string)
				commands <- command{action: action, key: key, date: date, value: value, db: db_int, result: result}
				response := <-result
				if response == "CLOSECONN" {
					conn.Write([]byte(response + "\n"))
					conn.Close()
				}
				conn.Write([]byte(response + "\n"))
			}
		} else {
			auth := strings.ToUpper(parts[0])
			if auth == "CSTRING" {
				if parts[1] == cstring {
					conn.Write([]byte("AUTH SUCCESSFUL\n"))
					isAuthed = true
				} else {
					isAuthed = false
					//can_continue = false
				}
			} else {
				conn.Write([]byte("Commands Must Start With DB Integer\n"))
			}
		}

		if !isAuthed {
			conn.Write([]byte("Authentication required\n"))
			conn.Close()
		}
	}

	/* We'll see how this goes, a rewrite in C or Rust(shiver) might be a better alternative if the GC impacts
	performance too much */
	//runtime.GC()
	log.Println("Client disconnected:", conn.RemoteAddr())
}
