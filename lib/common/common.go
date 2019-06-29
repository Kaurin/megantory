package common

// Result describes 1 element of a result
type Result struct {
	Account      string
	Region       string
	Service      string
	ResourceType string
	ResourceID   string
	ResourceJSON string
}

// Regions - Provided a map of regionsServices, this returns a slice of regions
func Regions(regionsServices map[string][]string) []string {
	regions := []string{}
	for region := range regionsServices {
		regions = append(regions, region)
	}
	return regions
}
