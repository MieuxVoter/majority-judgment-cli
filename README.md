# Majority Judgment CLI tool

> **WORK IN PROGRESS**
> Although the core mechanics are here and ranking does work,
> the features described in this README are not all implemented.
> We're merely doc-driving this tool.

- [ ] Read from stdin with `-`
- [x] Read `CSV` file
- [x] `--sort`
- [x] `--format text`
- [x] `--format json`
- [ ] `--format csv`
- [ ] `--format yml`
- [ ] `--format svg`
- [ ] `--chart`


## Download

> From CI artifacts?
> Perhaps we'll handcraft a few releases

Meanwhile, grab the source and build it:

    git clone https://github.com/MieuxVoter/majority-judgment-cli

## Usage

    ./mj example.csv > results.json

or use `-` to read from `stdin`:

    cat example.csv | ./mj - > results.json

You can specify the format:

    ./mj example.csv --format csv > results.csv
    ./mj example.csv --format svg > merit.svg

And the kind of chart you want:

    ./mj example.csv --format svg --chart opinion > opinion.svg

Available charts:
- `merit_linear` (default)
- `merit_circular`
- `opinion`


## Build

Install [golang](https://golang.org/doc/install).

Example for Ubuntu: `sudo snap install go --classic`

Then go into this project directory and run:

    go get
    go build -o mj


### Build distributables

    go build -ldflags "-s -w" -o mj && upx mj

Yields a `mj` binary of about `2 Mio`.

> They say we should not `strip` go builds.


#### For Windows

    GOOS=windows GOARCH=amd64 go build -ldflags "-s -w" -o mj.exe