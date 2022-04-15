package main

import (
	"encoding/binary"
	"fmt"
)

// this large prime lets us encode information into the first 63 bits of uint64
var prime uint64 = 18446744073709551557

func encodeSecret(bytes []byte, requiredShareCount uint) PolynomialField {

	// well just use every 7 bytes since it will fit in our encoding, although its not as storage efficient as it could be
	neededCoefficients := uint(len(bytes)-1)/7 + 1
	if neededCoefficients > requiredShareCount {
		fmt.Printf("Secret needs more space to be encoded, need at least %d coefficients\n", neededCoefficients)
		return PolynomialField{}
	}

	coefficients := make([]uint64, requiredShareCount)
	// first fill the coefficients with non-zero values, then we will overwrite with our own
	for i := uint(0); i < requiredShareCount; i++ {
		coefficients[i] = '#'
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

func decodeSecret(shares []Point, degree uint) []byte {
	// make sure the provided information is valid
	if uint(len(shares)) < degree {
		fmt.Println("Not enough shares to decode")
		return []byte{}
	}
	for i, share1 := range shares {
		for j, share2 := range shares {
			if i == j {
				continue
			}
			if share1.x == share2.x {
				fmt.Println("Shares must be distinct")
				return []byte{}
			}
		}
	}

	decodedPolynomial := lagrangeInterpolate(shares[:degree], prime)

	bytes := make([]byte, 0)

	for _, coefficient := range decodedPolynomial.coefficients {
		b := make([]byte, 8)
		binary.LittleEndian.PutUint64(b, coefficient)
		bytes = append(bytes, b[:7]...)
	}
	return bytes
}
