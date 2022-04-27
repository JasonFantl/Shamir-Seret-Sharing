package shamir

import "fmt"

type Share struct {
	X, Y uint64
}

type PolynomialField struct {
	coefficients []uint64
	prime        uint64
}

func (polynomial1 PolynomialField) add(polynomial2 PolynomialField) PolynomialField {
	if polynomial1.prime != polynomial2.prime {
		fmt.Println("cannot add polynomial fields of different primes")
		return polynomial1
	}
	length := maxInt(len(polynomial1.coefficients), len(polynomial2.coefficients))
	resultPolynomial := PolynomialField{make([]uint64, length), polynomial1.prime}
	for i := 0; i < length; i++ {
		if i < len(polynomial1.coefficients) {
			resultPolynomial.coefficients[i] = add64Mod(
				resultPolynomial.coefficients[i],
				polynomial1.coefficients[i],
				resultPolynomial.prime)
		}
		if i < len(polynomial2.coefficients) {
			resultPolynomial.coefficients[i] = add64Mod(
				resultPolynomial.coefficients[i],
				polynomial2.coefficients[i],
				resultPolynomial.prime)
		}
	}
	return resultPolynomial
}

func (polynomial1 PolynomialField) sub(polynomial2 PolynomialField) PolynomialField {
	if polynomial1.prime != polynomial2.prime {
		fmt.Println("cannot sub polynomial fields with different primes")
		return polynomial1
	}
	length := maxInt(len(polynomial1.coefficients), len(polynomial2.coefficients))
	resultPolynomial := PolynomialField{make([]uint64, length), polynomial1.prime}
	for i := 0; i < length; i++ {
		if i < len(polynomial1.coefficients) {
			resultPolynomial.coefficients[i] = add64Mod(
				resultPolynomial.coefficients[i],
				polynomial1.coefficients[i],
				resultPolynomial.prime)
		}
		if i < len(polynomial2.coefficients) {
			resultPolynomial.coefficients[i] = sub64Mod(
				resultPolynomial.coefficients[i],
				polynomial2.coefficients[i],
				resultPolynomial.prime)
		}
	}
	return resultPolynomial
}

func (polynomial1 PolynomialField) mult(polynomial2 PolynomialField) PolynomialField {
	if polynomial1.prime != polynomial2.prime {
		fmt.Println("cannot mult polynomial fields with different primes")
		return polynomial1
	}
	length := len(polynomial1.coefficients) + len(polynomial2.coefficients) - 1
	resultPolynomial := PolynomialField{make([]uint64, length), polynomial1.prime}
	for i := 0; i < len(polynomial1.coefficients); i++ {
		for j := 0; j < len(polynomial2.coefficients); j++ {
			resultPolynomial.coefficients[i+j] = add64Mod(
				resultPolynomial.coefficients[i+j],
				mult64Mod(
					polynomial1.coefficients[i],
					polynomial2.coefficients[j],
					resultPolynomial.prime),
				resultPolynomial.prime)
		}
	}
	return resultPolynomial
}

func (polynomial PolynomialField) scale(s uint64) PolynomialField {
	resultPolynomial := PolynomialField{make([]uint64, len(polynomial.coefficients)), polynomial.prime}
	for i := 0; i < len(polynomial.coefficients); i++ {
		resultPolynomial.coefficients[i] = mult64Mod(
			polynomial.coefficients[i],
			s,
			resultPolynomial.prime)
	}
	return resultPolynomial
}

func (polynomial PolynomialField) eval(x uint64) Share {
	// Horner's method
	degree := len(polynomial.coefficients) - 1
	y := polynomial.coefficients[degree]
	for i := degree - 1; i >= 0; i-- {
		y = add64Mod(mult64Mod(y, x, polynomial.prime), polynomial.coefficients[i], polynomial.prime)
	}
	return Share{x, y}
}

// make sure the number of points is = degree(polynomial) + 1,
// more will give an incorrect answer.
func lagrangeInterpolateEval(points []Share, x, prime uint64) Share {
	n := len(points)
	var sum uint64 = 0
	for i := 0; i < n; i++ {
		var product uint64 = 1
		for j := 0; j < n; j++ {
			if j == i {
				continue
			}
			product = mult64Mod(
				product,
				div64Mod(
					sub64Mod(x, points[j].X, prime),
					sub64Mod(points[i].X, points[j].X, prime),
					prime),
				prime)
		}

		sum = add64Mod(
			sum,
			mult64Mod(product, points[i].Y, prime),
			prime)
	}

	return Share{x, sum}
}

func lagrangeInterpolate(points []Share, prime uint64) PolynomialField {
	n := len(points)
	sum := PolynomialField{make([]uint64, 0), prime}
	for i := 0; i < n; i++ {
		product := PolynomialField{[]uint64{1}, prime}
		for j := 0; j < n; j++ {
			if j == i {
				continue
			}
			frac := PolynomialField{[]uint64{sub64Mod(0, points[j].X, prime), 1}, prime}
			frac = frac.scale(inverse64Mod(sub64Mod(points[i].X, points[j].X, prime), prime))
			product = product.mult(frac)
		}
		product = product.scale(points[i].Y)
		sum = sum.add(product)
	}

	return sum
}
