// SPDX-License-Identifier: Apache-2.0

package formatter

import (
	"github.com/mieuxvoter/majority-judgment-library-go/judgment"
	"github.com/mieuxvoter/merit-profile-library-go/merit"
	"image/color"
)

// SvgMeritFormatter renders merit profiles as SVG
type SvgMeritFormatter struct{}

// Format the provided results in the gnuplot script form
func (t *SvgMeritFormatter) Format(
	_ *judgment.PollTally,
	result *judgment.PollResult,
	proposals []string,
	_ []string,
	options *Options,
) (string, error) {
	proposalsResults := result.Proposals
	if options.Sorted {
		proposalsResults = result.ProposalsSorted
	}

	svgProposals := make([]merit.Proposal, len(proposalsResults))

	for i, pr := range proposalsResults {
		svgProposal := merit.Proposal{
			Name:  proposals[proposalsResults[i].Index],
			Tally: pr.Tally.Tally,
		}
		svgProposals[i] = svgProposal
	}

	renderOptions := make([]merit.RenderOptions, 0)
	renderOptions = append(renderOptions, merit.WithBgColor(color.NRGBA{R: 13, G: 13, B: 13, A: 255}))
	if options.GreenToRed {
		renderOptions = append(renderOptions, merit.WithBestGradeOnLeft(true))
	}

	return merit.RenderLinearProfileSVG(svgProposals, renderOptions...)
}
