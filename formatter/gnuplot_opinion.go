package formatter

import (
	"bytes"
	"encoding/csv"
	"github.com/mieuxvoter/majority-judgment-library-go/judgment"
	"log"
	"strconv"
	"strings"
)

// GnuplotOpinionFormatter creates a script for gnuplot that shows the opinion profile
type GnuplotOpinionFormatter struct{}

// Format the provided results
func (t *GnuplotOpinionFormatter) Format(
	pollTally *judgment.PollTally,
	result *judgment.PollResult,
	proposals []string,
	grades []string,
	options *Options,
) (string, error) {
	proposalsResults := result.Proposals
	if options.Sorted {
		proposalsResults = result.ProposalsSorted
	}

	proposalsNames := proposals
	if options.Sorted {
		proposalsNames = make([]string, 0, 10)
		for _, proposalResult := range proposalsResults {
			proposalsNames = append(proposalsNames, proposals[proposalResult.Index])
		}
	}
	for i, proposalName := range proposalsNames {
		proposalsNames[i] = truncateString(proposalName, 16, '…')
	}

	buffer := new(bytes.Buffer)
	writer := csv.NewWriter(buffer)

	if err := writer.Error(); err != nil {
		log.Fatal(err)
	}

	colHeader := make([]string, 0, 10)
	colHeader = append(colHeader, "Grade \\ Proposal")
	colHeader = append(colHeader, proposalsNames...)
	headersWriteErr := writer.Write(colHeader)
	if nil != headersWriteErr {
		log.Fatal(headersWriteErr)
	}

	for gradeIndex, gradeName := range grades {
		row := make([]string, 0, 10)
		row = append(row, gradeName)

		for _, proposalResult := range proposalsResults {
			proposalTally := pollTally.Proposals[proposalResult.Index]
			row = append(row, strconv.FormatFloat(
				float64(proposalTally.Tally[gradeIndex])/options.Scale,
				'f', -1, 64))
		}

		writeErr := writer.Write(row)
		if nil != writeErr {
			log.Fatal(writeErr)
		}
	}

	writer.Flush()

	plotWidth := 400 + len(grades)*90

	gnuplotScript := `# This is a script for gnuplot http://www.gnuplot.info/
# You may pipe it into gnuplot directly like so:
# ./mj example.csv --format gnuplot --chart opinion | gnuplot -p

$tally << EOD
` + strings.TrimSpace(buffer.String()) + `
EOD
set datafile separator ","

set term wxt \
    persist \
    size ` + strconv.Itoa(plotWidth) + `, 600 \
    background rgb '#f0f0f0' \
    title 'Opinion Profile' \
    font ',14'

#set title "Opinion Profile"
#set key below height 200
#set xlabel 'Grades'
set ylabel 'Judges'

set border 11


set key samplen 2 spacing 0.85

set key \
    out \
    center bottom \
    horizontal \
    spacing 1 \
    box \
    maxrows 1 \
	autotitle \
	columnhead \
    width 0.8541

set xtics nomirror scale 0
set ytics out nomirror

set grid ytics lt 0 lw 1 lc rgb "#bbbbbb"

set style data histogram
set style histogram rowstacked
set style fill solid border -1
set boxwidth 0.8541

nb_proposals = ` + strconv.Itoa(len(proposals)) + `

plot for [col = 2 : nb_proposals+1] \
    "$tally" using col:xticlabels(1)

`

	return gnuplotScript, nil
}
