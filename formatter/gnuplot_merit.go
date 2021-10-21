package formatter

import (
	"bytes"
	"encoding/csv"
	"github.com/mieuxvoter/majority-judgment-library-go/judgment"
	"log"
	"strconv"
	"strings"
)

// GnuplotMeritFormatter creates a script for gnuplot that displays the merit profiles
type GnuplotMeritFormatter struct{}

// Format the provided results
func (t *GnuplotMeritFormatter) Format(
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
	colHeader = append(colHeader, "Proposal \\ Grade")
	colHeader = append(colHeader, grades...)
	headersWriteErr := writer.Write(colHeader)
	if nil != headersWriteErr {
		log.Fatal(headersWriteErr)
	}

	for _, proposalResult := range proposalsResults {
		proposalTally := pollTally.Proposals[proposalResult.Index]
		row := make([]string, 0, 10)
		row = append(row, TruncateString(proposals[proposalResult.Index], 23))

		for gradeIndex, _ := range grades {
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
# ./mj example.csv --format gnuplot | gnuplot -p
$data <<EOD
` + strings.TrimSpace(buffer.String()) + `
EOD
set datafile separator ','

set term wxt \
    persist \
    size 1000, 400 \
    position 300, 200 \
    background rgb '#f0f0f0' \
    title 'Merit Profiles' \
    font ',12'

set xrange [:]
set yrange [:] reverse

set key \
    out \
    center bottom \
    horizontal \
    spacing 1.5 \
    box \
    maxrows 1 \
    width 0.8

set style fill solid 1.0

set arrow \
    from 50,-0.5 \
    to 50,2.5 \
    nohead \
    dt 2 \
    front

set format x '%.0f%%'
set xtics out 20

# set title 'Merit profile'
# set bmargin at screen 0.2

unset mouse

#stats $data using 0

nb_grades = ` + strconv.Itoa(len(grades)) + `
box_width = 0.9
array colors = ['#e63333', '#fa850a', '#e0b800', '#99c21f', '#48a948', '#338033']

plot for [col=2: nb_grades + 1] \
    $data u col: 0 : \
    ( total = sum [i=2: nb_grades + 1] column(i), \
    ( sum [i=2: col-1] column(i) / total * 100)): \
    ( sum [i=2: col  ] column(i) / total * 100) : \
    ($0 - box_width / 2.) : \
    ($0 + box_width / 2.) : \
    ytic(1) \
    with boxxyerror \
    title columnhead(col) \
    lt rgb colors[col-1]
`
	return gnuplotScript, nil
}
