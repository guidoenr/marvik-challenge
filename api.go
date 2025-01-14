package main

import (
	"database/sql"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/lib/pq" // Import pq to handle array scanning
	_ "github.com/lib/pq"
	"github.com/rs/zerolog"
)

// declare the global logger and DB connection at the package level
var db *sql.DB
var log zerolog.Logger // declaring it globally for easy access

// user struct defines the data model for user
type User struct {
	ID            int      `json:"id"`
	Name          string   `json:"name"`
	Surname       string   `json:"surname"`
	Email         string   `json:"email"`
	Organizations []string `json:"organizations"`
}

// helloWorld just returns a ready message to check if the API is ready
func helloWorld(c *gin.Context) {
	log.Info().Msg("API is ready")
	c.JSON(http.StatusOK, gin.H{"apiStatus": "ready"})
}

// getUsers fetches the users from the database, applying filters if provided
func getUsers(c *gin.Context) {
	// get query parameters (filters)
	name := c.DefaultQuery("name", "") // empty string if not provided
	surname := c.DefaultQuery("surname", "")
	email := c.DefaultQuery("email", "")
	organization := c.DefaultQuery("organization", "")

	// start building the query
	baseQuery := `
		SELECT u.id, u.name, u.surname, u.email, array_agg(o.name) AS organizations
		FROM users u
		LEFT JOIN user_organizations uo ON u.id = uo.user_id
		LEFT JOIN organizations o ON o.id = uo.organization_id
	`

	// condition to apply if filters are provided
	var conditions []string
	var args []interface{}
	var argIndex int

	// dynamically build WHERE conditions based on query params
	// LOWER is because PSQL is case sensitive by default
	if name != "" {
		conditions = append(conditions, "LOWER(u.name) = LOWER($"+strconv.Itoa(argIndex+1)+")")
		args = append(args, name)
		argIndex++
	}
	if surname != "" {
		conditions = append(conditions, "LOWER(u.surname) = LOWER($"+strconv.Itoa(argIndex+1)+")")
		args = append(args, surname)
		argIndex++
	}
	if email != "" {
		conditions = append(conditions, "LOWER(u.email) = LOWER($"+strconv.Itoa(argIndex+1)+")")
		args = append(args, email)
		argIndex++
	}
	if organization != "" {
		conditions = append(conditions, "LOWER(o.name) = LOWER($"+strconv.Itoa(argIndex+1)+")")
		args = append(args, organization)
		argIndex++
	}

	// if there are conditions, add them to the query
	if len(conditions) > 0 {
		baseQuery += " WHERE " + stringJoin(conditions, " AND ")
	}

	// group and order by user ID
	baseQuery += " GROUP BY u.id"

	// execute the query
	rows, err := db.Query(baseQuery, args...)
	if err != nil {
		log.Error().Err(err).Msg("error fetching users")
		c.JSON(http.StatusInternalServerError, gin.H{"message": "error fetching users"})
		return
	}
	defer rows.Close()

	// hold users
	var users []User
	for rows.Next() {
		var user User
		var organizations pq.StringArray // (for handling psql arrays)

		// read each user data and the organizations they belong to
		err := rows.Scan(&user.ID, &user.Name, &user.Surname, &user.Email, &organizations)
		if err != nil {
			log.Error().Err(err).Msg("error reading data")
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "error reading data",
				"error":   err.Error(),
			})
			return
		}

		// store organizations and save the users
		user.Organizations = []string(organizations)
		users = append(users, user)
	}

	// check for row errors
	err = rows.Err()
	if err != nil {
		log.Error().Err(err).Msg("error getting rows")
		c.JSON(http.StatusInternalServerError, gin.H{"message": "error getting rows"})
		return
	}

	log.Info().Msgf("%d users fetched successfully", len(users))
	c.JSON(http.StatusOK, users)
}

func main() {
	// start the counter in a go routine
	// (no need to add a select{} to keep main running since the)
	// router is block
	go incrementCounter()

	// connect to the DB
	connectToDb()
	defer db.Close()

	// init logger and display banner
	log = InitLogger()
	banner()

	// set up the router
	router := gin.Default()

	// endpoints
	router.GET("/", helloWorld)
	router.GET("/users", getUsers)

	// start the server
	err := router.Run(":8080")
	if err != nil {
		log.Fatal().Err(err).Msg("error starting server")
	}
}
