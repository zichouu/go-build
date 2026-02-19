package main

type Versions struct {
	GOOS   string
	GOARCH string
}

var Ver = []Versions{
	{GOOS: "linux", GOARCH: "amd64"},
	{GOOS: "linux", GOARCH: "arm64"},
	{GOOS: "darwin", GOARCH: "arm64"},
	{GOOS: "darwin", GOARCH: "amd64"},
	{GOOS: "android", GOARCH: "arm64"},
	{GOOS: "windows", GOARCH: "amd64"},
}
