package main

import (
	"fmt"
)

func main() {

	requiredShareCount := uint(6)
	secret := "this is my secret"

	encodedPolynomial := encodeSecret([]byte(secret), requiredShareCount)

	shares := make([]Point, requiredShareCount+2)
	for i := 0; i < len(shares); i++ {
		// make sure not to let a share be at x=0 as that gives the first 8 bytes of the secret
		shares[i] = encodedPolynomial.eval(uint64(i + 2))
	}

	decoded := decodeSecret(shares, requiredShareCount)
	fmt.Println(decoded)
	fmt.Println(string(decoded))
}
