package config

import (
	"os"
	"strconv"
)

func Get(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}

func GetInt(key string, def int) int {
	if v := os.Getenv(key); v != "" {
		if n, err := strconv.Atoi(v); err == nil {
			return n
		}
	}
	return def
}
