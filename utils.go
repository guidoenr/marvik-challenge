package main

import (
	"fmt"
	"os"
	"os/exec"
	"time"
)

// banner clears the terminal screen and prints the banner message
func banner() {
	// clear the terminal
	cmd := exec.Command("clear")
	cmd.Stdout = os.Stdout
	cmd.Run()

	fmt.Println("=====================================")
	fmt.Println("   Marvik Challenge by @guidoenr     ")
	fmt.Println("=====================================")
	log.Info().Msg("APi is ready !")
	log.Info().Msg("checkout README.md for further details on usage :)")
	fmt.Println("-----------------")
	log.Info().Msg("endpoints:")
	log.Info().Msg("http://localhost:8080/users")
	log.Info().Msg("http://localhost:8080/organizations")
	log.Info().Msg("http://localhost:8080/counters")
	log.Info().Msg("http://localhost:8080/")
}

// startGlobalCounter starts the global counter counting.
func startGlobalCounter() {
	counter := 0
	for {
		time.Sleep(1 * time.Minute)
		counter++
		log.Info().Msgf("GLOBAL counter: %d", counter)
	}
}
