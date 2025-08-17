package sort

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"sort"
	"strconv"
	"strings"

	yup "github.com/yupsh/framework"
	"github.com/yupsh/framework/opt"
	localopt "github.com/yupsh/sort/opt"
)

// Flags represents the configuration options for the sort command
type Flags = localopt.Flags

// Command implementation
type command opt.Inputs[string, Flags]

// Sort creates a new sort command with the given parameters
func Sort(parameters ...any) yup.Command {
	cmd := command(opt.Args[string, Flags](parameters...))
	// Set default delimiter to whitespace
	if cmd.Flags.Delimiter == "" {
		cmd.Flags.Delimiter = " "
	}
	return cmd
}

func (c command) Execute(ctx context.Context, stdin io.Reader, stdout, stderr io.Writer) error {
	// If no files specified, read from stdin
	if len(c.Positional) == 0 {
		return c.processReader(ctx, stdin, stdout)
	}

	var allLines []string

	// Collect lines from all input sources
	sources, err := yup.CollectInputSources(c.Positional, stdin)
	if err != nil {
		yup.ErrorF(stderr, "sort", "", err)
		return err
	}
	defer yup.CloseInputSources(sources)

	for _, source := range sources {
		// Check for cancellation before each source
		if err := yup.CheckContextCancellation(ctx); err != nil {
			return err
		}

		lines, err := c.readLines(ctx, source.Reader)
		if err != nil {
			yup.ErrorF(stderr, "sort", source.Filename, err)
			continue
		}
		allLines = append(allLines, lines...)
	}

	return c.sortAndOutput(ctx, allLines, stdout)
}

func (c command) processReader(ctx context.Context, reader io.Reader, output io.Writer) error {
	lines, err := c.readLines(ctx, reader)
	if err != nil {
		return err
	}

	return c.sortAndOutput(ctx, lines, output)
}

func (c command) readLines(ctx context.Context, reader io.Reader) ([]string, error) {
	var lines []string
	scanner := bufio.NewScanner(reader)

	for yup.ScanWithContext(ctx, scanner) {
		lines = append(lines, scanner.Text())
	}

	// Check if context was cancelled
	if err := yup.CheckContextCancellation(ctx); err != nil {
		return lines, err
	}

	return lines, scanner.Err()
}

func (c command) sortAndOutput(ctx context.Context, lines []string, output io.Writer) error {
	if len(lines) == 0 {
		return nil
	}

	// Apply sorting
	sortedLines := c.applySorting(lines)

	// Apply unique filter if requested
	if c.Flags.Unique {
		sortedLines = c.removeDuplicates(sortedLines)
	}

	// Output results
	for i, line := range sortedLines {
		// Check for cancellation periodically (every 1000 lines for efficiency)
		if i%1000 == 0 {
			if err := yup.CheckContextCancellation(ctx); err != nil {
				return err
			}
		}
		fmt.Fprintln(output, line)
	}

	return nil
}

func (c command) applySorting(lines []string) []string {
	// Make a copy to avoid modifying original
	result := make([]string, len(lines))
	copy(result, lines)

	if c.Flags.Random {
		// Simple random shuffle (not cryptographically secure)
		for i := len(result) - 1; i > 0; i-- {
			j := i % (i + 1) // Simple pseudo-random
			result[i], result[j] = result[j], result[i]
		}
		return result
	}

	sort.Slice(result, func(i, j int) bool {
		return c.comparLines(result[i], result[j])
	})

	return result
}

func (c command) comparLines(line1, line2 string) bool {
	key1 := c.extractSortKey(line1)
	key2 := c.extractSortKey(line2)

	var less bool

	if c.Flags.Numeric {
		// Numeric comparison
		num1, err1 := strconv.ParseFloat(key1, 64)
		num2, err2 := strconv.ParseFloat(key2, 64)

		if err1 != nil && err2 != nil {
			// Both non-numeric, fall back to string comparison
			less = c.stringCompare(key1, key2)
		} else if err1 != nil {
			// key1 is non-numeric, key2 is numeric
			less = false
		} else if err2 != nil {
			// key1 is numeric, key2 is non-numeric
			less = true
		} else {
			// Both numeric
			less = num1 < num2
		}
	} else {
		// String comparison
		less = c.stringCompare(key1, key2)
	}

	// Apply reverse if requested
	if c.Flags.Reverse {
		less = !less
	}

	return less
}

func (c command) stringCompare(s1, s2 string) bool {
	if c.Flags.IgnoreCase {
		return strings.ToLower(s1) < strings.ToLower(s2)
	}
	return s1 < s2
}

func (c command) extractSortKey(line string) string {
	if c.Flags.Field == 0 {
		// Sort by whole line
		return line
	}

	// Split by delimiter and extract field
	var fields []string
	if c.Flags.Delimiter == " " {
		// Special case for whitespace - split on any whitespace
		fields = strings.Fields(line)
	} else {
		fields = strings.Split(line, string(c.Flags.Delimiter))
	}

	// Field numbers are 1-based
	fieldIndex := int(c.Flags.Field) - 1
	if fieldIndex >= 0 && fieldIndex < len(fields) {
		return fields[fieldIndex]
	}

	// Field doesn't exist, return empty string
	return ""
}

func (c command) removeDuplicates(lines []string) []string {
	seen := make(map[string]bool)
	var result []string

	for _, line := range lines {
		key := line
		if c.Flags.IgnoreCase {
			key = strings.ToLower(line)
		}

		if !seen[key] {
			seen[key] = true
			result = append(result, line)
		}
	}

	return result
}
