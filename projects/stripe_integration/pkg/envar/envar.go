package envar

import (
	"log"
	"os"
)

var keys = []string{
	"stripe_key",
	"stripe_secret",
}

// Validate validate env varibales
func Validate() {
	for _, key := range keys {
		if os.Getenv(key) == "" {
			log.Fatalf("%v is empty", key)
		}
	}
}
