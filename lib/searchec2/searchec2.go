package searchec2

import (
	"sync"

	"github.com/Kaurin/megantory/lib/common"
	"github.com/aws/aws-sdk-go-v2/aws/external"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	log "github.com/sirupsen/logrus"
)

type subSearchInput struct {
	client    *ec2.Client
	cResult   chan common.Result
	parentWg  *sync.WaitGroup
	profile   string
	searchStr string
}

// SearchProfilesRegions iterates provided profiles and regions and feeds the provided chan
func SearchProfilesRegions(searchInput common.SearchInput) {
	wg := sync.WaitGroup{}
	log.Infof("EC2: Started searching resources...")
	regions := common.Regions(searchInput.RegionsVServices)
	for _, profile := range searchInput.Profiles {
		cfg, err := external.LoadDefaultAWSConfig(
			external.WithSharedConfigProfile(profile),
		)
		if err != nil {
			log.Errorf("unable to load config: %v... Skipping...", err)
			continue
		}
		log.Debugf("EC2: Loaded '%s' profile...", profile)
		for _, region := range regions {
			// Don't handle a search if a region doesn't support the EC2 service
			foundServiceInRegion := false
			for _, service := range searchInput.RegionsVServices[region] {
				if service == "ec2" {
					foundServiceInRegion = true
					break
				}
			}
			if !foundServiceInRegion {
				log.Warnf("EC2: Region '%v' does not support this service. Skipping...", region)
				continue
			}
			log.Tracef("EC2: Currently searching region: %v", region)
			// Proceed with search
			cfg.Region = region
			client := ec2.New(cfg)
			subSearchInput := subSearchInput{
				client:    client,
				cResult:   searchInput.CResults,
				parentWg:  &wg,
				searchStr: searchInput.SearchStr,
			}
			wg.Add(1)
			go searchInstances(subSearchInput)
			wg.Add(1)
			go searchAddresses(subSearchInput)
		}
	}
	wg.Wait()
	searchInput.ParentWg.Done()
}
