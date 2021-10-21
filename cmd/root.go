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
	"encoding/csv"
	"fmt"
	"github.com/MieuxVoter/majority-judgment-cli/formatter"
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

Get different formats as output:

	mj example.csv --format json
	mj example.csv --format svg
	mj example.csv --format csv

Only positive integers are supported in tallies.
If you used normalization and have real, floating-point values,
multiply them beforehand by a big factor like 1 000 000 000.

`,
	Run: func(cmd *cobra.Command, args []string) {
		const ABC = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"

		if len(args) != 1 {
			// Our FILE positional argument is mandatory
			_ = cmd.Help()
			return
		}
		var outputFormatter formatter.Formatter
		format := cmd.Flags().Lookup("format").Value.String()
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
				os.Exit(3)
			}
		} else if "gnuplot-merit" == format || "gnuplot_merit" == format {
			outputFormatter = &formatter.GnuplotMeritFormatter{}
		} else if "gnuplot-opinion" == format || "gnuplot_opinion" == format {
			outputFormatter = &formatter.GnuplotOpinionFormatter{}
		} else if "svg" == format {
			panic("todo")
		} else {
			fmt.Printf("Format `%s` is not supported.  Supported formats: text, csv, json, yaml\n", format)
			os.Exit(3)
		}

		proposalsTallies := make([]*judgment.ProposalTally, 0, 10)

		fileParameter := strings.TrimSpace(args[0])
		var csvReader io.Reader
		if "-" == fileParameter {
			csvReader = bufio.NewReader(os.Stdin)
		} else {
			csvFile, err := os.Open(fileParameter)
			if err != nil {
				fmt.Println(err)
			}
			defer func(csvFile *os.File) {
				err := csvFile.Close()
				if err != nil {
					fmt.Println(err)
				}
			}(csvFile)
			csvReader = csvFile
		}
		csvRows, err := csv.NewReader(csvReader).ReadAll()
		if err != nil {
			fmt.Println("Failed to read input CSV:", err)
			os.Exit(2)
		}

		var grades []string
		var proposals []string
		hasGradesNamesRow := false
		hasProposalNamesColumn := false
		for rowIndex, row := range csvRows {
			if rowIndex == 0 {
				rowLen := len(row)
				for i := 1; i < rowLen; i++ {
					if "" == strings.TrimSpace(row[i]) {
						continue
					}
					_, err := ReadNumber(row[i])
					if err != nil {
						hasGradesNamesRow = true
						break
					}
				}
			}

			if !hasGradesNamesRow || 0 != rowIndex {
				if "" == strings.TrimSpace(row[0]) {
					continue
				}
				_, err := ReadNumber(row[0])
				if err != nil {
					hasProposalNamesColumn = true
				}
			}

		}

		for rowIndex, row := range csvRows {
			rowLen := len(row)
			if rowLen < 2 {
				continue
			}

			if rowIndex == 0 {
				if hasGradesNamesRow {
					if hasProposalNamesColumn {
						grades = row[1:]
					} else {
						grades = row[:]
					}
				} else {
					if hasProposalNamesColumn {
						grades = strings.Split(ABC[0:rowLen-1], "")
					} else {
						grades = strings.Split(ABC[0:rowLen], "")
					}
				}
			}

			if rowIndex > 0 || !hasGradesNamesRow {
				if hasProposalNamesColumn {
					proposals = append(proposals, strings.TrimSpace(row[0]))
				} else {
					j := len(proposals)
					proposals = append(proposals, "Proposal "+ABC[j:j+1])
				}
				proposalTally := &judgment.ProposalTally{Tally: ReadRow(row, hasProposalNamesColumn)}
				proposalsTallies = append(proposalsTallies, proposalTally)
			}
		}

		for gradeIndex, grade := range grades {
			grades[gradeIndex] = strings.TrimSpace(grade)
		}

		poll := &judgment.PollTally{
			Proposals: proposalsTallies,
		}
		deliberator := &judgment.MajorityJudgment{}
		result, err := deliberator.Deliberate(poll)
		if err != nil {
			fmt.Println("Deliberation Error:", err)
			os.Exit(3)
		}

		desiredWidth, widthErr := strconv.Atoi(cmd.Flags().Lookup("width").Value.String())
		if widthErr != nil || desiredWidth < 0 {
			desiredWidth = 79
		}
		options := &formatter.Options{
			Sorted: bool(cmd.Flags().Lookup("sort").Changed),
			Width:  desiredWidth,
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
			os.Exit(4)
		}
		fmt.Println(out)
	},
}

// ReadRow reads a proposal tally row from strings
func ReadRow(row []string, skipFirst bool) (tallies []uint64) {
	tallies = make([]uint64, 0, 10)
	for colIndex, gradeTally := range row {
		if skipFirst && colIndex == 0 {
			continue
		}
		n, err := ReadNumber(gradeTally)
		if err != nil {
			//fmt.Println("Err with ReadRow", err)
			n = 0 // or propagate, perhaps
		}
		tallies = append(tallies, uint64(n))
	}

	return tallies
}

// ReadNumber reads the number from the input string.
func ReadNumber(s string) (n float64, err error) {
	return strconv.ParseFloat(strings.TrimSpace(s), 64)
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
	rootCmd.Flags().StringP("width", "w", "79", "desired width, in characters")
	rootCmd.Flags().StringP("chart", "c", "merit", "One of merit, opinion")
	//rootCmd.PersistentFlags().StringVarP(&userLicense, "license", "l", "", "name of license for the project")
	//rootCmd.PersistentFlags().Bool("viper", true, "use Viper for configuration")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("sort", "s", false, "sort proposals by Rank")
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

		// Search config in home directory with name ".majority-judgment-cli" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigType("yaml")
		viper.SetConfigName(".majority-judgment-cli")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		_, _ = fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	}
}
