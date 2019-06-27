package search

import (
	"fmt"

	"github.com/Kaurin/megantory/lib/common"
	"github.com/Kaurin/megantory/lib/searchec2"
)

var input string
var profiles []string

// regions -- Region name is key, and it contains a slice of supported services
var regionsServices map[string][]string

// Search searches all accounts and all supported services in all regions
func Search(input string) {
	searchEc2(input)
}

func searchEc2(input string) {
	cResults := make(chan common.Result)
	go searchec2.SearchProfilesRegions(profiles, regionsServices, cResults, input)
	for result := range cResults {
		fmt.Printf("%s // %s // %s\n", result.Account, result.Region, result.ResourceID)
	}
}
