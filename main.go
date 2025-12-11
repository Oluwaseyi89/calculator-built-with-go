package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/Oluwaseyi89/calculator-built-with-go/calculator"
	"github.com/Oluwaseyi89/calculator-built-with-go/parser"
	"github.com/Oluwaseyi89/calculator-built-with-go/utils"
)

type CalculatorApp struct {
	calc    *calculator.Calculator
	parser  *parser.Parser // Add this line
	scanner *bufio.Scanner
	history []string
	config  *Config
}

type Config struct {
	AngleMode    string // "deg" or "rad"
	Precision    int
	Scientific   bool
	ShowHistory  bool
	ColorEnabled bool
}

func NewCalculatorApp() *CalculatorApp {
	calc := calculator.NewCalculator()
	return &CalculatorApp{
		calc:    calc,
		parser:  parser.NewParser(calc), // Initialize parser
		scanner: bufio.NewScanner(os.Stdin),
		history: make([]string, 0),
		config: &Config{
			AngleMode:    "rad",
			Precision:    10,
			Scientific:   false,
			ShowHistory:  true,
			ColorEnabled: true,
		},
	}
}

func (app *CalculatorApp) Run() {
	app.showWelcome()
	app.showQuickHelp()
	app.loadConfig()
	app.mainLoop()
}

func (app *CalculatorApp) mainLoop() {
	for {
		app.printPrompt()

		if !app.scanner.Scan() {
			break
		}

		input := strings.TrimSpace(app.scanner.Text())

		if app.handleCommand(input) {
			continue
		}

		if app.handleSpecial(input) {
			continue
		}

		if err := utils.ValidateExpression(input); err != nil {
			app.printError(fmt.Sprintf("Invalid expression: %v", err))
			continue
		}

		app.evaluateExpression(input)
	}
}

func (app *CalculatorApp) handleCommand(input string) bool {
	command := strings.ToLower(input)

	commandHandlers := map[string]func(){
		"exit":      app.handleExit,
		"quit":      app.handleExit,
		"help":      app.showFullHelp,
		"clear":     app.handleClearMemory,
		"history":   app.handleShowHistory,
		"mem":       app.handleShowMemory,
		"settings":  app.handleSettings,
		"mode":      app.handleModeToggle,
		"config":    app.handleConfig,
		"clearhist": app.handleClearHistory,
		"version":   app.handleVersion,
		"license":   app.handleLicense,
		"examples":  app.showExamples,
		"units":     app.showUnitConversions,
		"stats":     app.showStatistics,
	}

	if handler, exists := commandHandlers[command]; exists {
		handler()
		return true
	}

	return false
}

func (app *CalculatorApp) handleSpecial(input string) bool {
	specialHandlers := map[string]func(string){
		"save ": app.handleSave,
		"load ": app.handleLoad,
		// "export ":    app.handleExport,
		// "import ":    app.handleImport,
		"set ": app.handleSet,
		// "del ":       app.handleDelete,
		// "var ":       app.handleVariable,
		"deg ":       app.handleDegrees,
		"rad ":       app.handleRadians,
		"precision ": app.handlePrecision,
	}

	for prefix, handler := range specialHandlers {
		if strings.HasPrefix(strings.ToLower(input), prefix) {
			arg := strings.TrimSpace(input[len(prefix):])
			handler(arg)
			return true
		}
	}

	return false
}

func (app *CalculatorApp) evaluateExpression(expr string) {
	startTime := time.Now()

	// Set angle mode in parser
	app.parser.SetAngleMode(app.config.AngleMode)

	// Evaluate expression
	result, err := app.parser.EvaluateExpression(expr)

	duration := time.Since(startTime)

	if err != nil {
		app.printError(fmt.Sprintf("Evaluation error: %v", err))
		return
	}

	app.displayResult(expr, result, duration)

	app.addToCommandHistory(expr)
}

func (app *CalculatorApp) displayResult(expr string, result float64, duration time.Duration) {
	formattedResult := utils.FormatNumber(result)

	if app.config.ColorEnabled {
		app.printColorizedResult(expr, formattedResult, duration)
	} else {
		fmt.Printf("\n%s = %s\n", expr, formattedResult)
	}

	if duration > 10*time.Millisecond {
		fmt.Printf("  (calculated in %s)\n", utils.FormatDuration(duration))
	}

	// Show additional formats
	if app.config.Scientific {
		fmt.Printf("  Scientific: %.6e\n", result)
	}

	// Show fraction approximation for simple decimals
	if result != 0 && result == float64(int64(result)) {
		fmt.Printf("  Integer: %.0f\n", result)
	} else if result > 0.01 && result < 100 {
		num, den := utils.ConvertToFraction(result, 100)
		if den > 1 {
			fmt.Printf("  Fraction: %d/%d\n", num, den)
		}
	}
}

func (app *CalculatorApp) printColorizedResult(expr, formattedResult string, duration time.Duration) {
	// ANSI color codes
	const (
		reset   = "\033[0m"
		green   = "\033[32m"
		yellow  = "\033[33m"
		blue    = "\033[34m"
		magenta = "\033[35m"
		cyan    = "\033[36m"
		bold    = "\033[1m"
	)

	fmt.Printf("\n%s%s%s %s= %s%s%s\n",
		bold, blue, expr,
		yellow,
		bold, green, formattedResult)

	if duration > 10*time.Millisecond {
		fmt.Printf("  %s(%s)%s\n", cyan, utils.FormatDuration(duration), reset)
	} else {
		fmt.Print(reset)
	}
}

func (app *CalculatorApp) printPrompt() {
	if app.config.ColorEnabled {
		fmt.Print("\033[1;36mcalc>\033[0m ")
	} else {
		fmt.Print("\ncalc> ")
	}
}

func (app *CalculatorApp) printError(msg string) {
	if app.config.ColorEnabled {
		fmt.Printf("\033[1;31mError:\033[0m %s\n", msg)
	} else {
		fmt.Printf("Error: %s\n", msg)
	}
}

func (app *CalculatorApp) printSuccess(msg string) {
	if app.config.ColorEnabled {
		fmt.Printf("\033[1;32m%s\033[0m\n", msg)
	} else {
		fmt.Printf("%s\n", msg)
	}
}

func (app *CalculatorApp) printInfo(msg string) {
	if app.config.ColorEnabled {
		fmt.Printf("\033[1;34m%s\033[0m\n", msg)
	} else {
		fmt.Printf("%s\n", msg)
	}
}

// Command Handlers
func (app *CalculatorApp) handleExit() {
	app.printInfo("Goodbye!")
	app.saveConfig()
	os.Exit(0)
}

func (app *CalculatorApp) handleClearMemory() {
	app.calc.ClearMemory()
	app.printSuccess("Memory cleared")
}

func (app *CalculatorApp) handleShowMemory() {
	value := app.calc.GetMemory()
	app.printInfo(fmt.Sprintf("Memory value: %s", utils.FormatNumber(value)))
}

func (app *CalculatorApp) handleShowHistory() {
	app.calc.ShowHistory()
}

func (app *CalculatorApp) handleClearHistory() {
	app.calc.ClearHistory()
	app.printSuccess("History cleared")
}

func (app *CalculatorApp) handleSettings() {
	app.printInfo("Current Settings:")
	fmt.Printf("  Angle mode: %s\n", app.config.AngleMode)
	fmt.Printf("  Precision: %d digits\n", app.config.Precision)
	fmt.Printf("  Scientific mode: %v\n", app.config.Scientific)
	fmt.Printf("  Show history: %v\n", app.config.ShowHistory)
	fmt.Printf("  Color output: %v\n", app.config.ColorEnabled)
	fmt.Println("\nCommands: mode deg/rad, precision N, scientific on/off")
}

func (app *CalculatorApp) handleModeToggle() {
	if app.config.AngleMode == "deg" {
		app.config.AngleMode = "rad"
		app.printSuccess("Switched to radians mode")
	} else {
		app.config.AngleMode = "deg"
		app.printSuccess("Switched to degrees mode")
	}
	app.saveConfig()
}

func (app *CalculatorApp) handleConfig() {
	app.handleSettings()
}

func (app *CalculatorApp) handleVersion() {
	app.showWelcome()
	fmt.Println("Version: 2.0.0")
	fmt.Println("Built with Go 1.21+")
	fmt.Println("GitHub: github.com/Oluwaseyi89/calculator-built-with-go")
}

func (app *CalculatorApp) handleLicense() {
	fmt.Println("MIT License")
	fmt.Println("Copyright (c) 2024 Your Name")
	fmt.Println("\nPermission is hereby granted...")
}

func (app *CalculatorApp) handleSave(filename string) {
	if filename == "" {
		filename = "calculator_state.json"
	}
	app.printInfo(fmt.Sprintf("Saving state to %s...", filename))
	// Implementation for saving state
}

func (app *CalculatorApp) handleLoad(filename string) {
	if filename == "" {
		filename = "calculator_state.json"
	}
	app.printInfo(fmt.Sprintf("Loading state from %s...", filename))
	// Implementation for loading state
}

func (app *CalculatorApp) handleSet(arg string) {
	parts := strings.SplitN(arg, "=", 2)
	if len(parts) != 2 {
		app.printError("Usage: set variable=value")
		return
	}

	variable := strings.TrimSpace(parts[0])
	valueStr := strings.TrimSpace(parts[1])

	if !utils.IsValidVariableName(variable) {
		app.printError("Invalid variable name")
		return
	}

	// Evaluate the value using parser
	result, err := app.parser.EvaluateExpression(valueStr)
	if err != nil {
		app.printError(fmt.Sprintf("Invalid value: %v", err))
		return
	}

	// Store variable in parser
	err = app.parser.SetVariable(variable, result)
	if err != nil {
		app.printError(fmt.Sprintf("Failed to set variable: %v", err))
		return
	}

	app.printSuccess(fmt.Sprintf("Set %s = %s", variable, utils.FormatNumber(result)))
}

func (app *CalculatorApp) handleDegrees(expr string) {
	// Temporarily set angle mode to degrees
	app.parser.SetAngleMode("deg")
	result, err := app.parser.EvaluateExpression(expr)
	if err != nil {
		app.printError(fmt.Sprintf("Error: %v", err))
		app.parser.SetAngleMode(app.config.AngleMode) // Restore original mode
		return
	}
	app.parser.SetAngleMode(app.config.AngleMode) // Restore original mode
	app.displayResult(fmt.Sprintf("deg(%s)", expr), result, 0)
}

func (app *CalculatorApp) handleRadians(expr string) {
	// Temporarily set angle mode to radians
	app.parser.SetAngleMode("rad")
	result, err := app.parser.EvaluateExpression(expr)
	if err != nil {
		app.printError(fmt.Sprintf("Error: %v", err))
		app.parser.SetAngleMode(app.config.AngleMode) // Restore original mode
		return
	}
	app.parser.SetAngleMode(app.config.AngleMode) // Restore original mode
	app.displayResult(fmt.Sprintf("rad(%s)", expr), result, 0)
}

func (app *CalculatorApp) handlePrecision(arg string) {
	precision, err := strconv.Atoi(arg)
	if err != nil || precision < 1 || precision > 20 {
		app.printError("Precision must be between 1 and 20")
		return
	}
	app.config.Precision = precision
	app.printSuccess(fmt.Sprintf("Precision set to %d digits", precision))
	app.saveConfig()
}

func (app *CalculatorApp) addToCommandHistory(cmd string) {
	app.history = append(app.history, cmd)
	if len(app.history) > 100 {
		app.history = app.history[1:]
	}
}

func (app *CalculatorApp) loadConfig() {
	// Load configuration from file
	// For now, use defaults
}

func (app *CalculatorApp) saveConfig() {
	// Save configuration to file
}

// Display Functions
func (app *CalculatorApp) showWelcome() {
	title := "SCIENTIFIC CALCULATOR v2.0"
	border := strings.Repeat("=", len(title)+4)

	if app.config.ColorEnabled {
		fmt.Printf("\033[1;35m%s\033[0m\n", border)
		fmt.Printf("\033[1;35m  %s  \033[0m\n", title)
		fmt.Printf("\033[1;35m%s\033[0m\n", border)
	} else {
		fmt.Println(border)
		fmt.Printf("  %s  \n", title)
		fmt.Println(border)
	}
}

func (app *CalculatorApp) showQuickHelp() {
	fmt.Println("\nType 'help' for full help, 'exit' to quit")
	fmt.Println("Type expressions like: 2 + 3 * 4, sin(pi/2), sqrt(16)")
}

func (app *CalculatorApp) showFullHelp() {
	app.printInfo("=== CALCULATOR HELP ===")

	helpSections := []struct {
		title   string
		content string
	}{
		{
			"BASIC COMMANDS",
			`  exit, quit      - Exit calculator
  help           - Show this help
  clear          - Clear memory
  mem            - Show memory value
  history        - Show calculation history
  clearhist      - Clear history
  settings       - Show current settings
  mode           - Toggle deg/rad mode`,
		},
		{
			"EXPRESSION SYNTAX",
			`  + - * / ^      - Basic arithmetic
  %              - Modulus/Percentage
  !              - Factorial
  ( )            - Parentheses for grouping
  pi, e          - Mathematical constants
  ans            - Previous result`,
		},
		{
			"MATHEMATICAL FUNCTIONS",
			`  sin(x), cos(x), tan(x)    - Trigonometric
  asin(x), acos(x), atan(x) - Inverse trigonometric
  sinh(x), cosh(x), tanh(x) - Hyperbolic
  sqrt(x), cbrt(x)          - Square/cube root
  log(x), log10(x)          - Natural/base-10 log
  exp(x)                    - Exponential e^x
  abs(x)                    - Absolute value`,
		},
		{
			"ADVANCED COMMANDS",
			`  set var=expr    - Set variable
  deg expr       - Evaluate in degrees mode
  rad expr       - Evaluate in radians mode
  precision N    - Set display precision (1-20)
  examples       - Show usage examples
  units          - Show unit conversions
  stats          - Statistical functions help`,
		},
		{
			"EXAMPLES",
			`  2 + 3 * 4           = 14
  sin(pi/2)           = 1
  2^3 + sqrt(16)      = 12
  log(e^2)            = 2
  deg sin(90)         = 1
  set x = 5           (assign variable)
  x^2 + 3*x + 2       (use variable)`,
		},
	}

	for _, section := range helpSections {
		app.printInfo(fmt.Sprintf("\n%s", section.title))
		fmt.Println(section.content)
	}
}

func (app *CalculatorApp) showExamples() {
	app.printInfo("=== EXAMPLE EXPRESSIONS ===")

	examples := []struct {
		expr     string
		expected string
	}{
		{"2 + 3 * 4", "14"},
		{"(2 + 3) * 4", "20"},
		{"10 % 3", "1"},
		{"5!", "120"},
		{"sin(pi/2)", "1"},
		{"cos(pi)", "-1"},
		{"sqrt(16)", "4"},
		{"2^10", "1024"},
		{"log(e^2)", "2"},
		{"abs(-5.5)", "5.5"},
		{"deg sin(90)", "1"},
		{"rad sin(1.5708)", "≈1"},
		{"exp(1)", "≈2.71828"},
		{"log10(100)", "2"},
		{"asin(0.5)", "0.5236"},
	}

	for _, ex := range examples {
		fmt.Printf("  %-20s → %s\n", ex.expr, ex.expected)
	}
}

func (app *CalculatorApp) showUnitConversions() {
	app.printInfo("=== UNIT CONVERSIONS (Future Feature) ===")
	fmt.Println("  length: meters, feet, inches, miles")
	fmt.Println("  weight: kg, lbs, ounces")
	fmt.Println("  temperature: C, F, K")
	fmt.Println("  angle: deg, rad, grad")
	fmt.Println("\nExample (future): convert(100, 'm', 'ft')")
}

func (app *CalculatorApp) showStatistics() {
	app.printInfo("=== STATISTICAL FUNCTIONS (Future Feature) ===")
	fmt.Println("  mean(values)     - Arithmetic mean")
	fmt.Println("  median(values)   - Median")
	fmt.Println("  mode(values)     - Mode")
	fmt.Println("  stddev(values)   - Standard deviation")
	fmt.Println("  variance(values) - Variance")
	fmt.Println("  min(values)      - Minimum value")
	fmt.Println("  max(values)      - Maximum value")
	fmt.Println("  sum(values)      - Sum of values")
	fmt.Println("\nExample (future): mean([1, 2, 3, 4, 5]) = 3")
}

// Main function
func main() {
	app := NewCalculatorApp()

	// Handle command line arguments
	if len(os.Args) > 1 {
		app.handleCommandLineArgs(os.Args[1:])
		return
	}

	app.Run()
}

func (app *CalculatorApp) handleCommandLineArgs(args []string) {
	switch args[0] {
	case "--version", "-v":
		app.handleVersion()
	case "--help", "-h":
		app.showFullHelp()
	case "--interactive", "-i":
		app.Run()
	case "--eval", "-e":
		if len(args) > 1 {
			expr := strings.Join(args[1:], " ")
			app.evaluateExpression(expr)
		} else {
			fmt.Println("Error: No expression provided for --eval")
		}
	default:
		fmt.Printf("Unknown option: %s\n", args[0])
		fmt.Println("Usage: calculator [--help|--version|--eval EXPR]")
	}
}
