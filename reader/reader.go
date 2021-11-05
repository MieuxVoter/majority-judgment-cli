package reader

import (
	"fmt"
	"io"
	"regexp"
	"strconv"
	"strings"
)

// Reader to implement to make another reader
type Reader interface {
	// Read the input into judgment and tally data, as well as poll metadata.
	// Most outputs are allowed to be empty (nil?),
	// but you must return at least either `tallies` or `judgments`.
	Read(
		input *io.Reader,
	) (
		judgments [][]int, // for each participant, the grade index per proposal, or -1
		tallies [][]float64, // for each proposal, the tallies of each grade
		proposals []string, // in the order they were submitted
		grades []string, // from "worst" to "best", just like in tally above
		err error,
	)
}

// sanitizeInput to help readers
func sanitizeInput(input string) string {
	sanitized := input // inefficient, but makes code below more modular — TBD

	// Remove duplicate spaces
	sanitized = regexp.MustCompile(` +`).ReplaceAllString(sanitized, " ")

	// …

	return sanitized
}

// ReadTallyRow reads a proposal tally row from strings
func ReadTallyRow(row []string, skipFirst bool) ([]float64, error) {
	tallies := make([]float64, 0, 7)
	for colIndex, gradeTally := range row {
		if skipFirst && colIndex == 0 {
			continue
		}
		gradeTallyFloat, err := ReadNumber(gradeTally)
		if err != nil {
			return nil, fmt.Errorf("failed to read `%s` as number: %s", gradeTally, err.Error())
		}
		if gradeTallyFloat < 0 {
			return nil, fmt.Errorf("strictly negative numbers are not allowed, but got `%s`", gradeTally)
		}
		tallies = append(tallies, gradeTallyFloat)
	}

	return tallies, nil
}

// ReadNamesRow reads a bunch of names as strings
func ReadNamesRow(row []string, skipFirst bool) (names []string) {
	names = make([]string, 0, 10)
	for i, name := range row {
		if skipFirst && 0 == i {
			continue
		}
		names = append(names, strings.TrimSpace(name))
	}

	return names
}

// ReadNumber reads the number from the input string.
func ReadNumber(s string) (n float64, err error) {
	s = strings.TrimSpace(s)
	if s == "" {
		return 0.0, nil
	}
	return strconv.ParseFloat(s, 64)
}

const alphabet = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"

// GenerateDummyGradeNames generates dummy grade names in reverse alphabetical order
func GenerateDummyGradeNames(thatMany int) (grades []string, err error) {
	if thatMany < 0 {
		err = fmt.Errorf("cannot generate negative amounts of grades (tried %d)", thatMany)
		return
	}
	if thatMany > len(alphabet) {
		err = fmt.Errorf("no more than %d different grades can be generated (tried %d)", len(alphabet), thatMany)
		return
	}
	grades = strings.Split(alphabet[0:thatMany], "")
	for i, j := 0, thatMany-1; i < j; i, j = i+1, j-1 {
		grades[i], grades[j] = grades[j], grades[i]
	}
	for i, grade := range grades {
		grades[i] = fmt.Sprintf("Grade %s", grade)
	}

	return
}

func readFirstRune(str string) (first rune) {
	for _, someRune := range str {
		first = someRune
		break
	}
	return
}
