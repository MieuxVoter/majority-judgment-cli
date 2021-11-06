/*
Copyright © 2021 Unescoop

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package cmd

import (
	"bufio"
	"fmt"
	"github.com/MieuxVoter/majority-judgment-cli/formatter"
	"github.com/MieuxVoter/majority-judgment-cli/reader"
	"github.com/MieuxVoter/majority-judgment-cli/version"
	"github.com/spf13/cobra"
	"io"
	"strings"

	"os"
	"strconv"

	"github.com/mieuxvoter/majority-judgment-library-go/judgment"
	"github.com/spf13/viper"
)

var configurationFilePath string

const errorConfiguring = 1
const errorReading = 2
const errorBalancing = 3
const errorDeliberating = 4
const errorFormatting = 5

var rootCmd = &cobra.Command{
	Use:     "mj FILE",
	Version: version.GitSummary,
	Short:   "Resolve and inspect Majority Judgment polls",
	Long: `Resolve Majority Judgment polls from an input CSV.

Say you have the following tally in a CSV (or TSV) file named example.csv:

         , reject, poor, fair, good, very good, excellent
    Pizza,      3,    2,    1,    4,         4,        2
    Chips,      2,    3,    0,    4,         3,        4
    Pasta,      4,    5,    1,    4,         0,        2

You could run:

	mj example.csv

or use - to read from stdin:

	cat example.csv | mj -

You probably want to sort the proposals by rank, as well:

	mj example.csv --sort

You can get different formats as output:

	mj example.csv --format json
	mj example.csv --format yml
	mj example.csv --format csv
	mj example.csv --format gnuplot
	mj example.csv --format gnuplot --chart opinion

Gnuplots are meant to be piped as scripts to gnuplot http://www.gnuplot.info

	mj example.csv --sort --format gnuplot | gnuplot

The --width parameter only applies to the default format (text).

`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 1 {
			// Our FILE positional argument is mandatory.
			_ = cmd.Help()
			return
		}
		format := cmd.Flags().Lookup("format").Value.String()
		defaultTo := cmd.Flags().Lookup("default").Value.String()
		amountOfJudgesStr := cmd.Flags().Lookup("judges").Value.String()
		chart := cmd.Flags().Lookup("chart").Value.String()
		normalize := cmd.Flags().Lookup("normalize").Changed
		colorize := !cmd.Flags().Lookup("no-color").Changed
		_, hasNoColorEnv := os.LookupEnv("NO_COLOR") // https://no-color.org/
		if hasNoColorEnv {
			colorize = false
		}

		var outputFormatter formatter.Formatter
		outputFormatter = &formatter.TextFormatter{}
		if "text" == format || "txt" == format {
			if "opinion" == chart {
				outputFormatter = &formatter.TextOpinionFormatter{}
			}
		} else if "json" == format {
			outputFormatter = &formatter.JsonFormatter{}
		} else if "csv" == format {
			outputFormatter = &formatter.CsvFormatter{}
		} else if "yml" == format || "yaml" == format {
			outputFormatter = &formatter.YamlFormatter{}
		} else if "gnuplot" == format || "plot" == format {
			if "merit" == chart {
				outputFormatter = &formatter.GnuplotMeritFormatter{}
			} else if "opinion" == chart {
				outputFormatter = &formatter.GnuplotOpinionFormatter{}
			} else {
				fmt.Printf("Chart `%s` is not supported.  Supported charts: merit, opinion\n", chart)
				os.Exit(errorConfiguring)
			}
		} else if "gnuplot-merit" == format || "gnuplot_merit" == format {
			outputFormatter = &formatter.GnuplotMeritFormatter{}
		} else if "gnuplot-opinion" == format || "gnuplot_opinion" == format {
			outputFormatter = &formatter.GnuplotOpinionFormatter{}
		} else if "svg" == format {
			panic("todo: see issue https://github.com/MieuxVoter/majority-judgment-cli/issues/11")
		} else {
			fmt.Printf("Format `%s` is not supported.  Supported formats: text, csv, json, yaml\n", format)
			os.Exit(errorConfiguring)
		}

		proposalsTallies := make([]*judgment.ProposalTally, 0, 10)

		fileParameter := strings.TrimSpace(args[0])
		var csvReader io.Reader
		if "-" == fileParameter {
			csvReader = bufio.NewReader(os.Stdin)
		} else {
			csvFile, errOpen := os.Open(fileParameter)
			if errOpen != nil {
				fmt.Println(errOpen)
			}
			// a bit nasty ; should we just defer close() and ignore err?
			defer func(csvFile *os.File) {
				errClosing := csvFile.Close()
				if errClosing != nil {
					fmt.Println(errClosing)
				}
			}(csvFile)
			csvReader = csvFile
		}

		var tallyReader reader.Reader

		tallyReader = reader.CsvTallyReader{}

		_, tallies, proposals, grades, errReader := tallyReader.Read(&csvReader)
		if errReader != nil {
			fmt.Printf("Failed to read input: " + errReader.Error() + "\n")
			os.Exit(errorReading)
		}

		if normalize {
			for proposalTallyIndex, proposalTallyAsFloats := range tallies {
				proposalTotal := 0.0
				for _, gradeTallyAsFloat := range proposalTallyAsFloats {
					proposalTotal += gradeTallyAsFloat
				}
				for gradeIndex, gradeTallyAsFloat := range proposalTallyAsFloats {
					tallies[proposalTallyIndex][gradeIndex] = gradeTallyAsFloat * 100.0 / proposalTotal
				}
			}
		}

		maximumPrecisionScale := 1000000.0
		precisionScale := 1.0
		for _, proposalTallyAsFloats := range tallies {
			for _, gradeTallyAsFloat := range proposalTallyAsFloats {
				if precisionScale >= maximumPrecisionScale {
					break
				}
				for float64(uint64(gradeTallyAsFloat*precisionScale)) != gradeTallyAsFloat*precisionScale {
					if precisionScale >= maximumPrecisionScale {
						break
					}
					precisionScale *= 10.0
				}
			}
			if precisionScale > maximumPrecisionScale {
				break
			}
		}

		for _, proposalTallyAsFloats := range tallies {
			proposalTallyAsInts := make([]uint64, 0, 7)
			for _, gradeTallyAsFloat := range proposalTallyAsFloats {
				proposalTallyAsInts = append(proposalTallyAsInts, uint64(gradeTallyAsFloat*precisionScale))
			}
			proposalTally := &judgment.ProposalTally{Tally: proposalTallyAsInts}
			proposalsTallies = append(proposalsTallies, proposalTally)
		}

		poll := &judgment.PollTally{
			Proposals: proposalsTallies,
		}

		amountOfJudges, amountOfJudgesErr := strconv.ParseInt(amountOfJudgesStr, 10, 64)
		if nil != amountOfJudgesErr || amountOfJudges < 0 {
			fmt.Printf("Unrecognized --judges amount `%s`.  "+
				"Use a positive integer, like so: --judges 42\n", amountOfJudgesStr)
			os.Exit(errorConfiguring)
		}

		if amountOfJudges > 0 {
			poll.AmountOfJudges = uint64(amountOfJudges)
		} else {
			poll.GuessAmountOfJudges()
		}

		var balancerErr error
		defaultGradeIndex := indexOf(defaultTo, grades)
		if -1 == defaultGradeIndex {
			if "majority" == defaultTo || "median" == defaultTo {
				balancerErr = poll.BalanceWithMedianDefault()
			} else {
				defaultGrade, defaultToErr := reader.ReadNumber(defaultTo)
				if nil != defaultToErr {
					fmt.Printf("Unrecognized --default grade `%s`.\n", defaultTo)
					os.Exit(errorConfiguring)
				}
				balancerErr = poll.BalanceWithStaticDefault(uint8(defaultGrade))
			}
		} else {
			balancerErr = poll.BalanceWithStaticDefault(uint8(defaultGradeIndex))
		}
		if balancerErr != nil {
			fmt.Println("Balancing Error:", balancerErr)
			os.Exit(errorBalancing)
		}

		mj := &judgment.MajorityJudgment{}
		result, deliberationErr := mj.Deliberate(poll)
		if deliberationErr != nil {
			fmt.Println("Deliberation Error:", deliberationErr)
			os.Exit(errorDeliberating)
		}

		desiredWidth, widthErr := strconv.Atoi(cmd.Flags().Lookup("width").Value.String())
		if widthErr != nil || desiredWidth < 0 {
			desiredWidth = 79
		}
		options := &formatter.Options{
			Colorized: colorize,
			Scale:     precisionScale,
			Sorted:    cmd.Flags().Lookup("sort").Changed,
			Width:     desiredWidth,
		}

		out, formatterErr := outputFormatter.Format(
			poll,
			result,
			proposals,
			grades,
			options,
		)
		if formatterErr != nil {
			fmt.Println("Formatter Error:", formatterErr)
			os.Exit(errorFormatting)
		}
		fmt.Println(out)
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}

func init() {
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	rootCmd.PersistentFlags().StringVar(&configurationFilePath, "config", "", "config file (default is $HOME/.mj.yaml)")
	//rootCmd.PersistentFlags().StringVar(&configurationFilePath, "config", "", "config file (default is $HOME/.cobra.yaml)")
	rootCmd.Flags().StringP("format", "f", "text", "desired format of the output")
	rootCmd.Flags().StringP("default", "d", "0", "default grade to use when unbalanced")
	rootCmd.Flags().StringP("width", "w", "79", "desired width, in characters")
	rootCmd.Flags().StringP("chart", "c", "merit", "one of merit, opinion")
	rootCmd.Flags().Int64P("judges", "j", 0, "amount of judges participating (overrides our guess)")
	//rootCmd.PersistentFlags().StringVarP(&userLicense, "license", "l", "", "name of license for the project")
	//rootCmd.PersistentFlags().Bool("viper", true, "use Viper for configuration")
	rootCmd.Flags().BoolP("sort", "s", false, "sort proposals by their rank")
	rootCmd.Flags().BoolP("normalize", "n", false, "normalize input to balance proposal participation")
	rootCmd.Flags().Bool("no-color", false, "do not use colors in the text outputs")
	rootCmd.SetVersionTemplate("{{.Version}}\n" + version.BuildDate + "\n")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if configurationFilePath != "" {
		// Use config file from the flag.
		viper.SetConfigFile(configurationFilePath)
	} else {
		// Find home directory.
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		// Search config in home directory with name ".mj.yaml"
		viper.AddConfigPath(home)
		viper.SetConfigType("yaml")
		viper.SetConfigName(".mj.yaml")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		_, _ = fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	}
}

// indexOf searches the data for the element, and returns its index, or -1
// Go's typing is pretty strict, hence the need for a grunt function like this.
func indexOf(element string, data []string) int {
	for k, v := range data {
		if element == v {
			return k
		}
	}
	return -1
}
