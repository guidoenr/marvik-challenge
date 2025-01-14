package main

import (
	"database/sql"
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/lib/pq" // Import pq to handle array scanning
	_ "github.com/lib/pq"
	"github.com/rs/zerolog"
)

// declare the global logger and DB connection at the package level
var db *sql.DB
var log zerolog.Logger // declaring it globally for easy access

// declare a global map to track endpoint counters
var counters = make(map[string]int)
var mu sync.Mutex // mutex to protect concurrent access to the counters map

// user struct defines the data model for user
type User struct {
	ID            int      `json:"id"`
	Name          string   `json:"name"`
	Surname       string   `json:"surname"`
	Email         string   `json:"email"`
	Organizations []string `json:"organizations"`
}

// countersHandler returns the number of times each endpoint was accessed
func countersHandler(c *gin.Context) {
	updateCounter("/counters")
	mu.Lock()
	defer mu.Unlock()

	// counters map
	c.JSON(http.StatusOK, counters)
}

// getUsers fetches the users from the database, applying filters if provided
func getUsers(c *gin.Context) {
	updateCounter("/users")

	// Get query parameters (filters)
	filters := map[string]string{
		"name":         c.DefaultQuery("name", ""),
		"surname":      c.DefaultQuery("surname", ""),
		"email":        c.DefaultQuery("email", ""),
		"organization": c.DefaultQuery("organization", ""),
	}

	// Build SQL query based on filters
	baseQuery, args := BuildQuery(filters)

	// Execute the query
	rows, err := db.Query(baseQuery, args...)
	if err != nil {
		log.Error().Err(err).Msg("error fetching users")
		c.JSON(http.StatusInternalServerError, gin.H{"message": "error fetching users"})
		return
	}
	defer rows.Close()

	// Store users
	var users []User
	for rows.Next() {
		var user User
		var organizations pq.StringArray
		err := rows.Scan(&user.ID, &user.Name, &user.Surname, &user.Email, &organizations)
		if err != nil {
			log.Error().Err(err).Msg("error reading data")
			c.JSON(http.StatusInternalServerError, gin.H{"message": "error reading data", "error": err.Error()})
			return
		}

		user.Organizations = []string(organizations)
		users = append(users, user)
	}

	if err := rows.Err(); err != nil {
		log.Error().Err(err).Msg("error getting rows")
		c.JSON(http.StatusInternalServerError, gin.H{"message": "error getting rows"})
		return
	}

	log.Info().Msgf("%d users fetched successfully", len(users))
	c.JSON(http.StatusOK, users)
}

// helloWorld returns a ready message to check if the API is ready
func helloWorld(c *gin.Context) {
	updateCounter("/")
	log.Info().Msg("API is ready ...")
	c.JSON(http.StatusOK, gin.H{"status": "ready"})
}

func main() {
	// GLOBAL counter
	go incrementGlobalCounter()

	connectToDb()
	defer db.Close()

	log = InitLogger()
	banner()

	// endpoints!!!
	router := gin.Default()

	router.GET("/", helloWorld)
	router.GET("/users", getUsers)
	router.GET("/counters", countersHandler)

	// set release and NOT debug mode
	gin.SetMode(gin.ReleaseMode)

	// run the server
	err := router.Run(":8080")
	if err != nil {
		log.Fatal().Err(err).Msg("error starting server")
	}
}
