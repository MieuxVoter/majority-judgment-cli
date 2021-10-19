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
	"encoding/csv"
	"fmt"
	"github.com/spf13/cobra"
	"strings"

	//"log"
	"os"
	"strconv"
	//"unicode"

	"github.com/mieuxvoter/majority-judgment-library-go/judgment"
	"github.com/spf13/viper"
	//gointtoletters "github.com/arturwwl/gointtoletters"
)

var cfgFile string

var rootCmd = &cobra.Command{
	Use:   "majority-judgment-cli",
	Short: "Resolve Majority Judgment polls",
	Long: `Resolve majority judgment polls from an input CSV.

Say you have the following tally in a CSV file named example.csv:

	, reject, poor, fair, good, very good, excellent
	Candidate A, 4, 2, 3, 4, 5, 2
	Candidate B, 4, 2, 3, 4, 5, 2
	Candidate C, 4, 2, 3, 4, 5, 2

You could run:

	mj example.csv

or

	cat example.csv > mj -

Get different formats as output:

	mj example.csv --format json
	mj example.csv --format svg
	mj example.csv --format csv

`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	Run: func(cmd *cobra.Command, args []string) {
		//_, _ = fmt.Fprintln(os.Stdout, "Reading input CSV file…")

		for i, s := range args {
			fmt.Println(i, s)
		}

		//fmt.Println("format:", cmd.Flags().Lookup("format").Value)

		proposalsTallies := make([]*judgment.ProposalTally, 0, 10)

		csvFile, err := os.Open(args[0])
		if err != nil {
			fmt.Println(err)
		}
		//fmt.Println("Successfully Opened CSV file")
		defer csvFile.Close()

		csvRows, err := csv.NewReader(csvFile).ReadAll()
		if err != nil {
			fmt.Println(err)
		}
		grades := []string{}
		headerLooksLikeTally := true
		for rowIndex, row := range csvRows {

			if len(row) < 2 {
				continue
			}

			if rowIndex == 0 {
				grades = row[1:]
				for i := 0; i < len(grades); i++ {
					_, err := ReadNumber(grades[i])
					if err != nil {
						headerLooksLikeTally = false
						break
					}
				}
				if headerLooksLikeTally {
					//fmt.Println("Header suspiciously looks like a tally.")
					const abc = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
					grades = []string{abc[0:len(grades)]}
				}
			}

			if rowIndex > 0 || headerLooksLikeTally {
				proposalTally := &judgment.ProposalTally{Tally: ReadRow(row)}
				proposalsTallies = append(proposalsTallies, proposalTally)
			}

			//fmt.Println(row)
		}
		fmt.Println("grades", grades)

		poll := &judgment.PollTally{
			Proposals: proposalsTallies,
		}
		deliberator := &judgment.MajorityJudgment{}
		result, err := deliberator.Deliberate(poll)
		if err != nil {
			fmt.Println("Error", err)
			os.Exit(3)
		}
		fmt.Println("result", result)
		for resultIndex, proposalResult := range result.Proposals {
			fmt.Println("Rank", proposalResult.Rank, proposalsTallies[resultIndex].Tally)
		}
	},
}

// ReadRow reads a proposal tally row from strings
func ReadRow(row []string) (tallies []uint64) {
	tallies = make([]uint64, 0, 10)
	for colIndex, gradeTally := range row {
		if colIndex == 0 {
			continue // skip proposals column — todo: make it smarter
		}
		n, err := ReadNumber(gradeTally)
		if err != nil {
			fmt.Println("Err with ReadRow", err)
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
	rootCmd.Flags().StringP("format", "f", "json", "desired format of the output")
	//rootCmd.PersistentFlags().StringVarP(&userLicense, "license", "l", "", "name of license for the project")
	//rootCmd.PersistentFlags().Bool("viper", true, "use Viper for configuration")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
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
