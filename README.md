# Majority Judgment CLI tool

> **WORK IN PROGRESS**
> Although the core mechanics are here and ranking does work,
> the features described in this README are not all implemented.
> We're merely doc-driving this tool, and our doc is ambitious.

- [ ] Read from stdin with `-`
- [x] Read `CSV` file
- [x] `--sort`
- [x] `--format text`
- [x] `--format json`
- [x] `--format csv`
- [ ] `--format yml`
- [ ] `--format svg`
- [ ] `--chart`


## Download

> From CI artifacts?
> Perhaps we'll handcraft a few releases

Meanwhile, grab the source and build it:

    git clone https://github.com/MieuxVoter/majority-judgment-cli

## Usage

    mj example.csv

or use `-` to read from `stdin`:

    cat example.csv | mj -

You can specify the format:

    mj example.csv --format json > results.json
    mj example.csv --format csv > results.csv
    mj example.csv --format svg > merit.svg

And the kind of chart you want:

    mj example.csv --format svg --chart opinion > opinion.svg

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

    go build -ldflags "-s -w" -o mj

Yields a `mj` binary of about `5 Mio`.

> They say we should not `strip` go builds.

You can run `upx` on the binary to reduce its size:

    upx mj


#### For Windows

    GOOS=windows GOARCH=amd64 go build -ldflags "-s -w" -o mj.exe

Packing the Windows executable with `upx` appears to trigger antivirus software.

