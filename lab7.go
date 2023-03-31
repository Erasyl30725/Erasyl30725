package main

import (
	"fmt"
	"math"
)

func prime(num int) bool {
	if num < 2 {
		return false
	}
	for i := 2; i <= int(math.Sqrt(float64(num))); i++ {
		if num%i == 0 {
			return false
		}
	}
	return true
}

func main() {
	var n int
	fmt.Scan(&n)
	primes := make([]int, n)
	var result []int
	for i := 0; i < n; i++ {
		fmt.Scan(&primes[i])
	}
	for i := 0; i < n; i++ {
		if prime(primes[i]) {
			result = append(result, primes[i])
		}
	}
	fmt.Println(result)
}