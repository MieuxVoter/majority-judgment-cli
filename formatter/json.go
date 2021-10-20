package formatter

import (
	"encoding/json"
	"github.com/mieuxvoter/majority-judgment-library-go/judgment"
)

// JsonFormatter formats the output as … JSON
// All on one line, feel free to make an option if you can figure out how to beautify it
type JsonFormatter struct{}

// Format the provided results
func (t *JsonFormatter) Format(
	tally *judgment.PollTally,
	result *judgment.PollResult,
	proposals []string,
	grades []string,
	options *Options,
) (string, error) {

	// JSON can ignore options.Sorted because it always sends back everything

	jsonBytes, jsonErr := json.Marshal(struct {
		Proposals []string             `json:"proposals"`
		Grades    []string             `json:"grades"`
		Tally     *judgment.PollTally  `json:"tally"`
		Result    *judgment.PollResult `json:"result"`
	}{
		Proposals: proposals,
		Grades:    grades,
		Tally:     tally,
		Result:    result,
	})

	if jsonErr != nil {
		return "", jsonErr
	}

	return string(jsonBytes), nil
}
