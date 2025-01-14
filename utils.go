package main

import (
	"fmt"
	"os"
	"os/exec"
	"time"
)

// startGlobalCounter starts the global counter counting.
func startGlobalCounter() {
	counter := 0
	for {
		time.Sleep(1 * time.Minute)
		counter++
		log.Info().Msgf("GLOBAL counter: %d", counter)
	}
}

// updateCounter increments the counter for a given endpoint
func updateCounter(endpoint string) {
	mu.Lock()
	defer mu.Unlock()
	log.Debug().Msgf("updating [%s] counter ...", endpoint)
	counters[endpoint]++
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
	log.Info().Msg("API is ready ...")
	log.Info().Msg("checkout http://localhost:8080")
	log.Info().Msg("and please check README.md for further details :)")
}
