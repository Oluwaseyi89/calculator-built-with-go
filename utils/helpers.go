package utils

import (
	"fmt"
	"math"
	"strconv"
	"strings"
	"time"
)

// FormatNumber formats a float64 with appropriate precision
func FormatNumber(value float64) string {
	// Use scientific notation for very large or small numbers
	if math.Abs(value) >= 1e12 || (math.Abs(value) < 1e-6 && value != 0) {
		return fmt.Sprintf("%.6e", value)
	}

	// Format based on magnitude
	if math.Abs(value) >= 1000 {
		return fmt.Sprintf("%.2f", value)
	} else if math.Abs(value) >= 1 {
		return fmt.Sprintf("%.4f", value)
	} else {
		return fmt.Sprintf("%.6f", value)
	}
}

// DegreesToRadians converts degrees to radians
func DegreesToRadians(degrees float64) float64 {
	return degrees * math.Pi / 180
}

// RadiansToDegrees converts radians to degrees
func RadiansToDegrees(radians float64) float64 {
	return radians * 180 / math.Pi
}

// IsScientificNotation checks if a string is in scientific notation
func IsScientificNotation(s string) bool {
	return strings.Contains(strings.ToLower(s), "e")
}

// ParseScientificNotation parses scientific notation strings
func ParseScientificNotation(s string) (float64, error) {
	parts := strings.Split(strings.ToLower(s), "e")
	if len(parts) != 2 {
		return 0, fmt.Errorf("invalid scientific notation")
	}

	mantissa, err := strconv.ParseFloat(parts[0], 64)
	if err != nil {
		return 0, err
	}

	exponent, err := strconv.ParseFloat(parts[1], 64)
	if err != nil {
		return 0, err
	}

	return mantissa * math.Pow(10, exponent), nil
}

// FormatTime formats duration for history display
func FormatTime(t time.Time) string {
	return t.Format("15:04:05")
}

// FormatDuration formats a duration in human-readable form
func FormatDuration(d time.Duration) string {
	if d < time.Second {
		return fmt.Sprintf("%dms", d.Milliseconds())
	}
	return fmt.Sprintf("%.2fs", d.Seconds())
}

// RoundTo rounds a number to specified decimal places
func RoundTo(value float64, decimals int) float64 {
	multiplier := math.Pow(10, float64(decimals))
	return math.Round(value*multiplier) / multiplier
}

// TruncateString truncates a string to specified length with ellipsis
func TruncateString(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen-3] + "..."
}

// ParseComplexNumber parses complex numbers (for future enhancement)
func ParseComplexNumber(s string) (float64, float64, error) {
	s = strings.ReplaceAll(s, " ", "")

	if !strings.Contains(s, "i") {
		real, err := strconv.ParseFloat(s, 64)
		return real, 0, err
	}

	// Simple complex number parser: a+bi or a-bi
	s = strings.ReplaceAll(s, "i", "")
	parts := strings.Split(s, "+")
	if len(parts) == 2 {
		real, err1 := strconv.ParseFloat(parts[0], 64)
		imag, err2 := strconv.ParseFloat(parts[1], 64)

		if err1 != nil || err2 != nil {
			return 0, 0, fmt.Errorf("invalid complex number format")
		}
		return real, imag, nil
	}

	parts = strings.Split(s, "-")
	if len(parts) == 3 { // Case: a-b-c (negative real and imaginary)
		real, err1 := strconv.ParseFloat("-"+parts[1], 64)
		imag, err2 := strconv.ParseFloat("-"+parts[2], 64)

		if err1 != nil || err2 != nil {
			return 0, 0, fmt.Errorf("invalid complex number format")
		}
		return real, imag, nil
	} else if len(parts) == 2 && s[0] != '-' { // Case: a-b (positive real, negative imaginary)
		real, err1 := strconv.ParseFloat(parts[0], 64)
		imag, err2 := strconv.ParseFloat("-"+parts[1], 64)

		if err1 != nil || err2 != nil {
			return 0, 0, fmt.Errorf("invalid complex number format")
		}
		return real, imag, nil
	}

	return 0, 0, fmt.Errorf("invalid complex number format")
}

// FormatComplex formats complex numbers as strings
func FormatComplex(real, imag float64) string {
	if imag == 0 {
		return FormatNumber(real)
	}

	op := "+"
	if imag < 0 {
		op = "-"
		imag = -imag
	}

	return fmt.Sprintf("%s %s %si", FormatNumber(real), op, FormatNumber(imag))
}

// GCD computes greatest common divisor (useful for fraction reduction)
func GCD(a, b int) int {
	for b != 0 {
		a, b = b, a%b
	}
	return a
}

// LCM computes least common multiple
func LCM(a, b int) int {
	return a * b / GCD(a, b)
}

// ConvertToFraction converts decimal to fraction approximation
func ConvertToFraction(decimal float64, maxDenominator int) (int, int) {
	// Using continued fraction approximation
	epsilon := 1.0 / float64(maxDenominator*maxDenominator)

	numerator := 1
	denominator := 1
	bestNumerator := 1
	bestDenominator := 1
	bestError := math.Abs(decimal - float64(bestNumerator)/float64(bestDenominator))

	for denominator <= maxDenominator {
		value := float64(numerator) / float64(denominator)

		if math.Abs(decimal-value) < bestError {
			bestError = math.Abs(decimal - value)
			bestNumerator = numerator
			bestDenominator = denominator
		}

		if value < decimal {
			numerator++
		} else {
			denominator++
		}

		if bestError < epsilon {
			break
		}
	}

	return bestNumerator, bestDenominator
}

// FactorialTable provides pre-computed factorial values for optimization
var FactorialTable = []float64{
	1,       // 0!
	1,       // 1!
	2,       // 2!
	6,       // 3!
	24,      // 4!
	120,     // 5!
	720,     // 6!
	5040,    // 7!
	40320,   // 8!
	362880,  // 9!
	3628800, // 10!
}

// GetFactorial returns factorial with caching
func GetFactorial(n int) float64 {
	if n < 0 {
		return math.NaN()
	}

	if n < len(FactorialTable) {
		return FactorialTable[n]
	}

	// Compute recursively
	result := FactorialTable[len(FactorialTable)-1]
	for i := len(FactorialTable); i <= n; i++ {
		result *= float64(i)
	}

	return result
}

// IsPrime checks if a number is prime
func IsPrime(n int) bool {
	if n < 2 {
		return false
	}
	if n == 2 || n == 3 {
		return true
	}
	if n%2 == 0 || n%3 == 0 {
		return false
	}

	for i := 5; i*i <= n; i += 6 {
		if n%i == 0 || n%(i+2) == 0 {
			return false
		}
	}
	return true
}
