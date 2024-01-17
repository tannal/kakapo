package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/chzyer/readline"
	"github.com/tannal/kakapo/lsm"
)

func main() {
	dbPath := "./db"
	db := lsm.New(dbPath, 1024) // Assuming 1024 is a reasonable buffer size for the in-memory table

	// Create a readline instance
	rl, err := readline.New("> ")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to create readline instance: %v\n", err)
		os.Exit(1)
	}
	defer rl.Close()

	fmt.Println("LSM Tree REPL")
	fmt.Println("Enter commands like 'set key value', 'get key', 'delete key', 'reset', or 'exit'.")

	for {
		line, err := rl.Readline()
		if err != nil {
			if err == readline.ErrInterrupt {
				continue
			} else {
				break
			}
		}

		args := strings.Fields(line)
		if len(args) == 0 {
			continue
		}

		switch args[0] {
		case "set":
			if len(args) != 3 {
				fmt.Println("Usage: set key value")
				continue
			}
			db.Set(args[1], []byte(args[2]))
			fmt.Printf("Set key '%s'\n", args[1])

		case "get":
			if len(args) != 2 {
				fmt.Println("Usage: get key")
				continue
			}
			value, exists := db.Get(args[1])
			if exists {
				fmt.Printf("Value: %s\n", string(value))
			} else {
				fmt.Println("Key not found")
			}

		case "delete":
			if len(args) != 2 {
				fmt.Println("Usage: delete key")
				continue
			}
			db.Delete(args[1])
			fmt.Printf("Deleted key '%s'\n", args[1])

		case "reset":
			db.ResetDB()
			fmt.Println("Database reset")

		case "exit", "quit":
			fmt.Println("Exiting REPL")
			return

		default:
			fmt.Println("Unknown command")
		}
	}
}
