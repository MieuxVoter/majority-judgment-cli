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

	biggestRank := 1
	for _, proposalResult := range proposalsResults {
		if biggestRank < proposalResult.Rank {
			biggestRank = proposalResult.Rank
		}
	}
	//amountOfDigitsForRank := countDigits(biggestRank)

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

	amountOfCharactersForGrade := 1
	maximumAmountOfCharactersForGrade := 20
	for _, grade := range grades {
		thatGradeLength := measureStringLength(grade)
		if thatGradeLength > amountOfCharactersForGrade {
			amountOfCharactersForGrade = thatGradeLength
		}
	}
	if amountOfCharactersForGrade > maximumAmountOfCharactersForGrade {
		amountOfCharactersForGrade = maximumAmountOfCharactersForGrade
	}

	maximumAmountOfJudgmentsForGrade := uint64(0)
	for gradeIndex := range grades {
		cumulatedAmountOfJudgmentsForGrade := uint64(0)
		for _, proposalTally := range proposalsTallies {
			cumulatedAmountOfJudgmentsForGrade += proposalTally.Tally[gradeIndex]
		}
		if cumulatedAmountOfJudgmentsForGrade > maximumAmountOfJudgmentsForGrade {
			maximumAmountOfJudgmentsForGrade = cumulatedAmountOfJudgmentsForGrade
		}
	}
	amountOfCharactersForTotal := countDigits(maximumAmountOfCharactersForGrade)

	for gradeIndex, gradeName := range grades {

		cumulatedAmountOfJudgmentsForGrade := uint64(0)
		for _, proposalTally := range proposalsTallies {
			cumulatedAmountOfJudgmentsForGrade += proposalTally.Tally[gradeIndex]
		}

		line := fmt.Sprintf("%*d ", amountOfCharactersForTotal, cumulatedAmountOfJudgmentsForGrade)

		line += fmt.Sprintf("%*s ", amountOfCharactersForGrade, truncateString(
			gradeName,
			amountOfCharactersForGrade,
			'…',
		))

		remainingWidth := expectedWidth - measureStringLength(line)

		line += makeAsciiOpinionProfile(
			proposalsTallies,
			gradeIndex,
			maximumAmountOfJudgmentsForGrade,
			remainingWidth,
		)

		out += line + "\n"
	}

	out += "\n   Legend:  "
	for proposalIndex, proposalResult := range proposalsResults {
		if proposalIndex > 0 {
			out += "  "
		}
		out += fmt.Sprintf(
			"%s=%s",
			getCharForIndex(proposalResult.Index),
			proposals[proposalResult.Index],
		)
	}

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
