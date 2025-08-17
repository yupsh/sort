package command

import (
	"math/rand"
	"sort"
	"strconv"
	"strings"

	yup "github.com/gloo-foo/framework"
)

type command yup.Inputs[yup.File, flags]

func Sort(parameters ...any) yup.Command {
	// Initialize automatically opens files when T is yup.File
	cmd := command(yup.Initialize[yup.File, flags](parameters...))
	if cmd.Flags.Delimiter == "" {
		cmd.Flags.Delimiter = " "
	}
	return cmd
}

func (p command) Executor() yup.CommandExecutor {
	// Wrap the helper executor so framework handles stdin vs files automatically
	return yup.Inputs[yup.File, flags](p).Wrap(
		yup.AccumulateAndProcess(func(lines []string) []string {
			return p.sortLines(lines)
		}).Executor(),
	)
}

// sortLines performs the actual sorting logic
func (p command) sortLines(lines []string) []string {
	sorted := make([]string, len(lines))
	copy(sorted, lines)

	if bool(p.Flags.Random) {
		// Random shuffle
		rand.Shuffle(len(sorted), func(i, j int) {
			sorted[i], sorted[j] = sorted[j], sorted[i]
		})
		return sorted
	}

	// Regular sort
	sort.Slice(sorted, func(i, j int) bool {
		lineI := sorted[i]
		lineJ := sorted[j]

		// Handle ignore case
		if bool(p.Flags.IgnoreCase) {
			lineI = strings.ToLower(lineI)
			lineJ = strings.ToLower(lineJ)
		}

		// Handle ignore leading blanks
		if bool(p.Flags.IgnoreLeadingBlanks) {
			lineI = strings.TrimLeft(lineI, " \t")
			lineJ = strings.TrimLeft(lineJ, " \t")
		}

		// Field-based sorting
		if p.Flags.Field > 0 {
			delim := string(p.Flags.Delimiter)
			fieldsI := strings.Split(sorted[i], delim)
			fieldsJ := strings.Split(sorted[j], delim)

			fieldIdx := int(p.Flags.Field) - 1
			if fieldIdx < len(fieldsI) {
				lineI = fieldsI[fieldIdx]
			}
			if fieldIdx < len(fieldsJ) {
				lineJ = fieldsJ[fieldIdx]
			}
		}

		// Numeric sort
		if bool(p.Flags.Numeric) {
			numI, errI := strconv.ParseFloat(strings.TrimSpace(lineI), 64)
			numJ, errJ := strconv.ParseFloat(strings.TrimSpace(lineJ), 64)

			if errI == nil && errJ == nil {
				if bool(p.Flags.Reverse) {
					return numI > numJ
				}
				return numI < numJ
			}
		}

		// String comparison
		result := lineI < lineJ
		if bool(p.Flags.Reverse) {
			return !result
		}
		return result
	})

	// Handle unique flag
	if bool(p.Flags.Unique) {
		unique := []string{}
		seen := make(map[string]bool)
		for _, line := range sorted {
			if !seen[line] {
				unique = append(unique, line)
				seen[line] = true
			}
		}
		return unique
	}

	return sorted
}
