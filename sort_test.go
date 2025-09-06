package sort

import (
	"context"
	"strings"
	"testing"
	"time"

	yup "github.com/yupsh/framework"

	"github.com/yupsh/sort/opt"
)

func TestSortBasic(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		flags    []any
		expected string
	}{
		{
			name:     "basic alphabetical sort",
			input:    "charlie\nbravo\nalpha\ndelta\n",
			flags:    []any{},
			expected: "alpha\nbravo\ncharlie\ndelta\n",
		},
		{
			name:     "reverse sort",
			input:    "alpha\nbravo\ncharlie\n",
			flags:    []any{opt.Reverse},
			expected: "charlie\nbravo\nalpha\n",
		},
		{
			name:     "numeric sort",
			input:    "10\n2\n1\n20\n",
			flags:    []any{opt.Numeric},
			expected: "1\n2\n10\n20\n",
		},
		{
			name:     "numeric sort with non-numeric values",
			input:    "10\nabc\n2\ndef\n1\n",
			flags:    []any{opt.Numeric},
			expected: "1\n2\n10\nabc\ndef\n",
		},
		{
			name:     "case insensitive sort",
			input:    "Apple\nbanana\nCherry\n",
			flags:    []any{opt.IgnoreCase},
			expected: "Apple\nbanana\nCherry\n",
		},
		{
			name:     "case sensitive sort (default)",
			input:    "Apple\nbanana\nCherry\n",
			flags:    []any{},
			expected: "Apple\nCherry\nbanana\n",
		},
		{
			name:     "unique sort",
			input:    "apple\nbanana\napple\ncherry\nbanana\n",
			flags:    []any{opt.Unique},
			expected: "apple\nbanana\ncherry\n",
		},
		{
			name:     "unique with case insensitive",
			input:    "Apple\nbanana\napple\nCherry\nBANANA\n",
			flags:    []any{opt.Unique, opt.IgnoreCase},
			expected: "Apple\nbanana\nCherry\n",
		},
		{
			name:     "empty input",
			input:    "",
			flags:    []any{},
			expected: "",
		},
		{
			name:     "single line",
			input:    "single\n",
			flags:    []any{},
			expected: "single\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := Sort(tt.flags...)

			var output strings.Builder
			var stderr strings.Builder

			ctx := context.Background()
			err := cmd.Execute(ctx, strings.NewReader(tt.input), &output, &stderr)

			if err != nil {
				t.Fatalf("Execute failed: %v\nStderr: %s", err, stderr.String())
			}

			result := output.String()
			if result != tt.expected {
				t.Errorf("Expected %q, got %q", tt.expected, result)
			}
		})
	}
}

func TestSortByField(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		field     opt.Field
		delimiter opt.Delimiter
		flags     []any
		expected  string
	}{
		{
			name:      "sort by second field (default space delimiter)",
			input:     "john 30\nalice 25\nbob 35\n",
			field:     opt.Field(2),
			delimiter: opt.Delimiter(""),
			flags:     []any{},
			expected:  "alice 25\njohn 30\nbob 35\n",
		},
		{
			name:      "sort by second field numeric",
			input:     "john 30\nalice 25\nbob 35\n",
			field:     opt.Field(2),
			delimiter: opt.Delimiter(""),
			flags:     []any{opt.Numeric},
			expected:  "alice 25\njohn 30\nbob 35\n",
		},
		{
			name:      "sort by first field with comma delimiter",
			input:     "charlie,30\nalpha,25\nbravo,35\n",
			field:     opt.Field(1),
			delimiter: opt.Delimiter(","),
			flags:     []any{},
			expected:  "alpha,25\nbravo,35\ncharlie,30\n",
		},
		{
			name:      "sort by second field with comma delimiter numeric",
			input:     "charlie,30\nalpha,25\nbravo,35\n",
			field:     opt.Field(2),
			delimiter: opt.Delimiter(","),
			flags:     []any{opt.Numeric},
			expected:  "alpha,25\ncharlie,30\nbravo,35\n",
		},
		{
			name:      "sort by non-existent field",
			input:     "one two\nthree four\n",
			field:     opt.Field(5),
			delimiter: opt.Delimiter(""),
			flags:     []any{},
			expected:  "one two\nthree four\n",
		},
		{
			name:      "sort by field with tab delimiter",
			input:     "name\tage\nalice\t25\nbob\t30\n",
			field:     opt.Field(2),
			delimiter: opt.Delimiter("\t"),
			flags:     []any{opt.Numeric},
			expected:  "alice\t25\nbob\t30\nname\tage\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			flags := append([]any{tt.field}, tt.flags...)
			if tt.delimiter != "" {
				flags = append(flags, tt.delimiter)
			}

			cmd := Sort(flags...)

			var output strings.Builder
			var stderr strings.Builder

			ctx := context.Background()
			err := cmd.Execute(ctx, strings.NewReader(tt.input), &output, &stderr)

			if err != nil {
				t.Fatalf("Execute failed: %v\nStderr: %s", err, stderr.String())
			}

			result := output.String()
			if result != tt.expected {
				t.Errorf("Expected %q, got %q", tt.expected, result)
			}
		})
	}
}

func TestSortCombinedFlags(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		flags    []any
		expected string
	}{
		{
			name:     "reverse numeric sort",
			input:    "1\n10\n2\n20\n",
			flags:    []any{opt.Numeric, opt.Reverse},
			expected: "20\n10\n2\n1\n",
		},
		{
			name:     "unique reverse sort",
			input:    "apple\nbanana\napple\ncherry\n",
			flags:    []any{opt.Unique, opt.Reverse},
			expected: "cherry\nbanana\napple\n",
		},
		{
			name:     "field sort with reverse and unique",
			input:    "alice 25\nbob 25\ncharlie 30\nalice 25\n",
			flags:    []any{opt.Field(2), opt.Numeric, opt.Reverse, opt.Unique},
			expected: "charlie 30\nalice 25\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := Sort(tt.flags...)

			var output strings.Builder
			var stderr strings.Builder

			ctx := context.Background()
			err := cmd.Execute(ctx, strings.NewReader(tt.input), &output, &stderr)

			if err != nil {
				t.Fatalf("Execute failed: %v\nStderr: %s", err, stderr.String())
			}

			result := output.String()
			if result != tt.expected {
				t.Errorf("Expected %q, got %q", tt.expected, result)
			}
		})
	}
}

func TestSortRandom(t *testing.T) {
	input := "line1\nline2\nline3\nline4\nline5\n"

	cmd := Sort(opt.Random)

	var output strings.Builder
	var stderr strings.Builder

	ctx := context.Background()
	err := cmd.Execute(ctx, strings.NewReader(input), &output, &stderr)

	if err != nil {
		t.Fatalf("Execute failed: %v\nStderr: %s", err, stderr.String())
	}

	result := output.String()
	lines := strings.Split(strings.TrimSpace(result), "\n")

	// Should have same number of lines
	if len(lines) != 5 {
		t.Errorf("Expected 5 lines, got %d", len(lines))
	}

	// Should contain all original lines (though order may be different)
	expectedLines := map[string]bool{
		"line1": true,
		"line2": true,
		"line3": true,
		"line4": true,
		"line5": true,
	}

	for _, line := range lines {
		if !expectedLines[line] {
			t.Errorf("Unexpected line in output: %q", line)
		}
	}
}

func TestSortContextCancellation(t *testing.T) {
	// Create a large input that would take time to sort
	var inputBuilder strings.Builder
	for i := 0; i < 10000; i++ {
		inputBuilder.WriteString("line")
		inputBuilder.WriteString(string(rune(i % 1000)))
		inputBuilder.WriteString("\n")
	}

	cmd := Sort()

	// Create a context that will be cancelled quickly
	ctx, cancel := context.WithCancel(context.Background())

	var output strings.Builder
	var stderr strings.Builder

	// Cancel context immediately
	cancel()

	err := cmd.Execute(ctx, strings.NewReader(inputBuilder.String()), &output, &stderr)

	// Should detect cancellation and return error
	if err == nil {
		t.Error("Expected context cancellation error, got nil")
	}

	if !strings.Contains(err.Error(), "context canceled") && !strings.Contains(err.Error(), "context cancelled") {
		t.Errorf("Expected context cancellation error, got: %v", err)
	}
}

func TestSortContextCancellationTimeout(t *testing.T) {
	// Create a large input
	var inputBuilder strings.Builder
	for i := 0; i < 10000; i++ {
		inputBuilder.WriteString("line")
		inputBuilder.WriteString(string(rune(i % 1000)))
		inputBuilder.WriteString("\n")
	}

	cmd := Sort()

	// Create a context with a very short timeout
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Millisecond)
	defer cancel()

	var output strings.Builder
	var stderr strings.Builder

	err := cmd.Execute(ctx, strings.NewReader(inputBuilder.String()), &output, &stderr)

	// Should timeout and return error (or succeed if fast enough)
	if err != nil {
		// If there's an error, it should be a context cancellation
		if !strings.Contains(err.Error(), "context") {
			t.Errorf("Expected context error, got: %v", err)
		}
	}
}

func TestSortWithFiles(t *testing.T) {
	// Test with positional arguments (files)
	cmd := Sort("testfile1", "testfile2")

	// Since these files don't exist, we expect an error
	var output strings.Builder
	var stderr strings.Builder

	ctx := context.Background()
	err := cmd.Execute(ctx, nil, &output, &stderr)

	// Should get file not found error
	if err == nil {
		t.Error("Expected file not found error, got nil")
	}
}

func TestSortEdgeCases(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		flags    []any
		expected string
	}{
		{
			name:     "lines with only whitespace",
			input:    "alpha\n\t\nbravo\n   \ncharlie\n",
			flags:    []any{},
			expected: "\t\n   \nalpha\nbravo\ncharlie\n",
		},
		{
			name:     "numbers with leading zeros",
			input:    "001\n10\n002\n",
			flags:    []any{opt.Numeric},
			expected: "001\n002\n10\n",
		},
		{
			name:     "decimal numbers",
			input:    "1.5\n1.10\n1.2\n",
			flags:    []any{opt.Numeric},
			expected: "1.10\n1.2\n1.5\n",
		},
		{
			name:     "negative numbers",
			input:    "-5\n-10\n-1\n0\n5\n",
			flags:    []any{opt.Numeric},
			expected: "-10\n-5\n-1\n0\n5\n",
		},
		{
			name:     "mixed case with ignore case",
			input:    "Apple\napple\nAPPLE\n",
			flags:    []any{opt.IgnoreCase, opt.Unique},
			expected: "Apple\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := Sort(tt.flags...)

			var output strings.Builder
			var stderr strings.Builder

			ctx := context.Background()
			err := cmd.Execute(ctx, strings.NewReader(tt.input), &output, &stderr)

			if err != nil {
				t.Fatalf("Execute failed: %v\nStderr: %s", err, stderr.String())
			}

			result := output.String()
			if result != tt.expected {
				t.Errorf("Expected %q, got %q", tt.expected, result)
			}
		})
	}
}

func TestSortInterface(t *testing.T) {
	// Verify that Sort command implements yup.Command interface
	var _ yup.Command = Sort()
}

func BenchmarkSortSmall(b *testing.B) {
	input := "delta\ncharlie\nbravo\nalpha\n"
	cmd := Sort()
	ctx := context.Background()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var output strings.Builder
		var stderr strings.Builder
		cmd.Execute(ctx, strings.NewReader(input), &output, &stderr)
	}
}

func BenchmarkSortLarge(b *testing.B) {
	var inputBuilder strings.Builder
	for i := 1000; i > 0; i-- {
		inputBuilder.WriteString("line")
		inputBuilder.WriteString(string(rune(i % 100)))
		inputBuilder.WriteString("\n")
	}
	input := inputBuilder.String()

	cmd := Sort()
	ctx := context.Background()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var output strings.Builder
		var stderr strings.Builder
		cmd.Execute(ctx, strings.NewReader(input), &output, &stderr)
	}
}

func BenchmarkSortNumeric(b *testing.B) {
	var inputBuilder strings.Builder
	for i := 1000; i > 0; i-- {
		inputBuilder.WriteString(string(rune(i)))
		inputBuilder.WriteString("\n")
	}
	input := inputBuilder.String()

	cmd := Sort(opt.Numeric)
	ctx := context.Background()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var output strings.Builder
		var stderr strings.Builder
		cmd.Execute(ctx, strings.NewReader(input), &output, &stderr)
	}
}

func BenchmarkSortByField(b *testing.B) {
	var inputBuilder strings.Builder
	for i := 1000; i > 0; i-- {
		inputBuilder.WriteString("name")
		inputBuilder.WriteString(string(rune(i % 100)))
		inputBuilder.WriteString(" ")
		inputBuilder.WriteString(string(rune(i)))
		inputBuilder.WriteString("\n")
	}
	input := inputBuilder.String()

	cmd := Sort(opt.Field(2), opt.Numeric)
	ctx := context.Background()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var output strings.Builder
		var stderr strings.Builder
		cmd.Execute(ctx, strings.NewReader(input), &output, &stderr)
	}
}

// Example tests for documentation
func ExampleSort() {
	cmd := Sort()
	ctx := context.Background()

	input := strings.NewReader("charlie\nbravo\nalpha\n")
	var output strings.Builder
	cmd.Execute(ctx, input, &output, &strings.Builder{})
	// Output would be: alpha\nbravo\ncharlie\n
}

func ExampleSort_numeric() {
	cmd := Sort(opt.Numeric)
	ctx := context.Background()

	input := strings.NewReader("10\n2\n1\n20\n")
	var output strings.Builder
	cmd.Execute(ctx, input, &output, &strings.Builder{})
	// Output would be: 1\n2\n10\n20\n
}

func ExampleSort_byField() {
	cmd := Sort(opt.Field(2), opt.Numeric)
	ctx := context.Background()

	input := strings.NewReader("alice 30\nbob 25\ncharlie 35\n")
	var output strings.Builder
	cmd.Execute(ctx, input, &output, &strings.Builder{})
	// Output would be: bob 25\nalice 30\ncharlie 35\n
}
