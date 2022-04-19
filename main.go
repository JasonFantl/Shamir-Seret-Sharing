package main

import (
	"fmt"

	"github.com/jasonfantl/secretSharing/shamir"
)

func main() {
	requiredShares := 6
	ss, err := shamir.NewSecretSharing([]byte("this is my secret"), requiredShares)
	if err != nil {
		fmt.Println(err)
		return
	}

	shares := make([]shamir.Point, 6+2)
	for i := 0; i < len(shares); i++ {
		shares[i] = ss.GenerateShare()
	}

	decoded, err := shamir.DecodeSecret(shares, requiredShares)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(decoded)
	fmt.Println(string(decoded))
}
