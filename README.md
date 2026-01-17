# mexe

A simple and lightweight maths expression evaluator designed for use in external scripts and applications

## Usage

```bash
mexe "2 + 2 * (3 - 1)"
```

### Syntax check only

Parse the expression without evaluating it:

```bash
mexe -c "10 / 2 + 5"
```

Exit code is 0 if syntax is valid, non-zero otherwise.

### Output precision

Control the number of significant digits displayed:

```bash
mexe -p 16 "2 ^ (1/7)"
```

Default precision is 12 significant digits.

### Degree mde

Make trigonometric functions interpret input in degrees instead of radians, and inverse functions return degrees:

```bash
mexe -d "sin(90) + cos(0)"
```

### Caveats

If your expression starts with a hyphen (`-`), you must add a `--` before the expression to prevent it from being interpreted as a flag:

```bash
mexe -- "-5 + 3"
```

For use in external programs where you're not the one controlling the input, it is best to do this always.

## Build

```bash
go build ./cmd/mexe
```

## Install

```bash
go install ./cmd/mexe
```

Or manually copy the built binary to a directory in your PATH.

## Supported syntax

### Operators

| Operator | Description                  |
| -------- | ---------------------------- |
| +        | Addition                     |
| -        | Subtraction / unary negation |
| \*       | Multiplication               |
| /        | Division                     |
| %        | Modulo                       |
| mod      | Modulo (infix keyword)       |
| ^        | Exponentiation               |
| !        | Factorial (postfix)          |
| \|x\|    | Absolute value               |

### Implicit multiplication

All of the following are valid:

```
2pi
2(3 + 1)
3sqrt(2)
sin(2pi)
```

## Constants

| Name | Description    |
| ---- | -------------- |
| pi   | π              |
| e    | Euler’s number |

## Functions

### Trigonometric

| Function                   | Description             |
| -------------------------- | ----------------------- |
| sin(x)                     | Sine                    |
| cos(x)                     | Cosine                  |
| tan(x)                     | Tangent                 |
| cot(x)                     | Cotangent               |
| asin(x), arcsin(x)         | Inverse sine            |
| acos(x), arccos(x)         | Inverse cosine          |
| atan(x), arctan(x)         | Inverse tangent         |
| atan2(x, y), arctan2(x, y) | Two-argument arctangent |

### Hyperbolic

| Function             | Description                |
| -------------------- | -------------------------- |
| sinh(x)              | Hyperbolic sine            |
| cosh(x)              | Hyperbolic cosine          |
| tanh(x)              | Hyperbolic tangent         |
| asinh(x), arcsinh(x) | Inverse hyperbolic sine    |
| acosh(x), arccosh(x) | Inverse hyperbolic cosine  |
| atanh(x), arctanh(x) | Inverse hyperbolic tangent |

### Logarithmic / Exponential

| Function  | Description       |
| --------- | ----------------- |
| ln(x)     | Natural logarithm |
| log(b, x) | Logarithm base b  |
| log2(x)   | Base-2 logarithm  |
| log10(x)  | Base-10 logarithm |
| exp(x)    | Exponential (eˣ)  |

### Roots and Powers

| Function    | Description             |
| ----------- | ----------------------- |
| sqrt(x)     | Square root             |
| cbrt(x)     | Cube root               |
| root(x, n)  | n-th root               |
| pow(x, y)   | x raised to the power y |
| hypot(x, y) | √(x² + y²)              |

### Rounding and Decomposition

| Function | Description              |
| -------- | ------------------------ |
| ceil(x)  | Ceiling                  |
| floor(x) | Floor                    |
| round(x) | Nearest integer          |
| trunc(x) | Truncate fractional part |
| fract(x) | Fractional part          |

### Arithmetic Helpers

| Function           | Description                            |
| ------------------ | -------------------------------------- |
| abs(x)             | Absolute value                         |
| sign(x)            | Sign of x                              |
| min(x, ...)        | Minimum of arguments                   |
| max(x, ...)        | Maximum of arguments                   |
| clamp(x, min, max) | Clamp x to range                       |
| mod(x, y)          | Modulo                                 |
| fact(n)            | Factorial (n must be a natural number) |

## My rationale

I developed this project for the launcher in my custom Hyprland shell to
massively speed up the responsiveness of mathematical calculations compared to
my old solution using Python's `eval()` with stripped globals.

What was especially important to me was the fact that if the expression has a syntax error,
no evaluation is performed at all, preventing unnecessary delays, which is crucial in a launcher.
