package reader

import (
	"encoding/csv"
	"errors"
	"fmt"
	"github.com/csimplestring/go-csv/detector"
	"io"
	"strings"
)

// CsvTallyReader reads a poll's tally in a CSV like so:
//     Nutriscore, G, F, E, D, C, B, A
//          Pizza, 4, 2, 3, 4, 5, 4, 1
//          Chips, 5, 3, 2, 4, 4, 3, 2
type CsvTallyReader struct{}

// Read the input CSV and return as much data as we can.
// Read does not fill the `judgments` because this data is absent from the profiles.
func (r CsvTallyReader) Read(input *io.Reader) (
	judgments [][]int,
	tallies [][]float64,
	proposals []string,
	grades []string,
	err error,
) {
	csvDelimiter := ',' // default value if our detector below fails
	csvQuote := '"'

	// I. Read the whole input at once.  Tried stream reading with io.Pipe but… buffer!
	allData, _ := io.ReadAll(*input)
	readerCloneA := strings.NewReader(string(allData))
	readerCloneB := strings.NewReader(SanitizeInput(string(allData)))

	// I.a Detect the delimiter between values in the input (default is comma `,`)
	delimiterDetector := detector.New()
	delimiters := delimiterDetector.DetectDelimiter(readerCloneA, byte(csvQuote))
	if 0 < len(delimiters) {
		csvDelimiter = readFirstRune(delimiters[0])
	} else {
		csvDelimiter = ' ' // fallback on space -- use another detector instead of this
	}
	if 1 < len(delimiters) {
		err = fmt.Errorf("too many delimiters: found `%s` and `%s`", delimiters[0], delimiters[1])
		return
	}

	// I.b Read the actual CSV contents
	csvReader := csv.NewReader(readerCloneB)
	csvReader.Comma = csvDelimiter
	csvRows, errReader := csvReader.ReadAll()
	if errReader != nil {
		err = errors.New("Failed to read input CSV: " + errReader.Error())
		return
	}

	// II. Detect the shape/structure of the input file
	hasGradesNamesRow, hasProposalNamesColumn := r.detectShape(csvRows)

	// III. Read the tallies, proposals, grades
	for rowIndex, row := range csvRows {
		rowLen := len(row)
		if rowLen < 2 {
			continue
		}

		// III.a Read the grades names on the first row, or generate some if missing
		if 0 == rowIndex {
			if hasGradesNamesRow {
				grades = ReadNamesRow(row[:], hasProposalNamesColumn)
			} else {
				var errGradesGen error
				if hasProposalNamesColumn {
					grades, errGradesGen = GenerateDummyGradeNames(rowLen - 1)
				} else {
					grades, errGradesGen = GenerateDummyGradeNames(rowLen)
				}
				if nil != errGradesGen {
					err = errors.New("Failed to generate default grades names: " + errGradesGen.Error())
					return
				}
			}
		}

		if rowIndex > 0 || !hasGradesNamesRow {
			// III.b Read the proposals' names
			if hasProposalNamesColumn {
				proposals = append(proposals, strings.TrimSpace(row[0]))
			} else {
				j := len(proposals)
				proposals = append(proposals, "Proposal "+alphabet[j:j+1])
			}

			// III.c Read the actual tallies
			proposalTallyOfFloats, tallyErr := ReadTallyRow(row, hasProposalNamesColumn)
			if nil != tallyErr {
				err = errors.New("Failed to read input tally: " + tallyErr.Error())
				return
			}
			tallies = append(tallies, proposalTallyOfFloats)
		}
	}

	return
}

func (r CsvTallyReader) detectShape(rows [][]string) (hasGradesNamesRow bool, hasProposalNamesColumn bool) {
	hasGradesNamesRow = false
	hasProposalNamesColumn = false

	for rowIndex, row := range rows {
		if rowIndex == 0 {
			for i := len(row) - 1; i >= 1; i-- {
				if "" == strings.TrimSpace(row[i]) {
					continue
				}
				_, errDetection := ReadNumber(row[i])
				if errDetection != nil {
					hasGradesNamesRow = true
					break
				}
			}
		}

		if !hasGradesNamesRow || 0 != rowIndex {
			if "" == strings.TrimSpace(row[0]) {
				continue
			}
			_, errDetection := ReadNumber(row[0])
			if errDetection != nil {
				hasProposalNamesColumn = true
			}
		}

	}

	return
}
