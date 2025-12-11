# Calculator Built with Go

A feature-rich, high-performance scientific calculator application built with Go. This calculator supports both command-line and interactive modes with advanced mathematical functions, variable support, and a clean architecture.

## Features

### 🧮 Core Operations
- **Basic Arithmetic**: Addition, subtraction, multiplication, division
- **Advanced Operations**: Exponentiation, modulus, percentage, factorial
- **Precision Control**: Configurable precision up to 20 digits
- **Memory Functions**: Store, recall, add to memory

### 📐 Scientific Functions
- **Trigonometric**: sin, cos, tan, asin, acos, atan (degrees/radians)
- **Hyperbolic**: sinh, cosh, tanh
- **Logarithmic**: log (natural), log10
- **Exponential**: exp, power functions
- **Roots**: sqrt, cbrt
- **Rounding**: round, floor, ceil with precision

### 🔧 Advanced Features
- **Variable Support**: Define and use variables with `set x=5`
- **Expression Parser**: Shunting-yard algorithm with operator precedence
- **Command History**: View and clear calculation history
- **Configuration**: Save/load settings, angle mode, precision, color themes
- **Error Handling**: Comprehensive validation and helpful error messages
- **Unit Conversions**: (Planned feature) length, weight, temperature, angles
- **Statistical Functions**: (Planned feature) mean, median, stddev, etc.

### 🎨 User Interface
- **Interactive Mode**: Color-coded prompt with syntax highlighting
- **Command-Line Mode**: Evaluate expressions directly from terminal
- **Color Output**: Optional ANSI color codes for better readability
- **Multi-format Results**: Display as decimal, scientific notation, and fraction

## Installation

### Prerequisites
- Go 1.21 or higher

### Build from Source
```bash
    # Clone the repository
    git clone https://github.com/Oluwaseyi89/calculator-built-with-go.git
    cd calculator-built-with-go

    # Build the application
    go build -o calculator

    # Run the calculator
    ./calculator
```

### Using go install
```bash
    go install github.com/Oluwaseyi89/calculator-built-with-go@latest
```

### Usage
#### Interactive Mode
```bash
    ./calculator
    # or
    ./calculator --interactive
```

##### Example interactive session:
```bash
    SCIENTIFIC CALCULATOR v2.0
    ==========================

    Type 'help' for full help, 'exit' to quit
    Type expressions like: 2 + 3 * 4, sin(pi/2), sqrt(16)

    calc> 2 + 3 * 4
    2 + 3 * 4 = 14

    calc> sin(pi/2)
    sin(pi/2) = 1
    Integer: 1

    calc> set radius = 5
    Set radius = 5

    calc> pi * radius^2
    pi * radius^2 = 78.53981634
    Scientific: 7.853982e+01
```

#### Command-Line Mode
```bash
    # Evaluate a single expression
    ./calculator --eval "2 + 3 * 4"
    ./calculator -e "sin(pi/4)"

    # Show help
    ./calculator --help

    # Show version
    ./calculator --version
```

#### Available Commands

| Command | Description |
|---------|-------------|
| `help` | Show comprehensive help |
| `exit` or `quit` | Exit the calculator |
| `clear` | Clear memory |
| `mem` | Show memory value |
| `history` | Show calculation history |
| `clearhist` | Clear history |
| `settings` | Show current settings |
| `mode` | Toggle between degrees and radians |
| `examples` | Show usage examples |
| `units` | Show unit conversion help |
| `stats` | Show statistical functions help |
| `set var=expr` | Set a variable |
| `deg expr` | Evaluate expression in degrees mode |
| `rad expr` | Evaluate expression in radians mode |
| `precision N` | Set display precision (1-20) |


### Expression Syntax
#### Basic Arithmetic

```
    2 + 3 * 4           # Standard operator precedence
    (2 + 3) * 4         # Parentheses for grouping
    10 % 3              # Modulus
    5!                  # Factorial
    2^10                # Exponentiation
```

#### Functions

```
    sin(pi/2)           # Trigonometric functions
    log(e^2)            # Natural logarithm
    sqrt(16)            # Square root
    exp(2)              # Exponential e^x
    abs(-5.5)           # Absolute value
    round(3.14159, 2)   # Round to 2 decimal places
```

#### Variables and Constants

```
    set x = 5           # Assign variable
    x^2 + 3*x + 2       # Use variable
    pi                  # π constant (3.14159...)
    e                   # Euler's number (2.71828...)
    ans                 # Previous result
```

### Project Structure

```
    calculator-built-with-go/
    ├── calculator/          # Core calculator engine
    │   ├── arithmetic.go   # Basic arithmetic operations
    │   ├── scientific.go   # Scientific functions
    │   ├── memory.go       # Memory and history management
    │   └── history.go      # History tracking
    ├── parser/             # Expression parsing
    │   └── expression.go   # Shunting-yard algorithm parser
    ├── utils/              # Utility functions
    │   ├── helpers.go      # Helper functions
    │   └── validators.go   # Input validation
    ├── main.go            # Main application entry point
    ├── go.mod             # Go module definition
    └── README.md          # This file
```

## Technical Details

### Parser Implementation
The calculator uses a modified Shunting-yard algorithm to parse mathematical expressions with:
- Operator precedence handling (`^` > `*/%` > `+-`)
- Unary minus support
- Function argument parsing
- Parentheses balancing
- Variable substitution

### Error Handling
Comprehensive validation includes:
- Syntax checking
- Division by zero detection
- Invalid function arguments
- Overflow/underflow detection
- Balanced parentheses validation

### Performance Features
- Concurrent-safe memory operations
- Time tracking for calculations
- Efficient history management (last 100 entries)
- Optimized mathematical operations using Go's math package

### Configuration
The calculator can be configured through commands or by editing the configuration file:

#### Settings Commands

```bash
    # Set angle mode
    mode                # Toggle between deg/rad
    deg expr            # Evaluate in degrees
    rad expr            # Evaluate in radians

    # Set precision
    precision 15        # Set to 15 decimal places

    # Toggle features
    # (Planned) scientific on/off
    # (Planned) color on/off
```

### Examples
#### Basic Calculations

```bash
    calc> 2 + 3 * 4
    14

    calc> (2 + 3) * 4
    20

    calc> 10 % 3
    1

    calc> 5!
    120
```

#### Scientific Calculations
```bash 
    calc> sin(pi/2)
    1

    calc> sqrt(16)
    4

    calc> log(e^2)
    2

    calc> 2^10
    1024
```

#### Using Variables
```bash
    calc> set radius = 5
    Set radius = 5

    calc> pi * radius^2
    78.53981634

    calc> set height = 10
    Set height = 10

    calc> pi * radius^2 * height
    785.3981634
```

#### Unit Conversions (Planned)
```bash
    # Future syntax
    calc> convert(100, 'm', 'ft')
    328.084

    calc> convert(25, 'C', 'F')
    77
```

### Development
#### Running Tests
```bash
    go test ./...
```

#### Building for Different Platforms
```bash
    # Linux
    GOOS=linux GOARCH=amd64 go build -o calculator-linux

    # macOS
    GOOS=darwin GOARCH=amd64 go build -o calculator-macos

    # Windows
    GOOS=windows GOARCH=amd64 go build -o calculator.exe
```

### Code Structure
- **`calculator` package**: Pure mathematical operations
- **`parser` package**: Expression parsing and evaluation
- **`utils` package**: Shared utilities and validation
- **`main` package**: User interface and application logic

### Planned Features
- Unit conversion system
- Statistical functions
- Matrix operations
- Complex number support
- Graphing capabilities
- Scripting support
- Plugin system for custom functions
- Web interface

### Contributing
1. Fork the repository
2. Create a feature branch (`git checkout -b feature/AmazingFeature`)
3. Commit your changes (`git commit -m 'Add some AmazingFeature'`)
4. Push to the branch (`git push origin feature/AmazingFeature`)
5. Open a Pull Request
