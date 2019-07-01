package searchrds

import (
	"context"
	"strings"
	"sync"

	"github.com/Kaurin/megantory/lib/common"
	"github.com/aws/aws-sdk-go-v2/service/rds"
	log "github.com/sirupsen/logrus"
)

// searchClusters searches single AWS RDS region
func searchClusters(rdsi rdsInput) {
	wg := sync.WaitGroup{}
	service := "rds"
	resourceType := "rds-cluster"
	bcr := common.BreadCrumbs(rdsi.profile, rdsi.region, service, resourceType)
	cClusters := make(chan rds.DBCluster)
	go describeClusters(rdsi.client, cClusters)
	for clusterL := range cClusters { // Blocked until describeClusters closes chan
		cluster := clusterL
		wg.Add(1)
		go func() {
			clusterLower := strings.ToLower(cluster.String())
			searchStrLower := strings.ToLower(rdsi.searchStr)
			if strings.Contains(clusterLower, searchStrLower) {
				result := common.Result{
					Account:      rdsi.profile,
					Region:       rdsi.region,
					Service:      service,
					ResourceType: resourceType,
					ResourceID:   *cluster.DBClusterIdentifier,
					ResourceJSON: cluster.String(),
				}
				log.Debugf("%s: Matched a cluster, sending back to the results channel.", bcr)
				rdsi.cResult <- result
			}
			wg.Done()
		}()
	}
	wg.Wait()
	rdsi.parentWg.Done()
}

// describeClusters wraps rds pagination for DescribeClusters
func describeClusters(client rds.Client, c chan<- rds.DBCluster) {
	defer close(c)
	input := &rds.DescribeDBClustersInput{}
	req := client.DescribeDBClustersRequest(input)
	p := rds.NewDescribeDBClustersPaginator(req)
	for p.Next(context.TODO()) {
		for _, cluster := range p.CurrentPage().DBClusters {
			c <- cluster
		}
	}
}
