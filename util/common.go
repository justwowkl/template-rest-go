package util

// Init init
func Init() {
	pgInit()
	redisInit()
	gcacheInit()
}

// Health health
func Health() bool {
	ch := make(chan bool, 3)
	go func() { ch <- pgHealth() }()
	go func() { ch <- redisHealth() }()
	go func() { ch <- gcacheHealth() }()

	// wait for goroutines to finish
	for i := 0; i < 3; i++ {
		if !<-ch {
			return false
		}
	}
	return true
}
