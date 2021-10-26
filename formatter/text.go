package formatter

import (
	"fmt"
	"github.com/mieuxvoter/majority-judgment-library-go/judgment"
	"math"
	"strconv"
	"strings"
)

// TextFormatter is the default formatter.
// It displays the proposals with their merit profiles and ranks.
// It does not use color (yet).  ANSI colors are appalling.
// Perhaps we can use xterm colors?
type TextFormatter struct{}

// Format the provided results
func (t *TextFormatter) Format(
	pollTally *judgment.PollTally,
	result *judgment.PollResult,
	proposals []string,
	grades []string,
	options *Options,
) (string, error) {
	out := ""

	expectedWidth := options.Width
	if 0 >= expectedWidth {
		expectedWidth = 79
	}

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
		thatProposalLength := measureStringLength(proposal)
		if thatProposalLength > amountOfCharactersForProposal {
			amountOfCharactersForProposal = thatProposalLength
		}
	}
	if amountOfCharactersForProposal > maximumAmountOfCharactersForProposal {
		amountOfCharactersForProposal = maximumAmountOfCharactersForProposal
	}

	for _, proposalResult := range proposalsResults {
		line := fmt.Sprintf(
			"#%0"+strconv.Itoa(amountOfDigitsForRank)+"d  ",
			proposalResult.Rank,
		)
		line += fmt.Sprintf(
			" %*s ",
			amountOfCharactersForProposal,
			truncateString(proposals[proposalResult.Index], amountOfCharactersForProposal, '…'),
		)

		remainingWidth := expectedWidth - measureStringLength(line)

		line += makeAsciiMeritProfile(
			pollTally.Proposals[proposalResult.Index],
			remainingWidth,
		)

		out += line + "\n"
	}

	out += "\n   Legend:  "
	for gradeIndex, gradeName := range grades {
		if 0 < gradeIndex {
			out += "  "
		}
		out += fmt.Sprintf("%s=%s", getCharForGrade(gradeIndex), gradeName)
	}
	out += "\n"

	return strings.TrimSpace(out), nil
}

func countDigits(i int) (count int) {
	for i > 0 {
		i = i / 10 // Euclid wuz hear
		count++
	}

	return
}

func makeAsciiMeritProfile(
	tally *judgment.ProposalTally,
	width int,
) (ascii string) {
	if width < 3 {
		width = 3
	}
	amountOfJudges := float64(tally.CountJudgments())
	for gradeIndex, gradeTallyInt := range tally.Tally {
		gradeTally := float64(gradeTallyInt)
		gradeRune := getCharForGrade(gradeIndex)
		ascii += strings.Repeat(
			gradeRune,
			int(math.Round(float64(width)*gradeTally/amountOfJudges)),
		)
	}

	for len(ascii) < width {
		ascii += ascii[len(ascii)-1:]
	}

	for len(ascii) > width {
		ascii = ascii[0 : len(ascii)-1]
	}

	ascii = replaceAtIndex(ascii, '|', width/2)

	return
}

func getCharForGrade(gradeIndex int) string {
	const chars = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	gradeIndex = gradeIndex % len(chars)
	return chars[gradeIndex : gradeIndex+1]
}
