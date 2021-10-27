package main

import (
	"bytes"
	"io"
	"os"
	"testing"
)

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

func BenchmarkBasicUsage(b *testing.B) {

	for _, tt := range testData {
		b.Run(tt.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				os.Args = []string{os.Args[0]}
				os.Args = append(os.Args, tt.args...)

				// Keep a backup of the real files.
				oldStdOut := os.Stdout

				// Let's use a pipe as buffer ( /!. pipes have limited buffer sizes /!. )
				pipeReader, pipeWriter, _ := os.Pipe()
				os.Stdout = pipeWriter

				main()

				outChan := make(chan string)
				// copy the output in a separate goroutine so printing can't block indefinitely
				go func() {
					var buf bytes.Buffer
					io.Copy(&buf, pipeReader)
					outChan <- buf.String()
				}()

				// back to normal state
				pipeWriter.Close()
				os.Stdout = oldStdOut // restoring the real stdout
				capturedOut := <-outChan

				// dummy usage of capturedOut, since we need the <-outChan
				for range capturedOut {
					break
				}

			}
		})
	}

}
