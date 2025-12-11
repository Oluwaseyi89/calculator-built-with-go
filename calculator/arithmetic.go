package calculator

import (
	"fmt"
	"math"
)

// Basic arithmetic operations
func Add(a, b float64) float64 {
	return a + b
}

func Subtract(a, b float64) float64 {
	return a - b
}

func Multiply(a, b float64) float64 {
	return a * b
}

func Divide(a, b float64) (float64, error) {
	if b == 0 {
		return 0, fmt.Errorf("division by zero")
	}
	return a / b, nil
}

func Power(base, exponent float64) float64 {
	return math.Pow(base, exponent)
}

func Modulus(a, b float64) (float64, error) {
	if b == 0 {
		return 0, fmt.Errorf("modulus by zero")
	}
	return math.Mod(a, b), nil
}

func Percentage(value, percent float64) float64 {
	return value * percent / 100
}

func Factorial(n float64) (float64, error) {
	if n < 0 || n != math.Trunc(n) {
		return 0, fmt.Errorf("factorial undefined for non-integer or negative numbers")
	}

	result := 1.0
	for i := 2.0; i <= n; i++ {
		result *= i
	}
	return result, nil
}
