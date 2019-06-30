package searchec2

import (
	"context"
	"strings"

	"github.com/Kaurin/megantory/lib/common"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	log "github.com/sirupsen/logrus"
)

// searchInstances searches single AWS EC2 region
func searchInstances(ec2i ec2Input) {
	service := "ec2"
	resourceType := "ec2-instance"
	bcr := common.BreadCrumbs(ec2i.profile, ec2i.region, service, resourceType)
	cInstances := make(chan *ec2.Instance)
	go describeInstances(ec2i.client, cInstances)
	for instance := range cInstances { // Blocked until describeInstances closes chan
		instanceLower := strings.ToLower(instance.String())
		searchStrLower := strings.ToLower(ec2i.searchStr)
		if strings.Contains(instanceLower, searchStrLower) {
			result := common.Result{
				Account:      ec2i.profile,
				Region:       ec2i.region,
				Service:      service,
				ResourceType: resourceType,
				ResourceID:   *instance.InstanceId,
				ResourceJSON: instance.String(),
			}
			log.Debugf("%s: Matched an instance, sending back to the results channel.", bcr)
			ec2i.cResult <- result
		}
	}
	ec2i.parentWg.Done()
}

// describeInstances wraps ec2 pagination for DescribeInstances
func describeInstances(client *ec2.Client, c chan<- *ec2.Instance) {
	defer close(c)
	input := &ec2.DescribeInstancesInput{}
	req := client.DescribeInstancesRequest(input)
	p := ec2.NewDescribeInstancesPaginator(req)
	for p.Next(context.TODO()) {
		for _, runInstancesOutput := range p.CurrentPage().Reservations {
			for _, instance := range runInstancesOutput.Instances {
				c <- &instance
			}
		}
	}
}
