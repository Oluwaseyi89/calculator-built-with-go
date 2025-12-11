package utils

import (
	"fmt"
	"math"
	"regexp"
	"strings"
)

var (
	// ValidExpressionRegex validates basic calculator expressions
	ValidExpressionRegex = regexp.MustCompile(`^[0-9+\-*/().^%!a-zπe\s]+$`)

	// ValidFunctionRegex validates function calls
	ValidFunctionRegex = regexp.MustCompile(`^[a-z]+\([^)]+\)$`)

	// ValidNumberRegex validates numeric strings
	ValidNumberRegex = regexp.MustCompile(`^-?\d*\.?\d+(?:[eE][+-]?\d+)?$`)

	// ValidIdentifierRegex validates variable/constant names
	ValidIdentifierRegex = regexp.MustCompile(`^[a-zA-Z_][a-zA-Z0-9_]*$`)
)

// ValidateExpression validates a mathematical expression
func ValidateExpression(expr string) error {
	if expr == "" {
		return fmt.Errorf("empty expression")
	}

	// Check for balanced parentheses
	if !HasBalancedParentheses(expr) {
		return fmt.Errorf("unbalanced parentheses")
	}

	// Check for consecutive operators
	if HasConsecutiveOperators(expr) {
		return fmt.Errorf("consecutive operators")
	}

	// Check for invalid characters
	if !ValidExpressionRegex.MatchString(expr) {
		return fmt.Errorf("invalid characters in expression")
	}

	// Check for division by zero patterns
	if strings.Contains(expr, "/0") && !strings.Contains(expr, "/0.") {
		// Allow /0. as in 5/0.5
		if !strings.Contains(expr, "/0.") {
			return fmt.Errorf("potential division by zero")
		}
	}

	// Check for invalid function calls
	if err := ValidateFunctionCalls(expr); err != nil {
		return err
	}

	// Check for valid number formats
	if err := ValidateNumbers(expr); err != nil {
		return err
	}

	return nil
}

// HasBalancedParentheses checks if parentheses are balanced
func HasBalancedParentheses(expr string) bool {
	balance := 0
	for _, ch := range expr {
		switch ch {
		case '(':
			balance++
		case ')':
			balance--
			if balance < 0 {
				return false
			}
		}
	}
	return balance == 0
}

// HasConsecutiveOperators checks for invalid operator sequences
func HasConsecutiveOperators(expr string) bool {
	operators := "+-*/^%"
	prev := ' '

	for _, ch := range expr {
		if strings.ContainsRune(operators, prev) && strings.ContainsRune(operators, ch) {
			// Allow +- or -- for signed numbers
			if !((prev == '+' || prev == '-') && (ch == '+' || ch == '-')) {
				return true
			}
		}
		prev = ch
	}
	return false
}

// ValidateFunctionCalls validates function syntax
func ValidateFunctionCalls(expr string) error {
	// Find all function calls
	funcPattern := regexp.MustCompile(`([a-z]+)\(`)
	matches := funcPattern.FindAllStringSubmatch(expr, -1)

	validFunctions := map[string]bool{
		"sin": true, "cos": true, "tan": true,
		"asin": true, "acos": true, "atan": true,
		"sinh": true, "cosh": true, "tanh": true,
		"sqrt": true, "cbrt": true,
		"log": true, "log10": true,
		"exp": true, "abs": true,
	}

	for _, match := range matches {
		funcName := match[1]
		if !validFunctions[funcName] {
			return fmt.Errorf("unknown function: %s", funcName)
		}
	}

	// Check for missing closing parentheses
	if strings.Count(expr, "(") != strings.Count(expr, ")") {
		return fmt.Errorf("mismatched parentheses")
	}

	return nil
}

// ValidateNumbers validates numeric strings in expression
func ValidateNumbers(expr string) error {
	// Extract numbers from expression
	numberPattern := regexp.MustCompile(`-?\d*\.?\d+(?:[eE][+-]?\d+)?`)
	numbers := numberPattern.FindAllString(expr, -1)

	for _, numStr := range numbers {
		if !ValidNumberRegex.MatchString(numStr) {
			return fmt.Errorf("invalid number format: %s", numStr)
		}

		// Check for multiple decimal points
		if strings.Count(numStr, ".") > 1 {
			return fmt.Errorf("invalid number with multiple decimal points: %s", numStr)
		}

		// Validate scientific notation
		if strings.ContainsAny(numStr, "eE") {
			parts := strings.Split(strings.ToLower(numStr), "e")
			if len(parts) != 2 {
				return fmt.Errorf("invalid scientific notation: %s", numStr)
			}

			// Parse to validate
			if _, err := ParseScientificNotation(numStr); err != nil {
				return fmt.Errorf("invalid scientific notation: %s", numStr)
			}
		}
	}

	return nil
}

// ValidateAngle checks if angle is valid for trigonometric functions
func ValidateAngle(angle float64, mode string) error {
	if mode == "deg" {
		if math.IsInf(angle, 0) || math.IsNaN(angle) {
			return fmt.Errorf("invalid angle value")
		}
		// No specific range check for degrees
	} else if mode == "rad" {
		if math.IsInf(angle, 0) || math.IsNaN(angle) {
			return fmt.Errorf("invalid angle value")
		}
	}
	return nil
}

// ValidateLogArgument validates argument for logarithmic functions
func ValidateLogArgument(x float64) error {
	if x <= 0 {
		return fmt.Errorf("logarithm undefined for non-positive numbers")
	}
	if math.IsInf(x, 0) || math.IsNaN(x) {
		return fmt.Errorf("invalid argument for logarithm")
	}
	return nil
}

// ValidateSqrtArgument validates argument for square root
func ValidateSqrtArgument(x float64) error {
	if x < 0 {
		return fmt.Errorf("square root undefined for negative numbers")
	}
	if math.IsInf(x, 0) || math.IsNaN(x) {
		return fmt.Errorf("invalid argument for square root")
	}
	return nil
}

// ValidateFactorialArgument validates argument for factorial
func ValidateFactorialArgument(x float64) error {
	if x < 0 || x != math.Trunc(x) {
		return fmt.Errorf("factorial undefined for non-integer or negative numbers")
	}
	if x > 170 { // 170! is close to max float64
		return fmt.Errorf("factorial too large for calculation")
	}
	return nil
}

// ValidateDivision validates division operation
func ValidateDivision(numerator, denominator float64) error {
	if denominator == 0 {
		return fmt.Errorf("division by zero")
	}
	if math.IsInf(numerator, 0) || math.IsInf(denominator, 0) {
		return fmt.Errorf("invalid values for division")
	}
	if math.IsNaN(numerator) || math.IsNaN(denominator) {
		return fmt.Errorf("NaN values for division")
	}
	return nil
}

// ValidatePower validates power operation
func ValidatePower(base, exponent float64) error {
	if base == 0 && exponent <= 0 {
		return fmt.Errorf("0^0 or 0^(negative) is undefined")
	}
	if base < 0 && exponent != math.Trunc(exponent) {
		return fmt.Errorf("negative base with non-integer exponent")
	}
	return nil
}

// IsValidVariableName checks if a string is a valid variable name
func IsValidVariableName(name string) bool {
	return ValidIdentifierRegex.MatchString(name) && !IsReservedKeyword(name)
}

// IsReservedKeyword checks if a string is a reserved keyword
func IsReservedKeyword(s string) bool {
	reserved := map[string]bool{
		"pi": true, "e": true, "ans": true,
		"sin": true, "cos": true, "tan": true,
		"asin": true, "acos": true, "atan": true,
		"sinh": true, "cosh": true, "tanh": true,
		"sqrt": true, "cbrt": true,
		"log": true, "log10": true,
		"exp": true, "abs": true,
		"exit": true, "quit": true,
		"help": true, "clear": true,
		"mem": true, "history": true,
	}

	return reserved[strings.ToLower(s)]
}

// ValidateRange checks if a value is within specified range
func ValidateRange(value, min, max float64, inclusive bool) error {
	if inclusive {
		if value < min || value > max {
			return fmt.Errorf("value %.6g out of range [%.6g, %.6g]", value, min, max)
		}
	} else {
		if value <= min || value >= max {
			return fmt.Errorf("value %.6g out of range (%.6g, %.6g)", value, min, max)
		}
	}
	return nil
}

// ValidateProbability checks if a value is a valid probability (0-1)
func ValidateProbability(p float64) error {
	if p < 0 || p > 1 {
		return fmt.Errorf("probability must be between 0 and 1")
	}
	if math.IsInf(p, 0) || math.IsNaN(p) {
		return fmt.Errorf("invalid probability value")
	}
	return nil
}

// ValidateMatrixDimensions validates matrix dimensions
func ValidateMatrixDimensions(rows, cols int) error {
	if rows <= 0 || cols <= 0 {
		return fmt.Errorf("matrix dimensions must be positive")
	}
	if rows > 1000 || cols > 1000 {
		return fmt.Errorf("matrix dimensions too large")
	}
	return nil
}

// ValidateComplexNumber validates complex number components
func ValidateComplexNumber(real, imag float64) error {
	if math.IsInf(real, 0) || math.IsInf(imag, 0) {
		return fmt.Errorf("complex number components cannot be infinite")
	}
	if math.IsNaN(real) || math.IsNaN(imag) {
		return fmt.Errorf("complex number components cannot be NaN")
	}
	return nil
}
