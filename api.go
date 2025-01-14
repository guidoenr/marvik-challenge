package main

import (
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/lib/pq" // Import pq to handle array scanning
	_ "github.com/lib/pq"
	"github.com/rs/zerolog"
)

// Declare the global logger and DB connection at the package level
var db *sql.DB
var log zerolog.Logger // Declaring it globally for easy access

// User struct defines the data model for user
type User struct {
	ID            int      `json:"id"`
	Name          string   `json:"name"`
	Surname       string   `json:"surname"`
	Email         string   `json:"email"`
	Organizations []string `json:"organizations"`
}

// connectToDb tries to connect to the postgreSQL
func connectToDb() {
	var err error
	connStr := "host=localhost port=5432 user=postgres password=mysecretpassword dbname=postgres sslmode=disable"
	db, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal().Err(err).Msg("error connecting to the database")
	}
}

// getUsers fetches the users from the database
func getUsers(c *gin.Context) {
	// select query
	rows, err := db.Query(`
		SELECT u.id, u.name, u.surname, u.email, array_agg(o.name) AS organizations
		FROM users u
		LEFT JOIN user_organizations uo ON u.id = uo.user_id
		LEFT JOIN organizations o ON o.id = uo.organization_id
		GROUP BY u.id
	`)
	if err != nil {
		log.Error().Err(err).Msg("error fetching users")
		c.JSON(http.StatusInternalServerError, gin.H{"message": "error fetching users"})
		return
	}
	defer rows.Close()

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
		user.Organizations = []string(organizations) // Convert pq.StringArray to []string
		users = append(users, user)
	}

	err = rows.Err()
	if err != nil {
		log.Error().Err(err).Msg("error getting rows")
		c.JSON(http.StatusInternalServerError, gin.H{"message": "error getting rows"})
		return
	}

	log.Info().Msg("Users fetched successfully")
	c.JSON(http.StatusOK, users)
}

// helloWorld just returns a ready message to check the api is ready
func helloWorld(c *gin.Context) {
	log.Info().Msg("API is ready")
	c.JSON(http.StatusOK, gin.H{"apiStatus": "ready"})
}

func main() {
	log = InitLogger()

	connectToDb()
	defer db.Close()

	// set up the router
	router := gin.Default()

	// endpoints!!!
	router.GET("/", helloWorld)
	router.GET("/users", getUsers)

	// start the server
	err := router.Run(":8080")
	if err != nil {
		log.Fatal().Err(err).Msg("error starting server")
	}
}
