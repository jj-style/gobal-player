package config

import "os"

// MustHomeDir gets the user's home directory path
// or panics if it fails
func MustHomeDir() string {
	dir, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}
	return dir
}
