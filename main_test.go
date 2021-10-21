package main

import (
	"os"
	"testing"
)

//func Test(t *testing.T) {
//	os.Args = []string{os.Args[0], "example/example6.csv", "--format", "gnuplot"}
//	main()
//}

//func Test2(t *testing.T) {
//	os.Args = []string{os.Args[0], "example/example6.csv", "--fo"}
//	main()
//}

var testData = []struct {
	name   string
	args   []string
	code   int
	stdout string
	stderr string
}{
	{
		name: "Basic usage, example.csv",
		args: []string{
			"example/example.csv",
		},
	},
	{
		name: "Basic usage, example1.csv",
		args: []string{
			"example/example1.csv",
		},
	},
	{
		name: "Basic usage, example2.csv",
		args: []string{
			"example/example2.csv",
		},
	},
	{
		name: "Basic usage, example3.csv",
		args: []string{
			"example/example3.csv",
		},
	},
	{
		name: "Basic usage, example4.csv",
		args: []string{
			"example/example4.csv",
		},
	},
	{
		name: "Basic usage, example5.csv",
		args: []string{
			"example/example5.csv",
		},
	},
	{
		name: "Basic usage, example6.csv",
		args: []string{
			"example/example6.csv",
		},
	},
	{
		name: "--sort usage, example.csv",
		args: []string{
			"example/example.csv",
			"--sort",
		},
	},
	{
		name: "--format gnuplot, example.csv",
		args: []string{
			"example/example.csv",
			"--format",
			"gnuplot",
		},
	},
}

func TestAll(t *testing.T) {
	for _, tt := range testData {
		t.Run(tt.name, func(t *testing.T) {

			os.Args = []string{os.Args[0]}
			os.Args = append(os.Args, tt.args...)
			main()

			// How to do?
			// Check return code against tt.code
			// Check stdout against tt.stdout
			// Check stderr against tt.stderr
		})
	}
}
