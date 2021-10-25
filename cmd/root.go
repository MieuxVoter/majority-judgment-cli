/*
Copyright © 2021 Unesco

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

var cfgFile string

const ErrorConfiguring = 1
const ErrorReading = 2
const ErrorBalancing = 3
const ErrorDeliberating = 4
const ErrorFormatting = 5

var rootCmd = &cobra.Command{
	Use:     "mj FILE",
	Version: version.GitSummary,
	Short:   "Rank proposals in Majority Judgment polls",
	Long: `Resolve majority judgment polls from an input CSV.

Say you have the following tally in a CSV file named example.csv:

	     , reject, poor, fair, good, very good, excellent
	Pizza,      3,    2,    1,    4,         4,        2
	Chips,      2,    3,    0,    4,         3,        4
	Pasta,      4,    5,    1,    4,         0,        2

You could run:

	mj example.csv

or

	cat example.csv > mj -

You probably want to sort the proposals by rank, as well:

	mj example.csv --sort

Get different formats as output:

	mj example.csv --format json
	mj example.csv --format yml
	mj example.csv --format csv
	mj example.csv --format gnuplot
	mj example.csv --format gnuplot --chart opinion

Gnuplots are meant to be piped as scripts to gnuplot http://www.gnuplot.info

	mj example.csv --sort --format gnuplot | gnuplot

`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 1 {
			// Our FILE positional argument is mandatory.
			_ = cmd.Help()
			return
		}
		var outputFormatter formatter.Formatter
		format := cmd.Flags().Lookup("format").Value.String()
		defaultTo := cmd.Flags().Lookup("default").Value.String()
		outputFormatter = &formatter.TextFormatter{}
		if "text" == format {
			//outputFormatter = &formatter.TextFormatter{}
		} else if "json" == format {
			outputFormatter = &formatter.JsonFormatter{}
		} else if "csv" == format {
			outputFormatter = &formatter.CsvFormatter{}
		} else if "yml" == format || "yaml" == format {
			outputFormatter = &formatter.YamlFormatter{}
		} else if "gnuplot" == format || "plot" == format {
			chart := cmd.Flags().Lookup("chart").Value.String()
			if "merit" == chart {
				outputFormatter = &formatter.GnuplotMeritFormatter{}
			} else if "opinion" == chart {
				outputFormatter = &formatter.GnuplotOpinionFormatter{}
			} else {
				fmt.Printf("Chart `%s` is not supported.  Supported charts: merit, opinion\n", chart)
				os.Exit(ErrorConfiguring)
			}
		} else if "gnuplot-merit" == format || "gnuplot_merit" == format {
			outputFormatter = &formatter.GnuplotMeritFormatter{}
		} else if "gnuplot-opinion" == format || "gnuplot_opinion" == format {
			outputFormatter = &formatter.GnuplotOpinionFormatter{}
		} else if "svg" == format {
			panic("todo: see issue ")
		} else {
			fmt.Printf("Format `%s` is not supported.  Supported formats: text, csv, json, yaml\n", format)
			os.Exit(ErrorConfiguring)
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
			os.Exit(ErrorReading)
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
				//proposalTallyAsInts = append(proposalTallyAsInts, uint64(gradeTallyAsFloat))
				proposalTallyAsInts = append(proposalTallyAsInts, uint64(gradeTallyAsFloat*precisionScale))
			}
			proposalTally := &judgment.ProposalTally{Tally: proposalTallyAsInts}
			proposalsTallies = append(proposalsTallies, proposalTally)
		}

		poll := &judgment.PollTally{
			Proposals: proposalsTallies,
		}
		poll.GuessAmountOfJudges()

		var balancerErr error
		defaultGradeIndex := indexOf(defaultTo, grades)
		if -1 == defaultGradeIndex {
			if "majority" == defaultTo || "median" == defaultTo {
				balancerErr = poll.BalanceWithMedianDefault()
			}
			defaultGrade, defaultToErr := reader.ReadNumber(defaultTo)
			if nil != defaultToErr {
				fmt.Printf("Unrecognized --default grade `%s`.\n", defaultTo)
				os.Exit(ErrorConfiguring)
			}
			balancerErr = poll.BalanceWithStaticDefault(uint8(defaultGrade))
		} else {
			balancerErr = poll.BalanceWithStaticDefault(uint8(defaultGradeIndex))
		}
		if balancerErr != nil {
			fmt.Println("Balancing Error:", balancerErr)
			os.Exit(ErrorBalancing)
		}

		deliberator := &judgment.MajorityJudgment{}
		result, err := deliberator.Deliberate(poll)
		if err != nil {
			fmt.Println("Deliberation Error:", err)
			os.Exit(ErrorDeliberating)
		}

		desiredWidth, widthErr := strconv.Atoi(cmd.Flags().Lookup("width").Value.String())
		if widthErr != nil || desiredWidth < 0 {
			desiredWidth = 79
		}
		options := &formatter.Options{
			Sorted: cmd.Flags().Lookup("sort").Changed,
			Width:  desiredWidth,
			Scale:  precisionScale,
		}

		out, formatterErr := outputFormatter.Format(
			poll,
			result,
			proposals,
			grades,
			options,
		)
		if formatterErr != nil {
			fmt.Println("Formatter Error:", err)
			os.Exit(ErrorFormatting)
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

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.majority-judgment-cli.yaml)")
	//rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.cobra.yaml)")
	rootCmd.Flags().StringP("format", "f", "text", "desired format of the output")
	rootCmd.Flags().StringP("default", "d", "0", "default grade to use when unbalanced")
	rootCmd.Flags().StringP("width", "w", "79", "desired width, in characters")
	rootCmd.Flags().StringP("chart", "c", "merit", "one of merit, opinion")
	//rootCmd.PersistentFlags().StringVarP(&userLicense, "license", "l", "", "name of license for the project")
	//rootCmd.PersistentFlags().Bool("viper", true, "use Viper for configuration")
	rootCmd.Flags().BoolP("sort", "s", false, "sort proposals by Rank")
	rootCmd.SetVersionTemplate("{{.Version}}\n" + version.BuildDate + "\n")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		// Search config in home directory with name ".mj.yml"
		viper.AddConfigPath(home)
		viper.SetConfigType("yaml")
		viper.SetConfigName(".mj.yml")
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
