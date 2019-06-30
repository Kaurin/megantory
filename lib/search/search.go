package search

import (
	"fmt"
	"sync"

	"github.com/Kaurin/megantory/lib/common"
	"github.com/Kaurin/megantory/lib/searchec2"
	log "github.com/sirupsen/logrus"
)

var searchStr string
var profiles []string

// regions -- Region name is key, and it contains a slice of supported services
var regionsServices map[string][]string

// Search searches all accounts and all supported services in all regions
func Search(searchStr string) {
	cResults := make(chan common.Result, 100) // Not sure how big the buffer should be
	parentWg := sync.WaitGroup{}
	searchInput := common.SearchInput{
		CResults:         cResults,
		Profiles:         detectProfiles(),
		RegionsVServices: getRegionsVServices(),
		SearchStr:        searchStr,
		ParentWg:         &parentWg,
	}
	log.Infoln("Starting the Search LIB")
	parentWg.Add(1) // "Done" signal sent at child
	go searchec2.SearchProfilesRegions(searchInput)

	for result := range cResults {
		fmt.Printf("%s // %s // %s // %s // %s\n",
			result.Account,
			result.Region,
			result.Service,
			result.ResourceType,
			result.ResourceID)
	}
	parentWg.Wait()
}
