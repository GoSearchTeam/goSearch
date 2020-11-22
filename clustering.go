package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net"
	"strings"
	"time"
)

var (
	MyLocalCluster  *LocalClusterGroup
	MyGlobalCluster *GlobalCluster
	MyClusterNode   *ClusterNode
	AllNodes        map[string]*ClusterNode
	GossipServer    net.Listener
)

// GlobalCluster A total cluster
type GlobalCluster struct {
	NodeCount         int // How many nodes in the global cluster
	Nodes             []ClusterNode
	ReplicationFactor int
	ID                string // ID of the global cluster
}

type LocalClusterGroup struct {
	ID        string                  // ID of the local cluster
	NodeCount int                     // How many nodes in local cluster
	Name      string                  // Name of the local cluster
	Nodes     map[string]*ClusterNode // Nodes in the cluster
}

// ClusterNode A single node within a Cluster
type ClusterNode struct {
	LocalCluster  string         // Parent local cluster
	GlobalCluster *GlobalCluster // Parent global cluster
	IP            string         // IP Address of the node
	Port          int            // Port of the node
	Name          string         // Name of the node
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

type GossipMessage struct {
	Type string `json:"type"`
	Data []byte `json:"data"` // Data JSON depends on the type
}

type GossipMessageTypeHello struct { // GossipMessage.Type == "hello"
	LocalCluster string `json:"localCluster"`
	Port         int    `json:"port"`
	Name         string `json:"name"`
}

func InitMyNode() {
	AllNodes = make(map[string]*ClusterNode)
	MyClusterNode = &ClusterNode{
		LocalCluster: *LocalClusterName,
		IP:           *NodeInterface,
		Port:         *NodePort,
		Name:         fmt.Sprintf("%s:%v", *NodeInterface, *NodePort),
	}
	AllNodes[MyClusterNode.Name] = MyClusterNode
	initMap := make(map[string]*ClusterNode)
	initMap[MyClusterNode.Name] = MyClusterNode
	MyLocalCluster = &LocalClusterGroup{
		Name:      *LocalClusterName,
		NodeCount: 1,
		Nodes:     initMap,
	}
}

func addNodeToCluster(localCluster string, port int, name string, tcpAddr *net.TCPAddr) error {
	if _, ok := AllNodes[name]; ok { // TODO: Suspect node stuff
		log.Println("New node pretending to be old node")
	}
	newNode := &ClusterNode{
		LocalCluster: localCluster,
		IP:           tcpAddr.IP.String(),
		Port:         port,
		Name:         name,
	}
	AllNodes[name] = newNode
	if localCluster == MyClusterNode.LocalCluster {
		MyLocalCluster.Nodes[name] = newNode
	}
	return nil
}

func handleGossipMessage(gospMsg GossipMessage, c net.Conn) {

	switch gospMsg.Type {
	case "hello":
		log.Println("New node just waved hello!")
		var gospMsgData GossipMessageTypeHello
		err := json.Unmarshal(gospMsg.Data, &gospMsgData)
		if err != nil {
			panic(err)
		}
		clientIP, ok := c.RemoteAddr().(*net.TCPAddr)
		if !ok {
			panic("Could not get TCPAddr!")
		}
		log.Println(gospMsgData)
		err = addNodeToCluster(gospMsgData.LocalCluster, gospMsgData.Port, gospMsgData.Name, clientIP)
		if err != nil {
			panic(err)
		}
		c.Write([]byte("got it\n"))
	default:
		log.Println("Unrecognized message")
		c.Write([]byte("Unrecognized message\n"))
	}
	c.Close()
	return
}

func handleConnection(c net.Conn) {
	fmt.Printf("Serving %s\n", c.RemoteAddr().String())
	c.SetReadDeadline(time.Now().Add(time.Second * 3))
	// defer func() {
	// 	c.Write([]byte("Read Timeout!\n"))
	// 	fmt.Printf("Closing connection to %s\n", c.RemoteAddr().String())
	// 	c.Close()
	// }()

	// request, err := bufio.NewReader(c).ReadString('\n')
	// result := "heyyy\n"
	// c.Write([]byte(string(result)))
	// }

	var gospMsg GossipMessage
	decoder := json.NewDecoder(c)
	err := decoder.Decode(&gospMsg)
	if err != nil {
		log.Println("Uh oh!")
		panic(err)
	}
	// for {
	switch err {
	case nil:
		// clientRequest := strings.TrimSpace(string(request))
		// if clientRequest == ":QUIT" {
		// 	log.Println("client requested server to close the connection so closing")
		// 	c.Close()
		// 	return
		// } else {
		// 	handleGossipMessage(&clientRequest, c)
		// }
		handleGossipMessage(gospMsg, c)
	case io.EOF:
		log.Println("client closed the connection by terminating the process")
		return
	default:
		log.Printf("error: %v\n", err)
		return
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
				if !strings.Contains(err.Error(), "use of closed network connection") { // Benign error
					fmt.Println("server err:")
					fmt.Println(err)
				}
				return
			}
			go handleConnection(c)
		}
	}()
}

func SendGossipMessage(nodeAddr string) {
	con, err := net.Dial("tcp4", nodeAddr)
	if err != nil {
		log.Println("Could not connect to fellow node!")
		log.Fatalln(err)
	}

	serverReader := bufio.NewReader(con)

	switch err {
	case nil:
		data := GossipMessageTypeHello{
			LocalCluster: MyClusterNode.LocalCluster,
			Name:         MyClusterNode.Name,
			Port:         MyClusterNode.Port,
		}
		msgData, err := json.Marshal(data)
		if err != nil {
			panic(err)
		}
		msg := GossipMessage{
			Data: msgData,
			Type: "hello",
		}
		jsonData, err := json.Marshal(msg)
		if err != nil {
			panic(err)
		}
		log.Println("sending", string(jsonData))
		if _, err = con.Write([]byte(string(jsonData) + "\n")); err != nil {
			log.Printf("failed to send the client request: %v\n", err)
		}
		log.Println("newline")
		// if _, err = con.Write([]byte("\n")); err != nil {
		// 	log.Printf("failed to send the client request: %v\n", err)
		// }
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
	fmt.Println("closing client connection")
	con.Close() // Close connection
}

func BeginClustering() {
	StartGossipServer()
	log.Println("Beginning cluster discovery...")
	if *FellowNodes != "" {
		nodeList := strings.Split(*FellowNodes, ",")
		log.Println("Nodelist:", nodeList)
		for _, node := range nodeList {
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
