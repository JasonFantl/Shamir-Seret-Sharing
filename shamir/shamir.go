package shamir

import (
	"encoding/binary"
	"fmt"
)

// this large prime lets us encode information into the first 63 bits of uint64
var prime uint64 = 18446744073709551557

type SecretSharing struct {
	polynomial PolynomialField
	counter    uint
}

func NewSecretSharing(secret []byte, requiredShareCount int) (SecretSharing, error) {
	poly, err := encodeSecret(secret, uint(requiredShareCount))
	if err != nil {
		return SecretSharing{}, err
	}

	return SecretSharing{
		polynomial: poly,
		counter:    1,
	}, nil
}

func (s *SecretSharing) GenerateShare() Point {
	share := s.polynomial.eval(uint64(s.counter))
	s.counter++
	return share
}

func encodeSecret(bytes []byte, requiredShareCount uint) (PolynomialField, error) {

	// well just use every 7 bytes since it will fit in our encoding, although its not as storage efficient as it could be
	neededCoefficients := uint(len(bytes)-1)/7 + 1
	if neededCoefficients > requiredShareCount {
		return PolynomialField{}, fmt.Errorf("Secret needs more space to be encoded, need at least %d coefficients\n", neededCoefficients)
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

	return PolynomialField{coefficients, prime}, nil
}

func DecodeSecret(shares []Point, degree int) ([]byte, error) {
	// make sure the provided information is valid
	if len(shares) < degree {
		return nil, fmt.Errorf("not enough shares to decode")
	}
	for i, share1 := range shares {
		for j, share2 := range shares {
			if i == j {
				continue
			}
			if share1.x == share2.x {
				return nil, fmt.Errorf("shares must be distinct")
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
	return bytes, nil
}
