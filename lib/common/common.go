package common

import (
	"strings"
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

// SubSearchInput is a struct that carries the input params towards the resource-specific search functions
type SubSearchInput struct {
	RegionsVServices map[string][]string
	Config           aws.Config
	CResult          chan Result
	ParentWg         *sync.WaitGroup
	Profile          string
	Region           string
	SearchStr        string
}

// CheckAwsErrors is meant to handle AWS errors
func CheckAwsErrors(bcr, reqType string, client *aws.Client, err error) {
	if aerr, ok := err.(awserr.Error); ok {
		switch aerr.Code() {
		default:
			log.Warnf("%s: Unable to request resource of type '%s': %s",
				bcr, reqType, err.Error())
		case "AuthFailure":
			log.Warnf("%s: Unable to request resource of type '%s' You might need to enable this region, or check your credentials. Recieved error: %s",
				bcr, reqType, err.Error())
		}
	} else {
		log.Warnf("%s: Unable to request resource of type '%s': %s",
			bcr, reqType, err.Error())
	}
}

// BreadCrumbs Creates a breadcrumb  string for log prefixing based on a variadic number of input strings
// profilename // region // service // resource
// or
// profilename
// or
// profilename // region
func BreadCrumbs(inputs ...string) string {
	if len(inputs) < 1 {
		log.Fatalf("BreadCrumbs function recieved less than 1 string element. This should never happen. Debug the code.")
	}
	return strings.Join(inputs, " // ")
}

// ServiceInRegion - returns false if service not supported in region.
func ServiceInRegion(regionsVServices map[string][]string, profile, region, service string) bool {
	// Don't handle a search if a region doesn't support the service
	bcr := BreadCrumbs(profile, region, service)
	foundServiceInRegion := false
	for _, service := range regionsVServices[region] {
		if service == service {
			foundServiceInRegion = true
			break
		}
	}
	if !foundServiceInRegion {
		log.Warnf("%s: Region '%v' does not support this service. Skipping...", bcr, region)
		return false
	}
	return true
}
