package searchec2

import (
	"sync"

	"github.com/Kaurin/megantory/lib/common"
	"github.com/aws/aws-sdk-go-v2/aws/external"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	log "github.com/sirupsen/logrus"
)

// SearchProfilesRegions iterates provided profiles and regions and feeds the provided chan
func SearchProfilesRegions(profiles []string, regionsServices map[string][]string, cResult chan<- common.Result, searchStr string) {
	wg := sync.WaitGroup{}
	defer close(cResult)
	log.Infof("EC2: Started searching resources...")
	regions := common.Regions(regionsServices)
	for _, profile := range profiles {
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
			for _, service := range regionsServices[region] {
				if service == "ec2" {
					foundServiceInRegion = true
					break
				}
			}
			if !foundServiceInRegion {
				log.Debugf("EC2: Region '%v' does not support this service. Skipping...", region)
				continue
			}
			log.Tracef("EC2: Currently searching region: %v", region)
			// Proceed with search
			cfg.Region = region
			client := ec2.New(cfg)
			wg.Add(1)
			go searchInstances(client, cResult, &wg, profile, searchStr)
			wg.Add(1)
			go searchAddresses(client, cResult, &wg, profile, searchStr)
		}
	}
	wg.Wait()
}
