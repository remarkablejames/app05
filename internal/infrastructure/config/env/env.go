package env

import (
	"log"
	"os"
	"strconv"
	"time"
)

func GetString(key, fallback string) string {
	val, ok := os.LookupEnv(key)
	if !ok {
		return fallback
	}

	return val
}

func GetInt(key string, fallback int) int {
	val, ok := os.LookupEnv(key)
	if !ok {
		return fallback
	}

	intVal, err := strconv.Atoi(val)
	if err != nil {
		log.Printf("\033[1;31mFailed to parse int for key %s: %v\033[0m", key, err)
		// Return fallback if parsing fails
		return fallback
	}

	return intVal
}

func GetDuration(key string, fallback time.Duration) time.Duration {
	val, ok := os.LookupEnv(key)
	if !ok {
		return fallback
	}

	durationVal, err := time.ParseDuration(val)
	if err != nil {
		log.Printf("\033[1;31mFailed to parse duration for key %s: %v\033[0m", key, err)
		// Return fallback if parsing fails
		return fallback
	}

	return durationVal
}

func GetBool(key string, fallback bool) bool {
	val, ok := os.LookupEnv(key)
	if !ok {
		return fallback
	}

	boolVal, err := strconv.ParseBool(val)
	if err != nil {
		log.Printf("\033[1;31mFailed to parse bool for key %s: %v\033[0m", key, err)
		// Return fallback if parsing fails
		return fallback
	}

	return boolVal
}

func GetFloat64(key string, fallback float64) float64 {
	val, ok := os.LookupEnv(key)
	if !ok {
		return fallback
	}

	floatVal, err := strconv.ParseFloat(val, 64)
	if err != nil {
		log.Printf("\033[1;31mFailed to parse float for key %s: %v\033[0m", key, err)
		// Return fallback if parsing fails
		return fallback
	}

	return floatVal
}
