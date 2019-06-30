package search

import (
	"os"

	log "github.com/sirupsen/logrus"

	"github.com/aws/aws-sdk-go-v2/aws/endpoints"
	"gopkg.in/ini.v1"
)

// detectProfiles finds all our profiles defined in ~/.aws/credentials
// It will skip profiles called "DEFAULT", but include "default"
// TODO: if AWS Credentials set in env, then override with that, I guess
func detectProfiles() []string {
	log.Debugln("Starting AWS profile detection from the credentials files")
	profs := []string{}
	homedir := os.Getenv("HOME")
	cfg, err := ini.Load(homedir + "/.aws/credentials")
	if err != nil {
		log.Fatalf("Fail to read credentials file: %v", err)
	}

	for _, sec := range cfg.Sections() {
		// Skip the DEFAULT section (not the 'default', though)
		if s := sec.Name(); s == ini.DEFAULT_SECTION {
			_, errini := sec.GetKey("aws_access_key_id")
			if errini == nil {
				log.Warnln("Due to restrictions of the ini library, please name your profile 'default' rather than 'DEFAULT'.")
			}
			continue
		}

		// If aws_access_key_id exists, it means that we want to include that profile
		_, err := sec.GetKey("aws_access_key_id")
		if err == nil {
			profs = append(profs, sec.Name())
		}
	}

	if len(profs) == 0 {
		log.Fatalf("Could not find any profiles defined in the AWS credentials file. Quitting.")
	}
	log.Infof("Found the following AWS profiles: %v", profs)
	return profs
}

// getRegionsVServices -- Keys are region names. Each key has a slice of supported services
func getRegionsVServices() map[string][]string {
	log.Debugln("Grabbing the map of service availability accross regions from the AWS SDK V2")
	resolver := endpoints.NewDefaultResolver()
	partitions := resolver.Partitions()

	regs := map[string][]string{}

	for _, p := range partitions {
		if p.ID() == "aws" { // TODO: Support govcloud and china partitions
			for id, region := range p.Regions() {
				regs[id] = []string{}
				for service := range region.Services() {
					regs[id] = append(regs[id], service)
				}
			}
		}
	}
	return regs
}
