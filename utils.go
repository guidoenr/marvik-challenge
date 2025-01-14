package main

import (
	"database/sql"
	"fmt"
	"os"
	"os/exec"
	"time"
)

// incrementCounter is the counter each 10 seconds
func incrementCounter() {
	counter := 0
	for {
		time.Sleep(10 * time.Second)
		counter++
		log.Info().Msgf("GLOBAL counter: %d", counter)
	}
}

// connectToDb tries to connect to the PostgreSQL DB
// TODO, this could be improved
func connectToDb() {
	var err error
	connStr := "host=localhost port=5432 user=postgres password=mysecretpassword dbname=postgres sslmode=disable"
	db, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal().Err(err).Msg("error connecting to the database")
	}
}

// stringJoin is a helper function to join strings with a separator
func stringJoin(parts []string, separator string) string {
	result := ""
	for i, part := range parts {
		if i > 0 {
			result += separator
		}
		result += part
	}
	return result
}

// banner clears the terminal screen and prints a custom banner message
func banner() {
	// clean
	cmd := exec.Command("clear")
	cmd.Stdout = os.Stdout
	cmd.Run()

	fmt.Println("=====================================")
	fmt.Println("   Marvik Challenge by @guidoenr     ")
	fmt.Println("=====================================")
}
