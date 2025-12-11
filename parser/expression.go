package parser

import (
	"fmt"
	"math"
	"regexp"
	"strconv"
	"strings"

	"github.com/Oluwaseyi89/calculator-built-with-go/calculator"
	"github.com/Oluwaseyi89/calculator-built-with-go/utils"
)

type Parser struct {
	calc       *calculator.Calculator
	angleMode  string // "deg" or "rad"
	variables  map[string]float64
	constants  map[string]float64
	ops        map[string]int
	funcRegex  *regexp.Regexp
	unaryMinus string // Marker for unary minus
}

func NewParser(calc *calculator.Calculator) *Parser {
	return &Parser{
		calc:      calc,
		angleMode: "rad",
		variables: make(map[string]float64),
		constants: map[string]float64{
			"pi":  math.Pi,
			"e":   math.E,
			"ans": 0,
		},
		ops: map[string]int{
			"+": 1, "-": 1,
			"*": 2, "/": 2, "%": 2,
			"^": 3,
		},
		funcRegex:  regexp.MustCompile(`([a-zA-Z_][a-zA-Z0-9_]*)\(([^)]*)\)`),
		unaryMinus: "~", // Use ~ to represent unary minus
	}
}

func (p *Parser) SetAngleMode(mode string) {
	if mode == "deg" || mode == "rad" {
		p.angleMode = mode
	}
}

func (p *Parser) SetVariable(name string, value float64) error {
	if !utils.IsValidVariableName(name) {
		return fmt.Errorf("invalid variable name: %s", name)
	}
	p.variables[name] = value
	return nil
}

func (p *Parser) GetVariable(name string) (float64, bool) {
	val, exists := p.variables[name]
	return val, exists
}

func (p *Parser) ClearVariables() {
	p.variables = make(map[string]float64)
}

func (p *Parser) EvaluateExpression(expr string) (float64, error) {
	if expr == "" {
		return 0, fmt.Errorf("empty expression")
	}

	// Validate expression first
	if err := utils.ValidateExpression(expr); err != nil {
		return 0, fmt.Errorf("invalid expression: %v", err)
	}

	// Update ans constant with last result
	p.constants["ans"] = p.calc.GetLastResult()

	// Preprocess expression
	processed, err := p.preprocess(expr)
	if err != nil {
		return 0, err
	}

	// Evaluate
	result, err := p.evaluate(processed)
	if err != nil {
		return 0, err
	}

	// Update calculator history
	p.calc.SetLastResult(result)
	p.calc.AddToHistory(expr, result)

	return result, nil
}

func (p *Parser) preprocess(expr string) (string, error) {
	// Convert to lowercase and remove spaces
	expr = strings.ToLower(strings.ReplaceAll(expr, " ", ""))

	// Replace constants
	for name, value := range p.constants {
		// Use word boundaries to avoid partial matches
		pattern := regexp.MustCompile(`\b` + regexp.QuoteMeta(name) + `\b`)
		expr = pattern.ReplaceAllString(expr, fmt.Sprintf("%.10f", value))
	}

	// Replace variables
	for name, value := range p.variables {
		pattern := regexp.MustCompile(`\b` + regexp.QuoteMeta(name) + `\b`)
		expr = pattern.ReplaceAllString(expr, fmt.Sprintf("%.10f", value))
	}

	// Handle functions
	expr, err := p.processFunctions(expr)
	if err != nil {
		return "", err
	}

	return expr, nil
}

func (p *Parser) processFunctions(expr string) (string, error) {
	for {
		matches := p.funcRegex.FindStringSubmatch(expr)
		if matches == nil {
			break
		}

		funcName := matches[1]
		argsStr := matches[2]

		// Handle built-in functions
		result, err := p.evaluateFunction(funcName, argsStr)
		if err != nil {
			return "", err
		}

		// Replace function call with result
		expr = strings.Replace(expr, matches[0], fmt.Sprintf("%.10f", result), 1)
	}

	return expr, nil
}

func (p *Parser) evaluateFunction(name, argsStr string) (float64, error) {
	// Parse arguments
	args, err := p.parseArguments(argsStr)
	if err != nil {
		return 0, err
	}

	// Evaluate arguments
	argValues := make([]float64, len(args))
	for i, arg := range args {
		val, err := p.EvaluateExpression(arg)
		if err != nil {
			return 0, fmt.Errorf("error in function argument: %v", err)
		}
		argValues[i] = val
	}

	// Dispatch to appropriate function handler
	switch name {
	case "sin", "cos", "tan", "asin", "acos", "atan",
		"sinh", "cosh", "tanh":
		if len(argValues) != 1 {
			return 0, fmt.Errorf("function %s expects 1 argument", name)
		}
		return p.evaluateTrigFunction(name, argValues[0])

	case "sqrt", "cbrt", "log", "log10", "exp", "abs":
		if len(argValues) != 1 {
			return 0, fmt.Errorf("function %s expects 1 argument", name)
		}
		return p.evaluateMathFunction(name, argValues[0])

	case "pow":
		if len(argValues) != 2 {
			return 0, fmt.Errorf("pow expects 2 arguments")
		}
		return calculator.Power(argValues[0], argValues[1]), nil

	case "min":
		if len(argValues) < 1 {
			return 0, fmt.Errorf("min expects at least 1 argument")
		}
		return p.min(argValues), nil

	case "max":
		if len(argValues) < 1 {
			return 0, fmt.Errorf("max expects at least 1 argument")
		}
		return p.max(argValues), nil

	case "round", "floor", "ceil":
		if len(argValues) == 1 {
			return p.evaluateRounding(name, argValues[0], 0)
		} else if len(argValues) == 2 {
			return p.evaluateRounding(name, argValues[0], int(argValues[1]))
		}
		return 0, fmt.Errorf("%s expects 1 or 2 arguments", name)

	default:
		return 0, fmt.Errorf("unknown function: %s", name)
	}
}

func (p *Parser) parseArguments(argsStr string) ([]string, error) {
	if argsStr == "" {
		return []string{}, nil
	}

	var args []string
	var current strings.Builder
	parenDepth := 0

	for _, ch := range argsStr {
		switch ch {
		case '(':
			parenDepth++
			current.WriteRune(ch)
		case ')':
			parenDepth--
			current.WriteRune(ch)
		case ',':
			if parenDepth == 0 {
				args = append(args, strings.TrimSpace(current.String()))
				current.Reset()
			} else {
				current.WriteRune(ch)
			}
		default:
			current.WriteRune(ch)
		}
	}

	if current.Len() > 0 {
		args = append(args, strings.TrimSpace(current.String()))
	}

	return args, nil
}

func (p *Parser) evaluateTrigFunction(name string, arg float64) (float64, error) {
	// Convert to radians if in degree mode
	if p.angleMode == "deg" && !strings.HasSuffix(name, "h") { // not hyperbolic
		arg = utils.DegreesToRadians(arg)
	}

	switch name {
	case "sin":
		return calculator.Sin(arg), nil
	case "cos":
		return calculator.Cos(arg), nil
	case "tan":
		return calculator.Tan(arg), nil
	case "asin":
		return calculator.Asin(arg)
	case "acos":
		return calculator.Acos(arg)
	case "atan":
		return calculator.Atan(arg), nil
	case "sinh":
		return calculator.Sinh(arg), nil
	case "cosh":
		return calculator.Cosh(arg), nil
	case "tanh":
		return calculator.Tanh(arg), nil
	default:
		return 0, fmt.Errorf("unknown trigonometric function: %s", name)
	}
}

func (p *Parser) evaluateMathFunction(name string, arg float64) (float64, error) {
	switch name {
	case "sqrt":
		return calculator.Sqrt(arg)
	case "cbrt":
		return calculator.Cbrt(arg), nil
	case "log":
		return calculator.Log(arg)
	case "log10":
		return calculator.Log10(arg)
	case "exp":
		return calculator.Exp(arg), nil
	case "abs":
		return calculator.Abs(arg), nil
	default:
		return 0, fmt.Errorf("unknown math function: %s", name)
	}
}

func (p *Parser) evaluateRounding(name string, value float64, precision int) (float64, error) {
	multiplier := math.Pow(10, float64(precision))

	switch name {
	case "round":
		return math.Round(value*multiplier) / multiplier, nil
	case "floor":
		return math.Floor(value*multiplier) / multiplier, nil
	case "ceil":
		return math.Ceil(value*multiplier) / multiplier, nil
	default:
		return value, nil
	}
}

func (p *Parser) min(values []float64) float64 {
	if len(values) == 0 {
		return 0
	}
	minVal := values[0]
	for _, v := range values[1:] {
		if v < minVal {
			minVal = v
		}
	}
	return minVal
}

func (p *Parser) max(values []float64) float64 {
	if len(values) == 0 {
		return 0
	}
	maxVal := values[0]
	for _, v := range values[1:] {
		if v > maxVal {
			maxVal = v
		}
	}
	return maxVal
}

func (p *Parser) evaluate(expr string) (float64, error) {
	// Handle factorial operator separately
	if strings.Contains(expr, "!") {
		return p.evaluateWithFactorial(expr)
	}

	// Tokenize
	tokens := p.tokenize(expr)

	// Handle unary operators
	tokens = p.processUnaryOperators(tokens)

	// Convert to RPN
	rpn, err := p.shuntingYard(tokens)
	if err != nil {
		return 0, err
	}

	// Evaluate RPN
	return p.evaluateRPN(rpn)
}

func (p *Parser) evaluateWithFactorial(expr string) (float64, error) {
	// Find the last factorial operator
	lastFactorial := strings.LastIndex(expr, "!")
	if lastFactorial == -1 {
		return 0, fmt.Errorf("no factorial operator found")
	}

	// Get the base expression (everything before the factorial)
	baseExpr := expr[:lastFactorial]
	factorialCount := len(expr) - lastFactorial // Count of ! characters

	// Evaluate the base expression
	base, err := p.evaluate(baseExpr)
	if err != nil {
		return 0, err
	}

	// Apply factorial repeatedly
	result := base
	for i := 0; i < factorialCount; i++ {
		result, err = calculator.Factorial(result)
		if err != nil {
			return 0, err
		}
	}

	return result, nil
}

func (p *Parser) tokenize(expr string) []string {
	var tokens []string
	var current strings.Builder

	for i, ch := range expr {
		if p.isDigit(ch) || ch == '.' || (ch == 'e' && i > 0 && p.isDigit(rune(expr[i-1]))) {
			current.WriteRune(ch)
		} else {
			if current.Len() > 0 {
				tokens = append(tokens, current.String())
				current.Reset()
			}
			if ch != ' ' {
				tokens = append(tokens, string(ch))
			}
		}
	}

	if current.Len() > 0 {
		tokens = append(tokens, current.String())
	}

	return tokens
}

func (p *Parser) processUnaryOperators(tokens []string) []string {
	result := make([]string, 0, len(tokens))

	for i, token := range tokens {
		if token == "-" && p.isUnaryMinus(tokens, i) {
			result = append(result, p.unaryMinus)
		} else if token == "+" && p.isUnaryPlus(tokens, i) {
			// Unary plus can be ignored
			continue
		} else {
			result = append(result, token)
		}
	}

	return result
}

func (p *Parser) isUnaryMinus(tokens []string, i int) bool {
	if i == 0 {
		return true
	}

	prev := tokens[i-1]
	return !p.isNumber(prev) && prev != ")" && prev != p.unaryMinus
}

func (p *Parser) isUnaryPlus(tokens []string, i int) bool {
	if i == 0 {
		return true
	}

	prev := tokens[i-1]
	return !p.isNumber(prev) && prev != ")" && prev != p.unaryMinus
}

func (p *Parser) shuntingYard(tokens []string) ([]string, error) {
	var output []string
	var stack []string

	for _, token := range tokens {
		if p.isNumber(token) {
			output = append(output, token)
		} else if token == "(" {
			stack = append(stack, token)
		} else if token == ")" {
			for len(stack) > 0 && stack[len(stack)-1] != "(" {
				output = append(output, stack[len(stack)-1])
				stack = stack[:len(stack)-1]
			}
			if len(stack) == 0 {
				return nil, fmt.Errorf("mismatched parentheses")
			}
			stack = stack[:len(stack)-1]
		} else if p.isOperator(token) {
			for len(stack) > 0 &&
				stack[len(stack)-1] != "(" &&
				p.precedence(stack[len(stack)-1]) >= p.precedence(token) {
				output = append(output, stack[len(stack)-1])
				stack = stack[:len(stack)-1]
			}
			stack = append(stack, token)
		} else {
			return nil, fmt.Errorf("invalid token: %s", token)
		}
	}

	for len(stack) > 0 {
		if stack[len(stack)-1] == "(" {
			return nil, fmt.Errorf("mismatched parentheses")
		}
		output = append(output, stack[len(stack)-1])
		stack = stack[:len(stack)-1]
	}

	return output, nil
}

func (p *Parser) evaluateRPN(rpn []string) (float64, error) {
	var stack []float64

	for _, token := range rpn {
		if p.isNumber(token) {
			val, err := strconv.ParseFloat(token, 64)
			if err != nil {
				return 0, fmt.Errorf("invalid number: %s", token)
			}
			stack = append(stack, val)
		} else if token == p.unaryMinus {
			if len(stack) < 1 {
				return 0, fmt.Errorf("insufficient operands for unary minus")
			}
			val := stack[len(stack)-1]
			stack[len(stack)-1] = -val
		} else if p.isOperator(token) {
			if len(stack) < 2 {
				return 0, fmt.Errorf("insufficient operands for operator %s", token)
			}

			b := stack[len(stack)-1]
			a := stack[len(stack)-2]
			stack = stack[:len(stack)-2]

			var result float64
			var err error

			switch token {
			case "+":
				result = calculator.Add(a, b)
			case "-":
				result = calculator.Subtract(a, b)
			case "*":
				result = calculator.Multiply(a, b)
			case "/":
				result, err = calculator.Divide(a, b)
			case "^":
				result = calculator.Power(a, b)
			case "%":
				result, err = calculator.Modulus(a, b)
			default:
				return 0, fmt.Errorf("unknown operator: %s", token)
			}

			if err != nil {
				return 0, err
			}

			stack = append(stack, result)
		} else {
			return 0, fmt.Errorf("unknown token in RPN: %s", token)
		}
	}

	if len(stack) != 1 {
		return 0, fmt.Errorf("invalid expression")
	}

	return stack[0], nil
}

// Helper methods
func (p *Parser) isDigit(ch rune) bool {
	return ch >= '0' && ch <= '9'
}

func (p *Parser) isNumber(s string) bool {
	_, err := strconv.ParseFloat(s, 64)
	return err == nil
}

func (p *Parser) isOperator(token string) bool {
	_, ok := p.ops[token]
	return ok || token == p.unaryMinus
}

func (p *Parser) precedence(op string) int {
	if op == p.unaryMinus {
		return 4
	}
	if p, ok := p.ops[op]; ok {
		return p
	}
	return 0
}

// Global function for backward compatibility
func EvaluateExpression(expr string, calc *calculator.Calculator) (float64, error) {
	parser := NewParser(calc)
	return parser.EvaluateExpression(expr)
}
