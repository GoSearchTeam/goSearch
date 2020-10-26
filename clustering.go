package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

// Cluster A total cluster
type Cluster struct {
	NodeCount         int
	Nodes             []ClusterNode
	ReplicationFactor int
}

// ClusterNode A single node within a Cluster
type ClusterNode struct {
	AZ   string
	IP   string
	Port int
	Name string
}

// ClusterDiscoverResponse The response struct from discovering nodes
type ClusterDiscoverResponse struct {
	Nodes []ClusterNode
}

func initCluster() *Cluster {
	az := "test-1"
	ip := "0.0.0.0"
	name := "test-name"
	rf := 2
	currentNode := ClusterNode{az, ip, 9889, name}
	nodeArr := make([]ClusterNode, 0)
	nodeArr = append(nodeArr, currentNode)
	cluster := &Cluster{1, nodeArr, rf}
	return cluster
}

// =============================================================================
// Cluster Methods
// =============================================================================

// DiscoverNodes Takes in a list of IP addresses and discoveres all other nodes in cluster
func (cluster *Cluster) DiscoverNodes(ipList *[]string) {
	// Contact ipList to see if they exist
	// Add each node to cluster nodes
	// Ask each node if they know of other nodes
	// Contact and add those nodes to list of they don't already exist
	log.Println("### DISCOVERING OTHER NODES ###")
	httpClient := &http.Client{Timeout: 10 * time.Second}
	for _, ip := range *ipList {
		// Contact item
		r, err := httpClient.Get(fmt.Sprintf("http://%s", ip))
		if err != nil {
			return
		}
		defer r.Body.Close()
		jsonResponse := &ClusterDiscoverResponse{}
		json.NewDecoder(r.Body).Decode(jsonResponse)

		log.Println(jsonResponse)
	}
}

// GetKnownNodes Gets the known ndoes in the Cluster, puts in response object
func (cluster *Cluster) GetKnownNodes() *ClusterDiscoverResponse {
	nodes := make([]ClusterNode, 0)
	// Build from addresses
	for _, i := range cluster.Nodes {
		nodes = append(nodes, i)
	}
	return &ClusterDiscoverResponse{nodes}
}

func (cluster *Cluster) addNode(r *ClusterDiscoverResponse) {
	for _, node := range r.Nodes {
		cluster.Nodes = append(cluster.Nodes, node)
	}
}
