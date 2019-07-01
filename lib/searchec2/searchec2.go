package searchec2

import (
	"sync"

	"github.com/Kaurin/megantory/lib/common"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
)

type funcResources []func(ec2Input)

type ec2Input struct {
	client    ec2.Client
	parentWg  *sync.WaitGroup
	profile   string
	searchStr string
	region    string
	cResult   chan common.Result
}

// Search through all supported EC2 resource types
func Search(ssi common.SubSearchInput) {
	if !common.ServiceInRegion(ssi.RegionsVServices, ssi.Profile, ssi.Region, "ec2") {
		ssi.ParentWg.Done()
		return
	}
	wg := sync.WaitGroup{}
	fResources := funcResources{
		searchInstances,
		searchAddresses,
	}
	client := ec2.New(ssi.Config)
	for _, f := range fResources {
		ec2i := ec2Input{
			client:    *client,
			parentWg:  &wg,
			profile:   ssi.Profile,
			searchStr: ssi.SearchStr,
			region:    ssi.Region,
			cResult:   ssi.CResult,
		}
		wg.Add(1)
		go f(ec2i)
	}
	wg.Wait()
	ssi.ParentWg.Done()
}
