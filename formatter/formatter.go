package formatter

import "github.com/mieuxvoter/majority-judgment-library-go/judgment"

// Options are shared between all formatters.
// Some formatters may ignore some options.
type Options struct {
	Colorized bool
	Scale     float64 // so we can use integers internally, and display floats
	Sorted    bool
	Width     int
}

const defaultWidth = 79

// Formatter to implement to make another formatter
// Keep in mind you need to add it to the "if else if" in root command as well
type Formatter interface {
	// Format the provided results
	Format(
		pollTally *judgment.PollTally,
		result *judgment.PollResult,
		proposals []string, // in the order they were submitted
		grades []string, // from "worst" to "best"
		options *Options,
	) (string, error)
}

// measureStringLength with support for unicode (hopefully)
// Heavy-duty replacement for len(str)
func measureStringLength(str string) int {
	count := 0
	for range str {
		count++
	}
	return count
}

// truncateString safely truncates a string (hopefully)
// from https://dev.to/takakd/go-safe-truncate-string-9h0
// with some tweaks, like the suffix ; the length includes the suffix
// Supports Japanese, see Range loops https://blog.golang.org/strings
// Provide a space as rune to disable the suffix
func truncateString(str string, length int, suffix rune) string {
	if length <= 0 {
		return ""
	}

	truncated := ""
	count := 0
	for _, char := range str {
		if count >= length {
			if suffix != ' ' {
				truncated = replaceAtIndex(truncated, suffix, length-1)
			}
			break
		}
		truncated += string(char)
		count++
	}
	return truncated
}

func replaceAtIndex(in string, r rune, i int) string {
	out := []rune(in)
	out[i] = r
	return string(out)
}
