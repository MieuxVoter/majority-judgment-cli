package version

// These properties are filled through -ldflags upon building using govvv

// GitSummary is a long, descriptive git version like 0.3.1 or 0.3.1-12-g3257b77
var GitSummary string

// BuildDate looks like this 2021-10-20T12:24:58Z
var BuildDate string
