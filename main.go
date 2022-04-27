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

	// distribute these shares to people, then when they gather together they will have this array
	shares := make([]shamir.Share, requiredShares+2)
	for i := 0; i < len(shares); i++ {
		shares[i] = ss.GenerateShare()
	}

	decoded, err := shamir.DecodeSecret(shares, requiredShares)
	if err != nil {
		fmt.Println(err)
		return
	}

	// note that the secret is in decoded, but will be padded, so decoded is not exactly the secret
	fmt.Println(decoded)
	fmt.Println(string(decoded))

	addSecretly()

}

func addSecretly() {
	// can we add numbers without sharing the numbers?
	// we will store the secrect in the first byte

	requiredShares := 3

	// person 1 will generate final point at x = 1, so distributes shares at 2 and 3
	s1 := []byte{byte(7)}
	p1, _ := shamir.NewSecretSharing(s1, requiredShares)
	share1a := p1.GenerateShareAt(2)
	share1b := p1.GenerateShareAt(3)

	s2 := []byte{byte(9)}
	p2, _ := shamir.NewSecretSharing(s2, requiredShares)
	share2a := p2.GenerateShareAt(1)
	share2b := p2.GenerateShareAt(3)

	s3 := []byte{byte(1)}
	p3, _ := shamir.NewSecretSharing(s3, requiredShares)
	share3a := p3.GenerateShareAt(1)
	share3b := p3.GenerateShareAt(2)

	// the generated shares (share(number)(letter)) should be sent securely to each other
	// person 1 generates their share of the final sum
	shareSum1 := shamir.Share{X: 1, Y: p1.GenerateShareAt(1).Y + share2a.Y + share3a.Y}

	shareSum2 := shamir.Share{X: 2, Y: share1a.Y + p2.GenerateShareAt(2).Y + share3b.Y}

	shareSum3 := shamir.Share{X: 3, Y: share1b.Y + share2b.Y + p3.GenerateShareAt(3).Y}

	// the shares should then be sent securely to each other
	decoded, err := shamir.DecodeSecret([]shamir.Share{shareSum1, shareSum2, shareSum3}, requiredShares)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(decoded)
	fmt.Println(int(decoded[0]))
}
