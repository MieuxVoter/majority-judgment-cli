# Majority Judgment CLI tool

> **WORK IN PROGRESS**
> The features described in this README are not yet implemented.
> We're merely doc-driving this tool.


## Download

> From CI artifacts?
> Perhaps we'll handcraft a few releases


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

    go get
    go build -o mj


### Build distributables

    go build -ldflags "-s -w" -o mj && upx mj

Yields a `mj` binary of about `2 Mio`.


#### For Windows

    GOOS=windows GOARCH=amd64 go build -ldflags "-s -w" -o mj.exe