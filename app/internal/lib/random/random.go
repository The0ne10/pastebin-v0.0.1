package random

import (
	"crypto/rand"
	"math/big"
)

func GenerateRandomString(length int) string {
	const op = "lib.random.random"
	chars := []rune(
		"QWERTYUIOPASDFGHJKLZXCVBNM" +
			"qwertyuiopasdfghjklzxcvbnm" +
			"0123456789",
	)

	sourceLength := big.NewInt(int64(len(chars)))

	randomIndices := make([]int64, length)

	for i := 0; i < length; i++ {
		randomIndex, _ := rand.Int(rand.Reader, sourceLength)
		randomIndices[i] = randomIndex.Int64()
	}

	var result string
	for _, index := range randomIndices {
		result += string(chars[index])
	}

	return result
}
