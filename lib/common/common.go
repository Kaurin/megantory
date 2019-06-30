package common

import (
	"sync"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/aws/awserr"
	log "github.com/sirupsen/logrus"
)

// Result describes 1 element of a result
type Result struct {
	Account      string
	Region       string
	Service      string
	ResourceType string
	ResourceID   string
	ResourceJSON string
}

// SearchInput is a struct that carries the input params towards the service-specific search functions
type SearchInput struct {
	Profiles         []string
	RegionsVServices map[string][]string
	CResults         chan Result
	SearchStr        string
	ParentWg         *sync.WaitGroup
}

// Regions - Provided a map of regionsServices, this returns a slice of regions
func Regions(regionsServices map[string][]string) []string {
	regions := []string{}
	for region := range regionsServices {
		regions = append(regions, region)
	}
	return regions
}

// CheckAwsErrors handles most types of errors
// Example from the searchec2 package (addresses):
//
//		// func describeAddresses(client *ec2.Client, profile string, c chan<- *ec2.Address) {
//		// 	defer close(c)
//		// 	reqType := "address"
//		// 	input := &ec2.DescribeAddressesInput{}
//		// 	req := client.DescribeAddressesRequest(input)
//		// 	addresses, err := req.Send(context.TODO())
//		// 	if err != nil {
//		// 		checkAwsErrors(profile, reqType, client.Client, err)
//		// 		return
//		// 	}
//		// 	for _, address := range addresses.Addresses {
//		// 		c <- &address
//		// 	}
//		// }
func CheckAwsErrors(profile, reqType string, client *aws.Client, err error) error {
	if aerr, ok := err.(awserr.Error); ok {
		switch aerr.Code() {
		default:
			log.Warnf("Account: '%v'. Unable to request resource of type '%s' in region '%s': %v",
				profile, reqType, client.Region, err)
			return err
		case "AuthFailure":
			log.Warnf("Account: '%v'. Unable to request resource of type '%s'. You might need to enable region %s, or check your credentials. Recieved error: %v",
				profile, reqType, client.Region, err)
			return err
		}
	} else {
		log.Warnf("Account: '%v'. Unable to request resource of type '%s' in region '%s': %v",
			profile, reqType, client.Region, err)
		return err
	}
}
