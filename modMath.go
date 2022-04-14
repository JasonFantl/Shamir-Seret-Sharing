package main

import (
	"fmt"
	"math"
)

// special operations to handle overflow and modulus

func add64Mod(left, right, prime uint64) uint64 {
	if left > math.MaxUint64-right { // overflow
		diff := (math.MaxUint64 % prime) + 1
		return add64Mod(left+right, diff, prime)
	}

	return (left + right) % prime
}

func sub64Mod(left, right, prime uint64) uint64 {
	// this makes the next steps easier
	left %= prime
	right %= prime

	if left < right { // overflow
		return prime - (right - left) // we know (right-left < prime), so this is safe
	}

	return left - right
}

// https://www.geeksforgeeks.org/how-to-avoid-overflow-in-modular-multiplication/
func mult64Mod(left, right, prime uint64) uint64 {
	left %= prime
	right %= prime

	var res uint64 = 0 // Initialize result
	for right > 0 {
		if right%2 == 1 {
			res = add64Mod(res, left, prime)
		}
		// Multiply 'left' with 2
		left = add64Mod(left, left, prime)
		right /= 2
	}
	return res
}

func div64Mod(left, right, prime uint64) uint64 {
	left %= prime
	right %= prime

	if right == 0 { // cannot divide by 0
		fmt.Println("cannot divide by zero")
		return 0
	}

	return mult64Mod(left, inverse(right, prime), prime)
}

// https://en.wikipedia.org/wiki/Extended_Euclidean_algorithm
func inverse(n, prime uint64) uint64 {
	n %= prime
	var old_r, r uint64 = n, prime
	var old_s, s uint64 = 1, 0
	// var old_t, t uint64 = 0, 1

	for r != 0 {
		quotient := old_r / r
		old_r, r = r, sub64Mod(old_r, quotient*r, prime)
		old_s, s = s, sub64Mod(old_s, quotient*s, prime)
		// old_t, t = t, sub64Mod(old_t, quotient*t, prime)
	}

	// Bézout coefficients: (old_s, old_t)
	// greatest common divisor: old_r
	// quotients by the gcd: (t, s)

	return old_s
}
