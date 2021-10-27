package reader

import (
	"encoding/csv"
	"errors"
	"io"
	"strings"
)

type CsvTallyReader struct{}

func (r CsvTallyReader) Read(input *io.Reader) (
	judgments [][]int,
	tallies [][]float64,
	proposals []string,
	grades []string,
	err error,
) {

	// I. Read the whole input at once.  Stream reading ?  Perhaps someday ?
	csvRows, errReader := csv.NewReader(*input).ReadAll()
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
				proposals = append(proposals, "Proposal "+Alphabet[j:j+1])
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

	// Do we still need this?  Right now we should be ok.
	//for gradeIndex, grade := range grades {
	//	grades[gradeIndex] = strings.TrimSpace(grade)
	//}

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
