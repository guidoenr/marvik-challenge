package main

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

/*
	TODO: this could be improved since it hits the real
	server, instead of mocking. But i find it uneeded at the moment
	since it was never a request from Marvik.
*/

func TestGetUsersWithNameFilter(t *testing.T) {
	// create a test HTTP request to /users with the name filter
	resp, err := http.Get("http://localhost:8080/users?name=Guido")
	if err != nil {
		t.Fatalf("error occurred: %s", err)
	}

	// assert that the response status code is 200 OK
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	// parse the response body into users slice
	var users []User
	err = json.NewDecoder(resp.Body).Decode(&users)
	if err != nil {
		t.Fatalf("error decoding response: %s", err)
	}

	// assert that each user contains 'John' in their name
	for _, user := range users {
		assert.Contains(t, user.Name, "Guido")
	}
}

func TestGetUsersWithEmailFilter(t *testing.T) {
	// create a test HTTP request to /users with the email filter
	resp, err := http.Get("http://localhost:8080/users?email=john.doe@example.com")
	if err != nil {
		t.Fatalf("error occurred: %s", err)
	}

	// assert that the response status code is 200 OK
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	// parse the response body into users slice
	var users []User
	err = json.NewDecoder(resp.Body).Decode(&users)
	if err != nil {
		t.Fatalf("error decoding response: %s", err)
	}

	// assert that the response contains the user with the specified email
	for _, user := range users {
		assert.Equal(t, user.Email, "john.doe@example.com")
	}
}

func TestGetUsersWithOrganizationFilter(t *testing.T) {
	// create a test HTTP request to /users with an organization filter
	resp, err := http.Get("http://localhost:8080/users?organization=Veritone")
	if err != nil {
		t.Fatalf("error occurred: %s", err)
	}

	// assert that the response status code is 200 OK
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	// parse the response body into users slice
	var users []User
	err = json.NewDecoder(resp.Body).Decode(&users)
	if err != nil {
		t.Fatalf("error decoding response: %s", err)
	}

	// assert that the users belong to 'Veritone' organization
	for _, user := range users {
		for _, org := range user.Organizations {
			assert.Contains(t, org.Name, "Veritone")
		}
	}
}

func TestGetOrganizationsWithUsers(t *testing.T) {
	// create a test HTTP request to /organizations
	resp, err := http.Get("http://localhost:8080/organizations")
	if err != nil {
		t.Fatalf("error occurred: %s", err)
	}

	// assert that the response status code is 200 OK
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	// parse the response body into organizations slice
	var organizations []Organization
	err = json.NewDecoder(resp.Body).Decode(&organizations)
	if err != nil {
		t.Fatalf("error decoding response: %s", err)
	}

	// assert that each organization has associated users
	for _, org := range organizations {
		// because hackerone contains no people
		if org.Name != "Hackerone" {
			assert.NotEmpty(t, org.Users)
		}
	}
}

func TestGetUsersWithMultipleFilters(t *testing.T) {
	// create a test HTTP request to /users with multiple filters (name and organization)
	resp, err := http.Get("http://localhost:8080/users?name=John&organization=Veritone")
	if err != nil {
		t.Fatalf("error occurred: %s", err)
	}

	// assert that the response status code is 200 OK
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	// parse the response body into users slice
	var users []User
	err = json.NewDecoder(resp.Body).Decode(&users)
	if err != nil {
		t.Fatalf("error decoding response: %s", err)
	}

	// assert that the response contains users matching both filters
	for _, user := range users {
		assert.Contains(t, user.Name, "John")
		for _, org := range user.Organizations {
			assert.Contains(t, org.Name, "Veritone")
		}
	}
}

func TestHelloWorld(t *testing.T) {
	// create a test HTTP request to / to check if the server is running
	resp, err := http.Get("http://localhost:8080/")
	if err != nil {
		t.Fatalf("error occurred: %s", err)
	}

	// assert that the response status code is 200 OK
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	// parse the response body to check if the server is ready
	var response map[string]string
	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		t.Fatalf("error decoding response: %s", err)
	}

	// assert that the server is ready
	assert.Equal(t, response["status"], "ready")
}
