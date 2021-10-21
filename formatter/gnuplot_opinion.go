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
			row = append(row, strconv.FormatUint(proposalTally.Tally[gradeIndex], 10))
		}

		writeErr := writer.Write(row)
		if nil != writeErr {
			log.Fatal(writeErr)
		}
	}

	writer.Flush()

	gnuplotScript := `# This is a script for gnuplot http://www.gnuplot.info/
# You may pipe it into gnuplot directly like so:
# ./mj example.csv --format gnuplot --chart opinion | gnuplot -p

$tally << EOD
` + strings.TrimSpace(buffer.String()) + `
EOD
set datafile separator ","

set title "Opinion Profile"
set xlabel 'Grades'
set ylabel 'Judges'

set border 11

set key \
	top left \
	outside horizontal \
	autotitle columnhead

set xtics nomirror scale 0
set ytics out nomirror

set grid ytics lt 0 lw 1 lc rgb "#bbbbbb"

set style data histogram
set style histogram rowstacked
set style fill solid border -1
set boxwidth 0.75

nb_proposals = ` + strconv.Itoa(len(proposals)) + `

plot for [col = 2 : nb_proposals+1] \
    "$tally" using col:xticlabels(1)

`

	return gnuplotScript, nil
}
