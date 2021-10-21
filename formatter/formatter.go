package formatter

import "github.com/mieuxvoter/majority-judgment-library-go/judgment"

// Options are shared between all formatters.
// Some formatters may ignore some options.
type Options struct {
	Sorted bool
	Width  int
}

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

// TruncateString safely truncates a string (hopefully)
// from https://dev.to/takakd/go-safe-truncate-string-9h0
func TruncateString(str string, length int) string {
	if length <= 0 {
		return ""
	}

	// This code cannot support Japanese
	// orgLen := len(str)
	// if orgLen <= length {
	//     return str
	// }
	// return str[:length]

	// Support Japanese
	// Ref: Range loops https://blog.golang.org/strings
	truncated := ""
	count := 0
	for _, char := range str {
		truncated += string(char)
		count++
		if count >= length {
			break
		}
	}
	return truncated
}
