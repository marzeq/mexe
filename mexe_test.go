package mexe

import (
	"strings"
	"testing"
)

func TestDefinitions(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{"phi=1.6180339887; phi*2", "3.2360679774"},
		{"f(x)=2x+3\nf(4)", "11"},
		{"sum(x,y,z)=x+y+z; sum(1,2,3)", "6"},
		{"a=2\n\nb=3; f(x,y)=a*x+b*y\n\nf(4,5)\n\n", "23"},
		{"pi=3; pi+1", "4"},
		{"f()=42; f()", "42"},
		{"f(x)=x+1; g(x,y)=f(x)*y; g(3,4)", "16"},
		{"x=10; f(x)=x+1; f(2)+x", "13"},
		{"f(x)=x+1; f(x)=2x; f(3)", "6"},
	}
	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			if err := CheckSyntax(tt.input); err != nil {
				t.Fatalf("syntax: %v", err)
			}
			got, err := Evaluate(tt.input, false)
			if err != nil {
				t.Fatal(err)
			}
			if got.Format(10) != tt.want {
				t.Fatalf("got %s, want %s", got.Format(10), tt.want)
			}
		})
	}
}

func TestDefinitionErrors(t *testing.T) {
	for _, input := range []string{
		"", " \n ; ", "1+2;3", "x=1", "f(x,x)=x;f(1,2)",
		"f(x)=x;f(1,2)=3;f(1)", "x=1; x=2",
	} {
		// The final statement is always an expression; x=1 alone is invalid.
		// Repeated definitions are permitted, so the final case is invalid only
		// because it attempts an assignment in the final statement.
		if err := CheckSyntax(input); err == nil {
			t.Errorf("expected syntax error for %q", input)
		}
	}
	for _, input := range []string{"f(x,y)=x+y;f(1)", "a=missing; a"} {
		if _, err := Evaluate(input, false); err == nil {
			t.Errorf("expected evaluation error for %q", input)
		}
	}
	if err := CheckSyntax("f(x,y)=x+y;f(1)"); err != nil {
		t.Fatalf("syntax check should not evaluate: %v", err)
	}
	if _, err := Evaluate("f(x,y)=x+y;f(1)", false); err == nil || !strings.Contains(err.Error(), "invalid number of arguments") {
		t.Fatalf("expected arity error, got %v", err)
	}
}
