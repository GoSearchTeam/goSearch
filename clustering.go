package main

import (
	// "encoding/json"
	"context"
	"github.com/perlin-network/noise"
	// "github.com/perlin-network/noise/kademlia"
	"log"
	"net"
	"strings"
	// "net/http"
)

var (
	MyLocalCluster  *LocalCluster
	MyGlobalCluster *GlobalCluster
	MyClusterNode   *ClusterNode
	AllNodesButMe   []*ClusterNode
	TestNode        *noise.Node
)

// GlobalCluster A total cluster
type GlobalCluster struct {
	NodeCount         int // How many nodes in the global cluster
	Nodes             []ClusterNode
	ReplicationFactor int
	ID                string // ID of the global cluster
}

type LocalCluster struct {
	ID        string // ID of the local cluster
	NodeCount int    // How many nodes in local cluster
	Name      string // Name of the local cluster
}

// ClusterNode A single node within a Cluster
type ClusterNode struct {
	LocalCluster  *LocalCluster  // Parent local cluster
	GlobalCluster *GlobalCluster // Parent global cluster
	IP            string         // IP Address of the node
	Port          int            // Port of the node
	Name          string         // Name of the node
	NoiseNode     *noise.Node
}

// ClusterDiscoverResponse The response struct from discovering nodes
type ClusterDiscoverResponse struct {
	Nodes []struct {
		LocalClusterID string // Parent local cluster
		IP             string // IP of node
		Port           int    // Port of node
		Name           string // Name of node
	}
}

func BeginClusterDiscovery() {
	log.Println("Beginning cluster discovery...")
	netIP := net.ParseIP(*NodeInterface)
	TestNode, _ = noise.NewNode(noise.WithNodeBindHost(netIP), noise.WithNodeBindPort(uint16(*NodePort)))
	// overlay := kademlia.New()
	// TestNode.Bind(overlay.Protocol())
	TestNode.Listen()
	log.Println("Node listening on:", TestNode.Addr())
	TestNode.Handle(func(ctx noise.HandlerContext) error {
		log.Printf("Got a message from Bob: '%s'\n", string(ctx.Data()))
		return nil
	})
	log.Println("handle...")
	if *FellowNodes != "" {
		nodeList := strings.Split(*FellowNodes, ",")
		log.Println("Nodelist:", nodeList)
		for _, node := range nodeList {
			log.Println("sending messag to", node)
			err := TestNode.Send(context.TODO(), node, []byte("hey"))
			if err != nil {
				panic(err)
			}
		}
	} else {
		log.Println("No other nodes")
	}
}

// func initCluster() *GlobalCluster {
// 	az := "test-1"
// 	ip := "0.0.0.0"
// 	name := "test-name"
// 	rf := 2
// 	currentNode := ClusterNode{az, ip, 9889, name}
// 	nodeArr := make([]ClusterNode, 0)
// 	nodeArr = append(nodeArr, currentNode)
// 	cluster := &GlobalCluster{1, nodeArr, rf}
// 	return cluster
// }

// // =============================================================================
// // Cluster Methods
// // =============================================================================

// // DiscoverNodes Takes in a list of IP addresses and discoveres all other nodes in cluster
// func (cluster *GlobalCluster) DiscoverNodes(ipList *[]string) {
// 	// Contact ipList to see if they exist
// 	// Add each node to cluster nodes
// 	// Ask each node if they know of other nodes
// 	// Contact and add those nodes to list of they don't already exist
// 	log.Println("### DISCOVERING OTHER NODES ###")
// 	httpClient := &http.Client{Timeout: 10 * time.Second}
// 	for _, ip := range *ipList {
// 		// Contact item
// 		r, err := httpClient.Get(fmt.Sprintf("http://%s", ip))
// 		if err != nil {
// 			return
// 		}
// 		defer r.Body.Close()
// 		jsonResponse := &ClusterDiscoverResponse{}
// 		json.NewDecoder(r.Body).Decode(jsonResponse)

// 		log.Println(jsonResponse)
// 	}
// }

// // GetKnownNodes Gets the known ndoes in the Cluster, puts in response object
// func (cluster *GlobalCluster) GetKnownNodes() *ClusterDiscoverResponse {
// 	nodes := make([]ClusterNode, 0)
// 	// Build from addresses
// 	for _, i := range cluster.Nodes {
// 		nodes = append(nodes, i)
// 	}
// 	return &ClusterDiscoverResponse{nodes}
// }

// func (cluster *Cluster) addNode(r *ClusterDiscoverResponse) {
// 	for _, node := range r.Nodes {
// 		cluster.Nodes = append(cluster.Nodes, node)
// 	}
// }
