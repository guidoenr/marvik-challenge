package main

import (
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"gorm.io/gorm"
)

// DB and Logger must be global.
var db *gorm.DB
var log zerolog.Logger

// counters map to handle each endpoint counter
var counters = make(map[string]int)
var mu sync.Mutex

// getCounters returns the number of times each endpoint was accessed
func getCounters(c *gin.Context) {
	updateCounter("/counters")
	mu.Lock()
	defer mu.Unlock()

	// Return counters map as JSON
	c.JSON(http.StatusOK, counters)
}

// getOrganizations fetch all the organizations with their associated users
func getOrganizations(c *gin.Context) {
	updateCounter("/organizations")

	// fetch organizations with preloaded users
	var organizations []Organization
	err := db.Preload("Users").Find(&organizations).Error
	if err != nil {
		log.Error().Err(err).Msg("error fetching organizations")
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "error fetching organizations",
			"error":   err.Error(),
		})
		return
	}

	log.Info().Msgf("%d organizations fetched successfully", len(organizations))
	c.JSON(http.StatusOK, organizations)
}

// getUsers fetch all the users from the DB
func getUsers(c *gin.Context) {
	updateCounter("/users")

	// get the query parameters (filters)
	filters := map[string]string{
		"name":         c.DefaultQuery("name", ""),
		"surname":      c.DefaultQuery("surname", ""),
		"email":        c.DefaultQuery("email", ""),
		"organization": c.DefaultQuery("organization", ""),
	}

	// preload the full Organizations struct
	query := db.Preload("Organizations")

	// apply filters if they are provided
	if filters["name"] != "" {
		query = query.Where("name ILIKE ?", "%"+filters["name"]+"%")
	}
	if filters["surname"] != "" {
		query = query.Where("surname ILIKE ?", "%"+filters["surname"]+"%")
	}
	if filters["email"] != "" {
		query = query.Where("email ILIKE ?", "%"+filters["email"]+"%")
	}
	if filters["organization"] != "" {
		query = query.Joins("JOIN user_organizations ON user_organizations.user_id = users.id").
			Joins("JOIN organizations ON organizations.id = user_organizations.organization_id").
			Where("organizations.name ILIKE ?", "%"+filters["organization"]+"%")
	}

	// fetch the users
	var users []User
	err := query.Find(&users).Error
	if err != nil {
		log.Error().Err(err).Msg("error fetching users")
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "error fetching users",
			"error":   err.Error(),
		})
		return
	}

	log.Info().Msgf("%d users fetched successfully", len(users))
	c.JSON(http.StatusOK, users)
}

// getUserOrganizations fetch all user-organization relationships from the DB
func getUserOrganizations(c *gin.Context) {
	updateCounter("/user_organizations")

	// Fetch all user-organization associations
	var userOrganizations []struct {
		UserID           uint   `json:"user_id"`
		UserName         string `json:"user_name"`
		OrganizationID   uint   `json:"organization_id"`
		OrganizationName string `json:"organization_name"`
	}

	// Perform the query with joins to get the related user and organization data
	err := db.Table("user_organizations").
		Select("user_organizations.user_id, users.name as user_name, user_organizations.organization_id, organizations.name as organization_name").
		Joins("JOIN users ON users.id = user_organizations.user_id").
		Joins("JOIN organizations ON organizations.id = user_organizations.organization_id").
		Find(&userOrganizations).Error

	if err != nil {
		log.Error().Err(err).Msg("error fetching user-organization associations")
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "error fetching user-organization associations",
			"error":   err.Error(),
		})
		return
	}

	log.Info().Msgf("%d user-organization associations fetched successfully", len(userOrganizations))
	c.JSON(http.StatusOK, userOrganizations)
}

// helloWorld returns a ready message to check if the API is ready
func helloWorld(c *gin.Context) {
	updateCounter("/")
	c.JSON(http.StatusOK, gin.H{"status": "ready"})
}

func main() {
	go startGlobalCounter()

	connectToDb()

	// close the DB when the app ends
	defer func() {
		sqlDB, err := db.DB()
		if err == nil {
			sqlDB.Close()
		}
	}()

	log = InitLogger()
	banner()

	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()

	router.GET("/", helloWorld)
	router.GET("/users", getUsers)
	router.GET("/counters", getCounters)
	router.GET("/organizations", getOrganizations)
	router.GET("/user_organizations", getUserOrganizations)

	err := router.Run(":8080")
	if err != nil {
		log.Fatal().Err(err).Msg("error starting server")
	}
}
