package formatter

import (
	"fmt"
	"github.com/mieuxvoter/majority-judgment-library-go/judgment"
	"strconv"
	"strings"
)

type TextFormatter struct{}

func (t *TextFormatter) Format(
	pollTally *judgment.PollTally,
	result *judgment.PollResult,
	proposals []string,
	grades []string,
	options *Options,
) (string, error) {
	out := ""

	proposalsResults := result.Proposals
	if options.Sorted {
		proposalsResults = result.ProposalsSorted
	}

	biggestRank := 1
	for _, proposalResult := range proposalsResults {
		if biggestRank < proposalResult.Rank {
			biggestRank = proposalResult.Rank
		}
	}
	amountOfDigitsForRank := countDigits(biggestRank)

	amountOfCharactersForProposal := 1
	maximumAmountOfCharactersForProposal := 30
	for _, proposal := range proposals {
		thatProposalLength := len(proposal)
		if thatProposalLength > amountOfCharactersForProposal {
			amountOfCharactersForProposal = thatProposalLength
		}
	}
	if amountOfCharactersForProposal > maximumAmountOfCharactersForProposal {
		amountOfCharactersForProposal = maximumAmountOfCharactersForProposal
	}

	for _, proposalResult := range proposalsResults {
		out += fmt.Sprintf(
			"#%0"+strconv.Itoa(amountOfDigitsForRank)+"d  ",
			proposalResult.Rank,
		)
		out += fmt.Sprintf(
			" %*s",
			amountOfCharactersForProposal,
			proposals[proposalResult.Index],
		)

		out += "\n"
	}

	return strings.TrimSpace(out), nil
}

func countDigits(i int) (count int) {
	for i > 0 {
		i = i / 10 // Euclid wuz hear
		count++
	}

	return
}
