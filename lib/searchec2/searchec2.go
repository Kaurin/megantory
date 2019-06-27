package searchec2

import (
	"context"
	"log"
	"strings"

	"github.com/Kaurin/megantory/lib/common"

	"github.com/aws/aws-sdk-go-v2/aws/external"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
)

// SearchProfilesRegions iterates provided profiles and regions and feeds the provided chan
func SearchProfilesRegions(profiles []string, regionsServices map[string][]string, cResult chan<- common.Result, input string) {
	regions := common.Regions(regionsServices)
	cBlock := make(chan struct{})
	for _, profile := range profiles {
		cfg, err := external.LoadDefaultAWSConfig(
			external.WithSharedConfigProfile(profile),
		)
		if err != nil {
			log.Printf("unable to load config: %v... Skipping...", err)
			continue
		}
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
				continue
			}

			// Proceed with search
			cfg.Region = region
			client := ec2.New(cfg)
			go search(client, cResult, cBlock, profile, input)
		}
	}
	for i := 1; i < len(regions)*len(profiles); i++ {
		<-cBlock
	}
	close(cResult)
}

// search searches single AWS EC2 region
func search(client *ec2.Client, cResult chan<- common.Result, cBlock chan<- struct{}, profile, input string) {
	cInstances := make(chan *ec2.Instance)
	go describeInstances(client, cInstances)
	for instance := range cInstances { // Blocked until describeInstances closes chan
		instanceLower := strings.ToLower(instance.String())
		inputLower := strings.ToLower(input)
		if strings.Contains(instanceLower, inputLower) {
			result := common.Result{
				Account:    profile,
				Region:     client.Region,
				ResourceID: *instance.InstanceId,
			}
			cResult <- result
		}
	}
	cBlock <- struct{}{}
}

// describeInstances wraps ec2 pagination for DescribeInstances
func describeInstances(svc *ec2.Client, c chan<- *ec2.Instance) {
	defer close(c)
	input := ec2.DescribeInstancesInput{}
	req := svc.DescribeInstancesRequest(&input)
	p := ec2.NewDescribeInstancesPaginator(req)
	ctx := context.Background()
	for p.Next(ctx) {
		for _, runInstancesOutput := range p.CurrentPage().Reservations {
			for _, instance := range runInstancesOutput.Instances {
				// fmt.Printf("%s", instance.String())
				c <- &instance
			}
		}
	}
}
