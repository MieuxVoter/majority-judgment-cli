package formatter

import (
	"fmt"
	"github.com/mieuxvoter/majority-judgment-library-go/judgment"
)

type TextFormatter struct{}

func (t *TextFormatter) Format(
	pollTally *judgment.PollTally,
	result *judgment.PollResult,
	proposals []string,
	grades []string,
) (string, error) {
	out := ""

	for resultIndex, proposalResult := range result.Proposals {
		//fmt.Println("Rank", proposalResult.Rank, proposalsTallies[resultIndex].Tally)
		out += fmt.Sprintf("Rank %d", proposalResult.Rank)
		out += fmt.Sprintf(" %s", proposals[resultIndex])

		out += "\n"
	}

	return out, nil
}
