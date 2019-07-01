package searchec2

import (
	"context"
	"strings"

	"github.com/Kaurin/megantory/lib/common"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	log "github.com/sirupsen/logrus"
)

// searchAddresses searches single AWS EC2 region
func searchAddresses(ec2i ec2Input) {
	cAddresses := make(chan ec2.Address)
	go describeAddresses(ec2i.client, ec2i.profile, cAddresses)
	for address := range cAddresses {
		addressLower := strings.ToLower(address.String())
		searchStrLower := strings.ToLower(ec2i.searchStr)
		if strings.Contains(addressLower, searchStrLower) {
			result := common.Result{
				Account:      ec2i.profile,
				Region:       ec2i.region,
				Service:      "ec2",
				ResourceType: "ec2-address",
				ResourceID:   *address.AllocationId,
				ResourceJSON: address.String(),
			}
			log.Debugln("EC2: Matched an address, sending back to the results channel.")
			ec2i.cResult <- result
		}
	}
	ec2i.parentWg.Done()
}

// describeAddresses Similar to describeInstances, but without pagination.
func describeAddresses(client ec2.Client, profile string, c chan<- ec2.Address) {
	defer close(c)
	reqType := "ec2-address"
	service := "ec2"
	bcr := common.BreadCrumbs(profile, client.Region, service, reqType)
	input := &ec2.DescribeAddressesInput{}
	req := client.DescribeAddressesRequest(input)
	addresses, err := req.Send(context.TODO())
	if err != nil {
		common.CheckAwsErrors(bcr, reqType, client.Client, err)
		return
	}
	for _, address := range addresses.Addresses {
		c <- address
	}
}
