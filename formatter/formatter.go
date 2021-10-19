package formatter

import "github.com/mieuxvoter/majority-judgment-library-go/judgment"

type Formatter interface {
	Format(
		pollTally *judgment.PollTally,
		result *judgment.PollResult,
		proposals []string, // in the order they were submitted
		grades []string, // from "worst" to "best"
	) (string, error)
}
