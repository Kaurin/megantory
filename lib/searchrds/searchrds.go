package searchrds

import (
	"sync"

	"github.com/Kaurin/megantory/lib/common"
	"github.com/aws/aws-sdk-go-v2/service/rds"
)

type funcResources []func(rdsInput)

type rdsInput struct {
	client    rds.Client
	parentWg  *sync.WaitGroup
	profile   string
	searchStr string
	region    string
	cResult   chan common.Result
}

// Search through all supported RDS resource types
func Search(ssi common.SubSearchInput) {
	if !common.ServiceInRegion(ssi.RegionsVServices, ssi.Profile, ssi.Region, "rds") {
		ssi.ParentWg.Done()
		return
	}
	wg := sync.WaitGroup{}
	fResources := funcResources{
		searchClusters,
	}
	client := rds.New(ssi.Config)
	for _, f := range fResources {
		rdsi := rdsInput{
			client:    *client,
			parentWg:  &wg,
			profile:   ssi.Profile,
			searchStr: ssi.SearchStr,
			region:    ssi.Region,
			cResult:   ssi.CResult,
		}
		wg.Add(1)
		go f(rdsi)
	}
	wg.Wait()
	ssi.ParentWg.Done()
}
