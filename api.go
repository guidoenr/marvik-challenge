package main

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

var db *sql.DB

// User struct defines the data model for user
type User struct {
	ID      int    `json:"id"`
	Name    string `json:"name"`
	Surname string `json:"surname"`
	Email   string `json:"email"`
}

// connectToDb tries to connect to the postgreSQL
// TODO: leaving here private keys is not good
func connectToDb() {
	var err error
	connStr := "host=localhost port=5432 user=postgres password=mysecretpassword dbname=postgres sslmode=disable"
	db, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("error connecting to the database: %v", err)
	}
}

// getUsers fetches the users from the database
func getUsers(c *gin.Context) {
	rows, err := db.Query("SELECT id, name, surname, email FROM users")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "error fetching users"})
		return
	}
	defer rows.Close()

	var users []User
	for rows.Next() {
		var user User
		if err := rows.Scan(&user.ID, &user.Name, &user.Surname, &user.Email); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "error reading data"})
			return
		}
		users = append(users, user)
	}

	if err := rows.Err(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "error getting rows"})
		return
	}

	c.JSON(http.StatusOK, users)
}

// helloWorld just returns a ready message to check the api is ready
func helloWorld(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"apiStatus": "ready"})
}

func main() {
	connectToDb()
	defer db.Close()

	// set up up the Gin router
	router := gin.Default()

	// define the API endpoint
	router.GET("/users", getUsers)
	router.GET("/", helloWorld)

	// run the server
	if err := router.Run(":8080"); err != nil {
		log.Fatalf("error starting server: %v", err)
	}
}
