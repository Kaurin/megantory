package searchec2

import (
	"context"
	"strings"
	"sync"

	"github.com/Kaurin/megantory/lib/common"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	log "github.com/sirupsen/logrus"
)

// searchAddresses searches single AWS EC2 region
func searchAddresses(client *ec2.Client, cResult chan<- common.Result, wg *sync.WaitGroup, profile, searchStr string) {
	cAddresses := make(chan *ec2.Address)
	go describeAddresses(client, profile, cAddresses)
	for address := range cAddresses {
		addressLower := strings.ToLower(address.String())
		searchStrLower := strings.ToLower(searchStr)
		if strings.Contains(addressLower, searchStrLower) {
			result := common.Result{
				Account:      profile,
				Region:       client.Region,
				Service:      "ec2",
				ResourceType: "ec2-address",
				ResourceID:   *address.AllocationId,
				ResourceJSON: address.String(),
			}
			log.Debugln("EC2: Matched an address, sending back to the results channel.")
			cResult <- result
		}
	}
	wg.Done()
}

// describeAddresses Similar to describeInstances, but without pagination.
func describeAddresses(client *ec2.Client, profile string, c chan<- *ec2.Address) {
	defer close(c)
	reqType := "address"
	input := &ec2.DescribeAddressesInput{}
	req := client.DescribeAddressesRequest(input)
	addresses, err := req.Send(context.TODO())
	if err != nil {
		common.CheckAwsErrors(profile, reqType, client.Client, err)
		return
	}

	for _, address := range addresses.Addresses {
		c <- &address
	}
}
