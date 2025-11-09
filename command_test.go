package command_test

import (
	"errors"
	"testing"

	"github.com/gloo-foo/testable/assertion"
	"github.com/gloo-foo/testable/run"
	command "github.com/yupsh/sort"
)

func TestSort_Basic(t *testing.T) {
	result := run.Command(command.Sort()).
		WithStdinLines("c", "a", "b").Run()
	assertion.NoError(t, result.Err)
	assertion.Lines(t, result.Stdout, []string{"a", "b", "c"})
}

func TestSort_Reverse(t *testing.T) {
	result := run.Command(command.Sort(command.Reverse)).
		WithStdinLines("a", "b", "c").Run()
	assertion.NoError(t, result.Err)
	assertion.Lines(t, result.Stdout, []string{"c", "b", "a"})
}

func TestSort_Numeric(t *testing.T) {
	result := run.Command(command.Sort(command.Numeric)).
		WithStdinLines("10", "2", "100").Run()
	assertion.NoError(t, result.Err)
	assertion.Lines(t, result.Stdout, []string{"2", "10", "100"})
}

func TestSort_Unique(t *testing.T) {
	result := run.Command(command.Sort(command.Unique)).
		WithStdinLines("a", "b", "a", "c").Run()
	assertion.NoError(t, result.Err)
	assertion.Lines(t, result.Stdout, []string{"a", "b", "c"})
}

func TestSort_EmptyInput(t *testing.T) {
	result := run.Quick(command.Sort())
	assertion.NoError(t, result.Err)
	assertion.Empty(t, result.Stdout)
}

func TestSort_SingleLine(t *testing.T) {
	result := run.Command(command.Sort()).
		WithStdinLines("only").Run()
	assertion.NoError(t, result.Err)
	assertion.Lines(t, result.Stdout, []string{"only"})
}

func TestSort_InputError(t *testing.T) {
	result := run.Command(command.Sort()).
		WithStdinError(errors.New("read failed")).Run()
	assertion.ErrorContains(t, result.Err, "read failed")
}

func TestSort_IgnoreCase(t *testing.T) {
	result := run.Command(command.Sort(command.IgnoreCase)).
		WithStdinLines("B", "a", "C").Run()
	assertion.NoError(t, result.Err)
	assertion.Count(t, result.Stdout, 3)
}

func TestSort_Field(t *testing.T) {
	result := run.Command(command.Sort(command.Field(2), command.Delimiter(","))).
		WithStdinLines("x,3", "y,1", "z,2").Run()
	assertion.NoError(t, result.Err)
	assertion.Count(t, result.Stdout, 3)
}

func TestSort_IgnoreLeadingBlanks(t *testing.T) {
	result := run.Command(command.Sort(command.IgnoreLeadingBlanks)).
		WithStdinLines("  c", " b", "a").Run()
	assertion.NoError(t, result.Err)
	assertion.Count(t, result.Stdout, 3)
}

func TestSort_Random(t *testing.T) {
	result := run.Command(command.Sort(command.Random)).
		WithStdinLines("a", "b", "c").Run()
	assertion.NoError(t, result.Err)
	assertion.Count(t, result.Stdout, 3)
}

func TestSort_VersionSort(t *testing.T) {
	result := run.Command(command.Sort(command.VersionSort)).
		WithStdinLines("v1.2", "v1.10", "v1.1").Run()
	assertion.NoError(t, result.Err)
	assertion.Count(t, result.Stdout, 3)
}

func TestSort_HumanNumeric(t *testing.T) {
	result := run.Command(command.Sort(command.HumanNumeric)).
		WithStdinLines("1K", "1M", "1G").Run()
	assertion.NoError(t, result.Err)
	assertion.Count(t, result.Stdout, 3)
}

func TestSort_MonthSort(t *testing.T) {
	result := run.Command(command.Sort(command.MonthSort)).
		WithStdinLines("Mar", "Jan", "Feb").Run()
	assertion.NoError(t, result.Err)
	assertion.Count(t, result.Stdout, 3)
}

func TestSort_StableSort(t *testing.T) {
	result := run.Command(command.Sort(command.StableSort)).
		WithStdinLines("a", "b", "c").Run()
	assertion.NoError(t, result.Err)
	assertion.Count(t, result.Stdout, 3)
}

func TestSort_NumericReverse(t *testing.T) {
	result := run.Command(command.Sort(command.Numeric, command.Reverse)).
		WithStdinLines("1", "10", "2").Run()
	assertion.NoError(t, result.Err)
	assertion.Lines(t, result.Stdout, []string{"10", "2", "1"})
}

func TestSort_TableDriven(t *testing.T) {
	tests := []struct {
		name     string
		input    []string
		expected []string
	}{
		{"abc", []string{"c", "b", "a"}, []string{"a", "b", "c"}},
		{"nums", []string{"3", "1", "2"}, []string{"1", "2", "3"}},
		{"same", []string{"x", "x"}, []string{"x", "x"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := run.Command(command.Sort()).
				WithStdinLines(tt.input...).Run()
			assertion.NoError(t, result.Err)
			assertion.Lines(t, result.Stdout, tt.expected)
		})
	}
}

