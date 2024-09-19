package util

import (
	"math/rand"
	"strings"
	"time"
)

const alphabet = "abcdefghijklmnopqrstuvwxyz"

// init seeds the random number generator using the current time in nanoseconds. This ensures that random numbers generated in each execution are different.
func init() {
	rand.Seed(time.Now().UnixNano())
}

// RandomInt generates a random integer between the provided min and max values (inclusive).
func RandomInt(min, max int64) int64 {
	// rand.Int63n generates a random number between 0 and (max-min),
	return min + rand.Int63n(max-min+1)
}

// RandomString generates a random string of length 'n' using lowercase English letters.
func RandomString(n int) string {
	// Create a string builder to efficiently build the string.
	var sb strings.Builder
	// k stores the length of the alphabet, which is 26 in this case.
	k := len(alphabet)

	// Loop 'n' times to generate 'n' random characters.
	for i := 0; i < n; i++ {
		// Generate a random index between 0 and 25 and select a character from the alphabet.
		c := alphabet[rand.Intn(k)]
		// Write the selected character to the string builder.
		sb.WriteByte(c)
	}
	// Convert the string builder contents to a string and return it.
	return sb.String()
}

// RandomOwner generates a random owner name, which is a random string of 6 characters.
func RandomOwner() string {
	// Calls RandomString with a fixed length of 6.
	return RandomString(6)
}

// RandomMoney generates a random money amount between 0 and 10,000.
func RandomMoney() int64 {
	// Calls RandomInt to generate a random integer between 0 and 10,000.
	return RandomInt(0, 10000)
}

// RandomCurrency selects a random currency from a predefined list (USD, THB).
func RandomCurrency() string {
	currencies := []string{USD, THB}
	// n stores the number of available currencies.
	n := len(currencies)
	// Select a random currency from the slice using a random index.
	return currencies[rand.Intn(n)]
}
