package main

import (
	// "encoding/json"
	"bufio"
	"github.com/perlin-network/noise"
	// "github.com/perlin-network/noise/kademlia"
	"fmt"
	"io"
	"log"
	"net"
	"strings"
	"time"
	// "net/http"
)

var (
	MyLocalCluster  *LocalCluster
	MyGlobalCluster *GlobalCluster
	MyClusterNode   *ClusterNode
	AllNodesButMe   []*ClusterNode
	GossipServer    net.Listener
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

func handleConnection(c net.Conn) {
	fmt.Printf("Serving %s\n", c.RemoteAddr().String())
	c.SetReadDeadline(time.Now().Add(time.Second * 3))
	defer func() {
		c.Write([]byte("Read Timeout!\n"))
		fmt.Printf("Closing connection to %s\n", c.RemoteAddr().String())
		c.Close()
	}()
	for {
		request, err := bufio.NewReader(c).ReadString('\n')
		switch err {
		case nil:
			clientRequest := strings.TrimSpace(string(request))
			if clientRequest == ":QUIT" {
				log.Println("client requested server to close the connection so closing")
				c.Close()
				return
			} else {
				log.Println(clientRequest)
			}
		case io.EOF:
			log.Println("client closed the connection by terminating the process")
			return
		default:
			log.Printf("error: %v\n", err)
			return
		}
		result := "heyyy\n"
		c.Write([]byte(string(result)))
	}
}

func StartGossipServer() {
	var err error
	GossipServer, err = net.Listen("tcp4", fmt.Sprintf("%s:%d", *NodeInterface, *NodePort))
	if err != nil {
		panic(err)
	}
	go func() {
		for {
			c, err := GossipServer.Accept()
			if err != nil {
				fmt.Println(err)
				return
			}
			go handleConnection(c)
		}
	}()
}

func SendGossipMessage(nodeAddr string) {
	con, err := net.Dial("tcp4", nodeAddr)
	defer con.Close()
	if err != nil {
		log.Fatalln(err)
	}
	defer con.Close()

	serverReader := bufio.NewReader(con)

	switch err {
	case nil:
		if _, err = con.Write([]byte("I am alive\n")); err != nil {
			log.Printf("failed to send the client request: %v\n", err)
		}
	case io.EOF:
		log.Println("client closed the connection")
		return
	default:
		log.Printf("client error: %v\n", err)
		return
	}

	// Waiting for the server response
	serverResponse, err := serverReader.ReadString('\n')

	switch err {
	case nil:
		log.Println("Got from server:", strings.TrimSpace(serverResponse))
	case io.EOF:
		log.Println("server closed the connection")
		return
	default:
		log.Printf("server error: %v\n", err)
		return
	}
}

func BeginClustering() {
	StartGossipServer()
	log.Println("Beginning cluster discovery...")
	if *FellowNodes != "" {
		nodeList := strings.Split(*FellowNodes, ",")
		log.Println("Nodelist:", nodeList)
		for _, node := range nodeList {
			log.Println("sending messag to", node)
			SendGossipMessage(node)
		}
	} else {
		log.Println("No other nodes")
	}
}

func LeaveGossipCluster() {
	// TODO: Broadcast leave
	GossipServer.Close()
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
