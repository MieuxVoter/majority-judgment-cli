# Majority Judgment CLI tool

[![MIT](https://img.shields.io/github/license/MieuxVoter/majority-judgment-cli?style=for-the-badge)](LICENSE)
[![Release](https://img.shields.io/github/v/release/MieuxVoter/majority-judgment-cli?include_prereleases&style=for-the-badge)](https://github.com/MieuxVoter/majority-judgment-cli/releases)
[![Build Status](https://img.shields.io/github/workflow/status/MieuxVoter/majority-judgment-cli/Go?style=for-the-badge)](https://github.com/MieuxVoter/majority-judgment-cli/actions/workflows/go.yml)
[![Code Quality](https://img.shields.io/codefactor/grade/github/MieuxVoter/majority-judgment-cli?style=for-the-badge)](https://www.codefactor.io/repository/github/mieuxvoter/majority-judgment-cli)
[![A+](https://img.shields.io/badge/go%20report-A+-brightgreen.svg?style=for-the-badge)](https://goreportcard.com/report/github.com/mieuxvoter/majority-judgment-cli)
![LoC](https://img.shields.io/tokei/lines/github/MieuxVoter/majority-judgment-cli?style=for-the-badge)
[![Discord Chat](https://img.shields.io/discord/705322981102190593.svg?style=for-the-badge)](https://discord.gg/rAAQG9S)


## Download

Hand-made builds are provided in the [Assets of each Release](https://github.com/MieuxVoter/majority-judgment-cli/releases).


## Usage

Say you have a tally CSV like so:

	     , reject, poor, fair, good, very good, excellent
	Pizza,      3,    2,    1,    4,         4,        2
	Chips,      2,    3,    0,    4,         3,        4
	Pasta,      4,    5,    1,    4,         0,        2

You can run

    ./mj example.csv

and get

    #2   Pizza 0000000000000111111111222233333333|333333334444444444444444455555555
    #1   Chips 0000000001111111111111333333333333|333344444444444445555555555555555
    #3   Pasta 0000000000000000011111111111111111|111222233333333333333333555555555
    
    Legend:  0=reject  1=poor  2=fair  3=good  4=very good  5=excellent


You probably want to `--sort` the proposals as well:

    #1   Chips 0000000001111111111111333333333333|333344444444444445555555555555555
    #2   Pizza 0000000000000111111111222233333333|333333334444444444444444455555555
    #3   Pasta 0000000000000000011111111111111111|111222233333333333333333555555555
    
    Legend:  0=reject  1=poor  2=fair  3=good  4=very good  5=excellent

or use `-` to read from `stdin`:

    cat example.csv | mj -


### Balancing

Majority Judgment, to stay fair, requires tallies to be balanced ; **all proposals must have received the same amount of judgments**.

If your tally is not balanced, you may use a _default judgment strategy_:

    mj example.csv --default 0
    mj example.csv --default excellent
    mj example.csv --default "très bien" --judges 42
    mj example.csv --default majority
    mj example.csv --normalize

The default balancing strategy is to replace missing votes with the "worst", most conservative vote, that is `--default 0`.

### Formats

You can specify the format of the output:

    ./mj example.csv --format json > results.json
    ./mj example.csv --format csv > results.csv
    ./mj example.csv --format yml > results.yml
    ./mj example.csv --format svg > merit.svg

And even format [gnuplot](http://www.gnuplot.info/) scripts that render charts:

    ./mj example.csv --format gnuplot | gnuplot

![Linear merit profiles okf the proposals of a poll](example/screenshot_merit.png)

You can specify the kind of chart you want:

    ./mj example.csv --format gnuplot --chart opinion | gnuplot

![Opinion chart, the cumulative amounts of judgments per grade](example/screenshot_opinion.png)

Available charts:
- [x] `merit` (default)
- [x] `opinion`
- [ ] …
- [ ] a LOT more would be possible with more detailed data, per participant


## Install

Copy the binary somewhere in your `PATH`.

Or don't, and use it from anywhere.


## Build

You can also grab the source and build it:

    git clone https://github.com/MieuxVoter/majority-judgment-cli

Install [golang](https://golang.org/doc/install).

Example:
- Ubuntu: `sudo snap install go --classic`

Then go into this project directory and run:

    go get
    go build -o mj
    ./mj


### Build distributables

We have a convenience script `build.sh` that will handle version embedding from git,
using the clever `govvv`.

But basically, it's:

    go build -ldflags "-s -w" -o mj

Yields a `mj` binary of about `5 Mio`.

> They say we should not `strip` go builds.

You can run `upx` on the binary to reduce its size:

    upx mj


#### For Windows

    GOOS=windows GOARCH=amd64 go build -ldflags "-s -w" -o mj.exe

Sometimes, Go builds for Windows are [falsely detected](https://golang.org/doc/faq#virus) by antiviral software.
