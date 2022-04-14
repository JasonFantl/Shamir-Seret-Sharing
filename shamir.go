package main

import (
	"encoding/binary"
	"fmt"
	"math"
)

func encodeSecret(bytes []byte, coefficientCount uint) PolynomialField {
	var prime uint64 = math.MaxUint64 - 58 // https://primes.utm.edu/lists/2small/0bit.html
	// this large prime lets us encode information into the first 63 bits of uint64

	// well just use every 7 bytes since it will fit in our encoding, although its not as storage efficient as it could be
	neededCoefficients := uint(len(bytes)-1)/7 + 1
	fmt.Println(bytes)
	fmt.Println(neededCoefficients)
	fmt.Println(coefficientCount)
	if neededCoefficients > coefficientCount {
		fmt.Printf("Secret needs more space to be encoded, need at least %d coefficients\n", neededCoefficients)
		return PolynomialField{}
	}

	coefficients := make([]uint64, coefficientCount)
	// make all the coefficients 1
	for i := uint(0); i < coefficientCount; i++ {
		coefficients[i] = 1
	}

	for i := uint(0); i < neededCoefficients-1; i++ {
		formatted := make([]byte, 8)
		copy(formatted, bytes[7*i:7*(i+1)])
		coefficients[i] = binary.LittleEndian.Uint64(formatted)
	}

	// there may be some bytes left over
	formatted := make([]byte, 8)
	copy(formatted, bytes[(neededCoefficients-1)*7:])
	coefficients[neededCoefficients-1] = binary.LittleEndian.Uint64(formatted)

	return PolynomialField{coefficients, prime}
}
