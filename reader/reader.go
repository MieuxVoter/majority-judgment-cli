package reader

import (
	"io"
)

// Reader to implement to make another reader
type Reader interface {
	// Read the input into judgment and tally data, as well as poll metadata.
	// Most outputs are allowed to be empty (nil?),
	// but you must return at least either `tallies` or `judgments`.
	Read(
		input *io.Reader,
	) (
		judgments [][]int, // for each participant, the grade index per proposal, or -1
		tallies [][]float64, // for each proposal, the tallies of each grade
		proposals []string, // in the order they were submitted
		grades []string, // from "worst" to "best", just like in tally above
		err error,
	)
}
