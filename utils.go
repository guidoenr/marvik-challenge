package main

import (
	"database/sql"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

// BuildQuery dynamically builds the SQL query with filters
func BuildQuery(filters map[string]string) (string, []interface{}) {
	var conditions []string
	var args []interface{}
	argIndex := 1

	baseQuery := `
		SELECT u.id, u.name, u.surname, u.email, array_agg(o.name) AS organizations
		FROM users u
		LEFT JOIN user_organizations uo ON u.id = uo.user_id
		LEFT JOIN organizations o ON o.id = uo.organization_id
	`

	for field, value := range filters {
		if value != "" {
			if field == "organization" {
				// Adjusted for the organization filter
				conditions = append(conditions, "LOWER(o.name) = LOWER($"+strconv.Itoa(argIndex)+")")
			} else {
				conditions = append(conditions, "LOWER(u."+field+") = LOWER($"+strconv.Itoa(argIndex)+")")
			}
			args = append(args, value)
			argIndex++
		}
	}

	if len(conditions) > 0 {
		baseQuery += " WHERE " + strings.Join(conditions, " AND ")
	}

	baseQuery += " GROUP BY u.id"
	return baseQuery, args
}

// incrementCounter is the counter each 1 minute
func incrementGlobalCounter() {
	counter := 0
	for {
		time.Sleep(1 * time.Minute)
		counter++
		log.Info().Msgf("GLOBAL counter (updatead each 1 minute): %d", counter)
	}
}

// incrementCounter increments the counter for a given endpoint
func updateCounter(endpoint string) {
	mu.Lock()
	defer mu.Unlock()
	counters[endpoint]++
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
