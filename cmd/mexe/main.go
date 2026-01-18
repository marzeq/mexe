package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/marzeq/mexe"
)

func main() {
	prec := flag.Int("p", 12, "output precision (significant digits)")
	degrees := flag.Bool("d", false, "use degrees mode for trigonometric functions and their inverses")
	check := flag.Bool("c", false, "check syntax only")
	flag.Parse()

	if flag.NArg() < 1 {
		fmt.Fprintf(os.Stderr, "Usage: %s [-p N] [-c] [-d] <expression>\n", os.Args[0])
		os.Exit(1)
	}

	var expr strings.Builder
	expr.WriteString(flag.Arg(0))
	for i := 1; i < flag.NArg(); i++ {
		expr.WriteString(" " + flag.Arg(i))
	}

	if *check {
		if err := mexe.CheckSyntax(expr.String()); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		return
	}

	val, err := mexe.Evaluate(expr.String(), *degrees)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	fmt.Println(formatNumber(val, *prec))
}
