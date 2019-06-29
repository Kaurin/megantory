package main

import (
	"runtime"

	"github.com/Kaurin/megantory/cmd"
)

var buildHash string
var buildVersion string
var buildDate string

func main() {
	cmd.BuildHash = buildHash
	cmd.BuildVersion = buildVersion
	cmd.BuildDate = buildDate
	cmd.BuildGoVersion = runtime.Version()
	cmd.Execute()
}
