package search

import (
	"fmt"

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
	profiles = detectProfiles()
	regionsServices = getAllRegions()
	log.Infoln("Starting the Search LIB")
	searchEc2(searchStr)
}

func searchEc2(searchStr string) {
	log.Infoln("Calling the EC2 search library")
	cResults := make(chan common.Result)
	go searchec2.SearchProfilesRegions(profiles, regionsServices, cResults, searchStr)
	for result := range cResults {
		fmt.Printf("%s // %s // %s // %s // %s\n",
			result.Account,
			result.Region,
			result.Service,
			result.ResourceType,
			result.ResourceID)
	}
}
