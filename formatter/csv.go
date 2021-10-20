package formatter

import (
	"bytes"
	"encoding/csv"
	"github.com/mieuxvoter/majority-judgment-library-go/judgment"
	"log"
	"strconv"
)

type CsvFormatter struct{}

func (t *CsvFormatter) Format(
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

	buffer := new(bytes.Buffer)
	writer := csv.NewWriter(buffer)

	if err := writer.Error(); err != nil {
		log.Fatal(err)
	}

	headersWriteErr := writer.Write([]string{
		"Rank",
		"Proposal",
		"Score",
		"MajorityGrade",
		"SecondMajorityGrade",
	})

	if nil != headersWriteErr {
		log.Fatal(headersWriteErr)
	}

	for _, proposalResult := range proposalsResults {

		writeErr := writer.Write([]string{
			strconv.Itoa(proposalResult.Rank),
			proposals[proposalResult.Index],
			proposalResult.Score,
			grades[proposalResult.Analysis.MedianGrade],
			grades[proposalResult.Analysis.SecondMedianGrade],
		})
		if nil != writeErr {
			log.Fatal(writeErr)
		}

	}

	writer.Flush() // I've also seen "defer" prefixed here.  Gotta RTFM

	return buffer.String(), nil
}
