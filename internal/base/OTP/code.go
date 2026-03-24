package otp

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"strings"
)

func GenerateNumericCode(length int) (string, error) {
	if length <= 0 {
		return "", fmt.Errorf("invalid code length: %d", length)
	}
	var sb strings.Builder
	for i := 0; i < length; i++ {
		n, err := rand.Int(rand.Reader, big.NewInt(10))
		if err != nil {
			return "", fmt.Errorf("generate random number: %w", err)
		}
		sb.WriteString(n.String())
	}
	return sb.String(), nil
}
