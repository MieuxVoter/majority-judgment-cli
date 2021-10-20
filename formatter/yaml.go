package formatter

import (
	"github.com/mieuxvoter/majority-judgment-library-go/judgment"
	"gopkg.in/yaml.v3"
)

// YamlFormatter formats the results in YAML
type YamlFormatter struct{}

// Format the provided results
func (t *YamlFormatter) Format(
	tally *judgment.PollTally,
	result *judgment.PollResult,
	proposals []string,
	grades []string,
	options *Options,
) (string, error) {

	// Can ignore options.Sorted because it always sends back everything

	yamlBytes, yamlErr := yaml.Marshal(struct {
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

	if yamlErr != nil {
		return "", yamlErr
	}

	return string(yamlBytes), nil
}
