package formatter

import (
	"fmt"
	"github.com/acarl005/stripansi"
	"github.com/mieuxvoter/majority-judgment-library-go/judgment"
	"github.com/muesli/termenv"

	"strconv"
	"strings"
)

const minimumDefinitionLength = 7

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

	colorized := options.Colorized
	palette := judgment.CreateDefaultPalette(len(grades))
	colorProfile := termenv.ColorProfile()

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

		if chartWidth%2 == 0 {
			chartWidth--
		}

		line += makeAsciiMeritProfile(
			pollTally.Proposals[proposalResult.Index],
			chartWidth,
			colorized,
		)

		out += line + "\n"
	}

	legendDefinitions := make([]string, 0, 16)
	for gradeIndex, gradeName := range grades {
		maximumDefinitionLength := chartWidth - 2
		if maximumDefinitionLength < minimumDefinitionLength {
			maximumDefinitionLength = minimumDefinitionLength
		}
		gradeChar := getCharForIndex(gradeIndex)
		if colorized {
			color := colorProfile.FromColor(palette[gradeIndex])
			s := termenv.String(gradeChar)
			s = s.Background(color)
			s = s.Foreground(color)
			gradeChar = s.String()
		}
		legendDefinitions = append(
			legendDefinitions,
			fmt.Sprintf(
				"%s=%s",
				gradeChar,
				truncateString(gradeName, maximumDefinitionLength, '…'),
			),
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
func makeTextLegend(
	title string,
	definitions []string,
	indentation int,
	maxWidth int,
) (legend string) {
	line := ""
	leftOnLine := maxWidth
	for i, def := range definitions {
		if i == 0 {
			line += fmt.Sprintf("%*s", indentation-1, title)
			leftOnLine -= indentation - 1
		}
		needed := measureStringLength(stripansi.Strip(def)) + 1
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
	colorized bool,
) (ascii string) {
	if width < 3 {
		width = 3
	}
	palette := judgment.CreateDefaultPalette(int(tally.CountAvailableGrades()))
	colorProfile := termenv.ColorProfile()

	for cursor := 0; cursor < width; cursor++ {
		ratio := float64(cursor) / float64(width)
		gradeIndex, _ := getGradeAtRatio(tally, ratio)
		gradeChar := getCharForIndex(gradeIndex)
		isMedian := (width)/2 == cursor
		if isMedian {
			gradeChar = "|"
		}

		if colorized {
			color := colorProfile.FromColor(palette[gradeIndex])
			s := termenv.String(gradeChar)
			s = s.Background(color)
			if !isMedian {
				s = s.Foreground(color)
			}
			gradeChar = s.String()
		}
		ascii += gradeChar
	}

	return
}

func getGradeAtRatio(
	tally *judgment.ProposalTally,
	ratio float64,
) (int, error) {
	targetIndex := int(ratio * float64(tally.CountJudgments()))
	cursorStart := 0
	cursor := 0
	for gradeIndex, gradeTallyInt := range tally.Tally {
		if 0 == gradeTallyInt {
			continue
		}
		cursorStart = cursor
		cursor = cursor + int(gradeTallyInt)

		if cursorStart <= targetIndex && targetIndex < cursor {
			return gradeIndex, nil
		}
	}

	return 0, fmt.Errorf("")
}

func getCharForIndex(gradeIndex int) string {
	const chars = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	gradeIndex = gradeIndex % len(chars)
	return chars[gradeIndex : gradeIndex+1]
}
