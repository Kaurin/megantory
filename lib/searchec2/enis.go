package searchec2

import (
	"context"
	"strings"
	"sync"

	"github.com/Kaurin/megantory/lib/common"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	log "github.com/sirupsen/logrus"
)

// searchEnis searches single AWS EC2 region
func searchEnis(ec2i ec2Input) {
	wg := sync.WaitGroup{}
	service := "ec2"
	resourceType := "ec2-eni"
	bcr := common.BreadCrumbs(ec2i.profile, ec2i.region, service, resourceType)
	cnetworkInterfaces := make(chan ec2.NetworkInterface)
	go describeEnis(ec2i.client, ec2i.profile, cnetworkInterfaces)
	for networkInterfaceL := range cnetworkInterfaces {
		networkInterface := networkInterfaceL
		wg.Add(1)
		go func() {
			networkInterfaceLower := strings.ToLower(networkInterface.String())
			searchStrLower := strings.ToLower(ec2i.searchStr)
			if strings.Contains(networkInterfaceLower, searchStrLower) {
				result := common.Result{
					Account:      ec2i.profile,
					Region:       ec2i.region,
					Service:      service,
					ResourceType: resourceType,
					ResourceID:   *networkInterface.NetworkInterfaceId,
					ResourceJSON: networkInterface.String(),
				}
				log.Debugf("%s: Matched an %s, sending back to the results channel.", bcr, resourceType)
				ec2i.cResult <- result
			}
			wg.Done()
		}()
	}
	wg.Wait()
	ec2i.parentWg.Done()
}

// describeEnis Similar to describenetworkinterfaces
func describeEnis(client ec2.Client, profile string, c chan<- ec2.NetworkInterface) {
	defer close(c)
	reqType := "ec2-networkInterface"
	service := "ec2"
	bcr := common.BreadCrumbs(profile, client.Region, service, reqType)
	input := &ec2.DescribeNetworkInterfacesInput{}
	req := client.DescribeNetworkInterfacesRequest(input)
	networkinterfaces, err := req.Send(context.TODO())
	if err != nil {
		common.CheckAwsErrors(bcr, reqType, client.Client, err)
		return
	}
	for _, networkInterface := range networkinterfaces.NetworkInterfaces {
		c <- networkInterface
	}
}
