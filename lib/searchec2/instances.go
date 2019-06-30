package searchec2

import (
	"context"
	"strings"

	"github.com/Kaurin/megantory/lib/common"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	log "github.com/sirupsen/logrus"
)

// search searches single AWS EC2 region
func searchInstances(ssi subSearchInput) {
	cInstances := make(chan *ec2.Instance)
	go describeInstances(ssi.client, cInstances)
	for instance := range cInstances { // Blocked until describeInstances closes chan
		instanceLower := strings.ToLower(instance.String())
		searchStrLower := strings.ToLower(ssi.searchStr)
		if strings.Contains(instanceLower, searchStrLower) {
			result := common.Result{
				Account:      ssi.profile,
				Region:       ssi.client.Region,
				Service:      "ec2",
				ResourceType: "ec2-instance",
				ResourceID:   *instance.InstanceId,
				ResourceJSON: instance.String(),
			}
			log.Debugln("EC2: Matched an instance, sending back to the results channel.")
			ssi.cResult <- result
		}
	}
	ssi.parentWg.Done()
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
