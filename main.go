package main

import (
	"fmt"
)

func main() {

	// var prime uint64 = math.MaxUint64 - 58 // https://primes.utm.edu/lists/2small/0bit.html
	// fmt.Println(prime)
	var prime uint64 = 67

	// poly := PolynomialField{[]uint64{5, 55, 10}, prime}
	// poly2 := PolynomialField{[]uint64{50, 2, 62}, prime}

	fmt.Println(inverse(60, prime))
	// https://www.desmos.com/calculator/2fn27pdvuy
	points := []Point{{0, 5}, {1, 6}, {3, 10}, {7, 88}}
	interpolated := lagrangeInterpolate(points, prime)
	fmt.Println(interpolated)
	fmt.Println(lagrangeInterpolate([]Point{interpolated.eval(5), interpolated.eval(10), interpolated.eval(30), interpolated.eval(90)}, interpolated.prime))

	// why does the above work, but not the below? is it because of overflow?
	encoded := encodeSecret([]byte("hi"), 3)
	fmt.Println(encoded)
	shares := []Point{encoded.eval(0), encoded.eval(1), encoded.eval(10)}
	decoded := lagrangeInterpolate(shares, encoded.prime)
	fmt.Println(decoded)

}
