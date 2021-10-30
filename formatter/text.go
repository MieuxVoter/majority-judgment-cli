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
	if expectedWidth <= 0 {
		expectedWidth = defaultWidth
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

	chartWidth := 0
	tableWidth := 0
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

		tableWidth = measureStringLength(line)
		remainingWidth := expectedWidth - tableWidth
		chartWidth = remainingWidth

		line += makeAsciiMeritProfile(
			pollTally.Proposals[proposalResult.Index],
			chartWidth,
		)

		out += line + "\n"
	}

	legendDefinitions := make([]string, 0, 16)
	for gradeIndex, gradeName := range grades {
		legendDefinitions = append(
			legendDefinitions,
			fmt.Sprintf("%s=%s", getCharForIndex(gradeIndex), truncateString(gradeName, chartWidth-3, '…')),
		)
	}

	out += "\n"
	out += makeTextLegend("Legend:", legendDefinitions, tableWidth, expectedWidth)

	return out, nil
}

// countDigits returns 1 for 0, 1 for 5, 3 for 421, 3 for -42
func countDigits(i int) (count int) {
	if i < 0 {
		count = countDigits(i * -1)
		return
	}
	if i == 0 {
		count = 1
		return
	}
	for i > 0 {
		i = i / 10 // Euclid wuz hear
		count++
	}

	return
}

// makeTextLegend makes a legend for an ASCII chart
// `title` should be shorter than `indentation` characters
func makeTextLegend(title string, definitions []string, indentation int, maxWidth int) (legend string) {
	line := ""
	leftOnLine := maxWidth
	for i, def := range definitions {
		if i == 0 {
			line += fmt.Sprintf("%*s", indentation-1, title)
			leftOnLine -= indentation - 1
		}
		needed := measureStringLength(def) + 1
		if needed > leftOnLine && i > 0 {
			legend += line + "\n"
			line = ""
			line += strings.Repeat(" ", indentation-1)
			leftOnLine = maxWidth - indentation - 1
		}
		line += fmt.Sprintf(" %s", def)
		leftOnLine -= needed
	}
	if strings.TrimSpace(line) != "" {
		legend += line
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
		gradeChar := getCharForIndex(gradeIndex)
		ascii += strings.Repeat(
			gradeChar,
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

func getCharForIndex(gradeIndex int) string {
	const chars = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	gradeIndex = gradeIndex % len(chars)
	return chars[gradeIndex : gradeIndex+1]
}
