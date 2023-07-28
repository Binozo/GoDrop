package utils

import (
	"log"
	"os/user"
)

// https://stackoverflow.com/a/66624820/14119471
func IsRoot() bool {
	currentUser, err := user.Current()
	if err != nil {
		log.Fatalf("[isRoot] Unable to get current user: %s", err)
	}
	return currentUser.Username == "root"
}
