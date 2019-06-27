package search

import (
	"log"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws/endpoints"
	"gopkg.in/ini.v1"
)

func init() {
	profiles = detectProfiles()
	regionsServices = getAllRegions()
}

// detectProfiles finds all our profiles defined in ~/.aws/credentials
// It will skip profiles called "DEFAULT", but include "default"
// TODO: if AWS Credentials set in env, then override with that, I guess
func detectProfiles() []string {
	profs := []string{}
	homedir := os.Getenv("HOME")
	cfg, err := ini.Load(homedir + "/.aws/credentials")
	if err != nil {
		log.Fatalf("Fail to read credentials file: %v", err)
	}

	for _, sec := range cfg.Sections() {
		if s := sec.Name(); s == ini.DEFAULT_SECTION { // Skip the DEFAULT section (not the 'default', though)
			_, errini := sec.GetKey("aws_access_key_id") // Warn if users actually have values under "DEFAULT" instead "default"
			if errini == nil {
				log.Printf("Due to restrictions of the ini library, please name your section 'default' rather than 'DEFAULT'.")
			}
			continue
		}
		_, err := sec.GetKey("aws_access_key_id")
		if err == nil { // If aws_access_key_id exists, it means that we want to include that profile
			profs = append(profs, sec.Name())
		}
	}

	return profs
}

// getAllRegions
// Keys are region names
// Each key has a []string, which is a slice of supported services
func getAllRegions() map[string][]string {
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
