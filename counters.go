package main

// this channel saves the 'update counters' tasks for each endpoint
// like 'Update /users endpoint +1'
// and then those tasks are performed by a routine
var counterUpdates chan string

// counters contains each endpoint counter
var counters = make(map[string]int)

// startCountersManager kicks off the routine that listens
// to the counterUpdates channel for updates
func startCountersManager() {
	log.Debug().Msg("initializing counters map")
	counters = map[string]int{
		"/":              0,
		"/users":         0,
		"/counters":      0,
		"/organizations": 0,
	}

	// channel to hold the ''update tasks''
	counterUpdates = make(chan string)

	// this routine listens to the channel and the do the updates
	go func() {
		for endpoint := range counterUpdates {
			counters[endpoint]++
			log.Info().Str("endpoint", endpoint).Msgf("counter updated: %d", counters[endpoint])
		}
	}()
}

// updateCounter increments the counter for a given endpoint
// (sends the endpoint name to the channel)
func updateCounter(endpoint string) {
	counterUpdates <- endpoint
	log.Debug().Str("endpoint", endpoint).Msg("added 'update endpoint' to channel ...")
}
