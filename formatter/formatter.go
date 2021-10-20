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
