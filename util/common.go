package util

// Init init
func Init() {
	pgInit()
	redisInit()
	gcacheInit()
}

// Health health
func Health() bool {
	if pgHealth() && redisHealth() && gcacheHealth() {
		return true
	}
	return false
}
