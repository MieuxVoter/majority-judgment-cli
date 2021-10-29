package formatter

import (
	"fmt"
	"github.com/mieuxvoter/majority-judgment-library-go/judgment"
	"math"
	"strings"
)

// TextOpinionFormatter formats opinion profiles in ASCII
type TextOpinionFormatter struct{}

// Format the provided results
func (t *TextOpinionFormatter) Format(
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

	proposalsTallies := pollTally.Proposals
	// → sort them as well?

	amountOfCharactersForProposalName := 1
	amountOfCharactersForProposalNameThreshold := 30
	for _, proposal := range proposals {
		thatProposalLength := measureStringLength(proposal)
		if thatProposalLength > amountOfCharactersForProposalName {
			amountOfCharactersForProposalName = thatProposalLength
		}
	}
	if amountOfCharactersForProposalName > amountOfCharactersForProposalNameThreshold {
		amountOfCharactersForProposalName = amountOfCharactersForProposalNameThreshold
	}

	amountOfCharactersForGrade := 1
	amountOfCharactersForGradeThreshold := 20
	for _, grade := range grades {
		thatGradeLength := measureStringLength(grade)
		if thatGradeLength > amountOfCharactersForGrade {
			amountOfCharactersForGrade = thatGradeLength
		}
	}
	if amountOfCharactersForGrade > amountOfCharactersForGradeThreshold {
		amountOfCharactersForGrade = amountOfCharactersForGradeThreshold
	}

	maximumAmountOfJudgmentsForGrade := uint64(0)
	for gradeIndex := range grades {
		cumulatedAmountOfJudgmentsForGrade := uint64(0)
		for _, proposalTally := range proposalsTallies {
			cumulatedAmountOfJudgmentsForGrade += proposalTally.Tally[gradeIndex] // * uint64(options.Scale)
		}
		if cumulatedAmountOfJudgmentsForGrade > maximumAmountOfJudgmentsForGrade {
			maximumAmountOfJudgmentsForGrade = cumulatedAmountOfJudgmentsForGrade
		}
	}
	amountOfCharactersForTotal := countDigits(int(maximumAmountOfJudgmentsForGrade))

	chartWidth := 0
	tableWidth := 0
	for gradeIndex, gradeName := range grades {

		cumulatedAmountOfJudgmentsForGrade := uint64(0)
		for _, proposalTally := range proposalsTallies {
			cumulatedAmountOfJudgmentsForGrade += proposalTally.Tally[gradeIndex]
		}

		line := ""
		if options.Scale == 1.0 {
			line += fmt.Sprintf("%*d ", amountOfCharactersForTotal, cumulatedAmountOfJudgmentsForGrade)
		} else {
			line += fmt.Sprintf("%*.2f ", amountOfCharactersForTotal, float64(cumulatedAmountOfJudgmentsForGrade)/options.Scale)
		}

		line += fmt.Sprintf("%*s ", amountOfCharactersForGrade, truncateString(
			gradeName,
			amountOfCharactersForGrade,
			'…',
		))

		tableWidth = measureStringLength(line)
		remainingWidth := expectedWidth - tableWidth
		chartWidth = remainingWidth

		line += makeAsciiOpinionProfile(
			proposalsTallies,
			gradeIndex,
			maximumAmountOfJudgmentsForGrade,
			chartWidth,
		)

		out += line + "\n"
	}

	legendDefinitions := make([]string, 0, 16)
	for _, proposalResult := range proposalsResults {
		legendDefinitions = append(
			legendDefinitions,
			fmt.Sprintf(
				"%s=%s",
				getCharForIndex(proposalResult.Index),
				proposals[proposalResult.Index],
			),
		)
	}

	out += "\n"
	out += makeTextLegend("Legend:", legendDefinitions, tableWidth, expectedWidth)

	return out, nil
}

func makeAsciiOpinionProfile(
	tallies []*judgment.ProposalTally,
	gradeIndex int,
	maximumValue uint64,
	width int,
) (ascii string) {
	if width < 3 {
		width = 3
	}

	widthFloat := float64(width)
	maximumValueFloat := float64(maximumValue)
	for proposalIndex, proposalTally := range tallies {
		gradeTallyInt := proposalTally.Tally[gradeIndex]
		gradeTally := float64(gradeTallyInt)
		proposalChar := getCharForIndex(proposalIndex)
		ascii += strings.Repeat(
			proposalChar,
			int(math.Round(widthFloat*gradeTally/maximumValueFloat)),
		)
	}

	return
}
