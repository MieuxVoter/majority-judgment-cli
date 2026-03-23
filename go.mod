module github.com/MieuxVoter/majority-judgment-cli

// I don't see any reason not to bump this from time to time as needed.
// We're not a lib so we don't have to be as low as possible, right?
// Also, 1.20 has the slices lib (slices.Reverse !), so…
go 1.17

require (
	// We use ANSI characters for color
	github.com/acarl005/stripansi v0.0.0-20180116102854-5a71ef0e047d
	// We accept CSV as input
	github.com/csimplestring/go-csv v0.0.0-20180328183906-5b8b3cd94f2c
	// The amazing Majority Judgment lib made by the goated devs of MieuxVoter
	github.com/mieuxvoter/majority-judgment-library-go v0.3.3
	// We use it to get the color profiles of the user's terminal
	github.com/muesli/termenv v0.16.0
	// Cobra & Viper are the CLI app framework we use
	github.com/spf13/cobra v1.2.1
	github.com/spf13/viper v1.9.0
	// We accept YAML as input
	gopkg.in/yaml.v3 v3.0.1
)

require (
	github.com/aymanbagabas/go-osc52/v2 v2.0.1 // indirect
	github.com/fsnotify/fsnotify v1.5.1 // indirect
	github.com/hashicorp/hcl v1.0.0 // indirect
	github.com/inconshreveable/mousetrap v1.0.0 // indirect
	github.com/lucasb-eyer/go-colorful v1.2.0 // indirect
	github.com/magiconair/properties v1.8.5 // indirect
	github.com/mattn/go-isatty v0.0.20 // indirect
	github.com/mattn/go-runewidth v0.0.13 // indirect
	github.com/mitchellh/mapstructure v1.4.2 // indirect
	github.com/pelletier/go-toml v1.9.4 // indirect
	github.com/rivo/uniseg v0.4.7 // indirect
	github.com/spf13/afero v1.6.0 // indirect
	github.com/spf13/cast v1.4.1 // indirect
	github.com/spf13/jwalterweatherman v1.1.0 // indirect
	github.com/spf13/pflag v1.0.5 // indirect
	github.com/subosito/gotenv v1.2.0 // indirect
	golang.org/x/sys v0.30.0 // indirect
	golang.org/x/text v0.3.8 // indirect
	gopkg.in/ini.v1 v1.63.2 // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
)
