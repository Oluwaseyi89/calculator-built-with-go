package calculator

import (
	"fmt"
	"math"
)

// Trigonometric functions (radians)
func Sin(x float64) float64 {
	return math.Sin(x)
}

func Cos(x float64) float64 {
	return math.Cos(x)
}

func Tan(x float64) float64 {
	return math.Tan(x)
}

// Inverse trigonometric
func Asin(x float64) (float64, error) {
	if x < -1 || x > 1 {
		return 0, fmt.Errorf("asin input must be between -1 and 1")
	}
	return math.Asin(x), nil
}

func Acos(x float64) (float64, error) {
	if x < -1 || x > 1 {
		return 0, fmt.Errorf("acos input must be between -1 and 1")
	}
	return math.Acos(x), nil
}

func Atan(x float64) float64 {
	return math.Atan(x)
}

// Hyperbolic functions
func Sinh(x float64) float64 {
	return math.Sinh(x)
}

func Cosh(x float64) float64 {
	return math.Cosh(x)
}

func Tanh(x float64) float64 {
	return math.Tanh(x)
}

// Roots and logarithms
func Sqrt(x float64) (float64, error) {
	if x < 0 {
		return 0, fmt.Errorf("square root of negative number")
	}
	return math.Sqrt(x), nil
}

func Cbrt(x float64) float64 {
	return math.Cbrt(x)
}

func Log(x float64) (float64, error) {
	if x <= 0 {
		return 0, fmt.Errorf("logarithm undefined for non-positive numbers")
	}
	return math.Log(x), nil
}

func Log10(x float64) (float64, error) {
	if x <= 0 {
		return 0, fmt.Errorf("log10 undefined for non-positive numbers")
	}
	return math.Log10(x), nil
}

// Exponential and absolute
func Exp(x float64) float64 {
	return math.Exp(x)
}

func Abs(x float64) float64 {
	return math.Abs(x)
}

// Constants
const (
	Pi = math.Pi
	E  = math.E
)
