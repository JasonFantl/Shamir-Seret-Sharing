# Shamir-Seret-Sharing

This project is a simple implementation of Shamir's Secret Sharing Algorithm [https://en.wikipedia.org/wiki/Shamir%27s_Secret_Sharing].

## Running

You must have Go installed, then you can simply run `go run main.go`, or run `go build` followed by executing the generated executable.

## Notes

`main.go` contains in its main function the process of encoding a secret into a polynomial, then generating unique shares to distribute, then finally recreating the polynomial using the shares, allowing you to access the original secret.

There is also an example using this algorithm to obtain the sum of some values without ever actually knowing the individual values. In essence this is done by summing each persons polynomial, which will contain the sum of each persons value. The only information actually shared are points on each persons polynomial, but not enough points to recreate their polynomial.

