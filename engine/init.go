package engine

import (
	"crypto/rand"
	"math/big"
	"time"
)

var globalSeed int64

// init generates a random seed when the package is initialized
func init() {
	// Generate a truly random seed
	if randomBig, err := rand.Int(rand.Reader, big.NewInt(1000000)); err == nil {
		globalSeed = randomBig.Int64()
	} else {
		// Fallback to time-based seed if crypto/rand fails
		globalSeed = time.Now().UnixNano() % 1000000
	}
}
