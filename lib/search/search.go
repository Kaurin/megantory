package search

import (
	"fmt"
	"sync"

	"github.com/Kaurin/megantory/lib/common"
	"github.com/Kaurin/megantory/lib/searchec2"
	"github.com/Kaurin/megantory/lib/searchrds"
	"github.com/aws/aws-sdk-go-v2/aws/external"
	log "github.com/sirupsen/logrus"
)

// FuncsServices contains a slice of functions, each of which handles searching a service.
type funcsServices []func(common.SubSearchInput)

// Search searches all accounts and all supported services in all regions
func Search(searchStr string) {
	cResults := make(chan common.Result, 100) // Not sure how big the buffer should be
	defer close(cResults)
	parentWg := sync.WaitGroup{}
	si := common.SearchInput{
		CResults:         cResults,
		Profiles:         detectProfiles(),
		RegionsVServices: getRegionsVServices(),
		SearchStr:        searchStr,
		ParentWg:         &parentWg,
	}

	fSvcs := funcsServices{ // IMPORTANT: fill this list with new supported services
		searchec2.Search,
		searchrds.Search,
	}

	// Invoke main loop
	parentWg.Add(1)
	go searchProfilesRegions(si, fSvcs)

	parentWg.Add(1)
	go func() {
		for result := range cResults {
			fmt.Println(common.BreadCrumbs(
				result.Account,
				result.Region,
				result.Service,
				result.ResourceType,
				result.ResourceID))
		}
		parentWg.Done()
	}()
	parentWg.Wait()
}

// SearchProfilesRegions iterates provided profiles and regions and feeds the provided chan
func searchProfilesRegions(si common.SearchInput, fSvcs funcsServices) {
	wg := sync.WaitGroup{}
	log.Infoln("Started iterating regions...")
	regions := regionsF(si.RegionsVServices)
	for _, profile := range si.Profiles {
		log.Infoln("Started iterating profiles...")
		bcp := common.BreadCrumbs(profile)
		cfg, err := external.LoadDefaultAWSConfig(
			external.WithSharedConfigProfile(profile),
		)
		if err != nil {
			log.Errorf("%s: unable to load config: %v... Skipping...", bcp, err)
			continue
		}
		log.Debugf("%s: Loaded '%s' profile...", bcp, profile)
		for _, region := range regions {
			cfg.Region = region
			ssi := common.SubSearchInput{
				Profile:          profile,
				Config:           cfg,
				Region:           region,
				CResult:          si.CResults,
				ParentWg:         &wg,
				SearchStr:        si.SearchStr,
				RegionsVServices: si.RegionsVServices,
			}
			for _, f := range fSvcs {
				wg.Add(1)
				go f(ssi)
			}
		}
	}
	wg.Wait()
	si.ParentWg.Done()
}
