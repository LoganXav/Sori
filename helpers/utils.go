package helpers

import (
	"math/rand"
	"time"
)

func RandomNumber(numLength int) string {
	rand.Seed(time.Now().UnixNano())
	numbers := make([]byte, numLength)

	for i := 0; i < numLength; i++ {
		numbers[i] = byte(rand.Intn(10)) + '0'
	}

	return string(numbers)
}

