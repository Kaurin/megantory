package main

import "github.com/Kaurin/megantory/cmd"

var buildHash string
var buildVersion string
var buildDate string

func main() {
	cmd.BuildHash = buildHash
	cmd.BuildVersion = buildVersion
	cmd.BuildDate = buildDate
	cmd.Execute()
}
