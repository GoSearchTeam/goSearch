package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"math/rand"
	"net"
	"net/http"
	"strings"
	"time"
)

var (
	MyLocalCluster  *LocalClusterGroup
	MyGlobalCluster *GlobalCluster
	MyClusterNode   *ClusterNode
	// AllNodes Name: Node
	AllNodes     map[string]*ClusterNode
	GossipServer net.Listener
	GMCache      *GossipMessageCache
)

// GossipMessageCache used to prevent duplicate action of gossip message
type GossipMessageCache map[uint64]void

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
	LocalCluster  string // Parent local cluster
	GlobalCluster string // Parent global cluster
	IP            string // IP Address of the node
	Port          int    // Port of the node
	Name          string // Name of the node
	APIPort       int    // Port for the API of the node
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
	Type    string   `json:"type"`
	Data    []byte   `json:"data"`    // Data JSON depends on the type
	Visited []string `json:"visited"` // Which nodes this message has already visited
	TTL     int      `json:"ttl"`     // TTL for message
	ID      uint64   `json:"id"`      // De-duplication id
}

type GossipMessageTypeHello struct { // GossipMessage.Type == "hello"
	LocalCluster  string `json:"localCluster"`
	Port          int    `json:"port"`
	Name          string `json:"name"`
	Interface     string `json:"interface"`
	GlobalCluster string `json:"globalCluster"`
	APIPort       int    `json:"apiPort"`
}

type GossipMessageTypeHelloResponse struct {
	ClusterNodes []ClusterNode `json:"clusterNodes"`
}

type GossipMessageTypeAddIndex struct {
	DocID  uint64                 `json:"docID"`
	Fields map[string]interface{} `json:"fields"`
	App    string                 `json:"app"`
}

func InitMyNode() {
	AllNodes = make(map[string]*ClusterNode)
	MyClusterNode = &ClusterNode{
		LocalCluster:  *LocalClusterName,
		GlobalCluster: *GlobalClusterName,
		IP:            *NodeInterface,
		Port:          *NodePort,
		Name:          fmt.Sprintf("%s:%v", *NodeInterface, *NodePort),
		APIPort:       *APIPort,
	}
	AllNodes[MyClusterNode.Name] = MyClusterNode
	initMap := make(map[string]*ClusterNode)
	initMap[MyClusterNode.Name] = MyClusterNode
	MyLocalCluster = &LocalClusterGroup{
		Name:      *LocalClusterName,
		NodeCount: 1,
		Nodes:     initMap,
	}
	GMCache = &GossipMessageCache{}
}

// isRelay is for if this message is being relayed by other nodes, we don't want to keep relaying it
func addNodeToCluster(localCluster string, port int, name string, tcpAddr string, apiPort int) error {
	if _, ok := AllNodes[name]; ok { // TODO: Suspect node stuff
		log.Println("New node pretending to be old node, or I've seen this already, dropping add...")
		return nil
	} else {
		newNode := &ClusterNode{
			LocalCluster:  localCluster,
			IP:            tcpAddr,
			Port:          port,
			Name:          name,
			GlobalCluster: *GlobalClusterName,
			APIPort:       apiPort,
		}
		AllNodes[name] = newNode
		if localCluster == MyClusterNode.LocalCluster {
			MyLocalCluster.Nodes[name] = newNode
		}
		// Send gossip to more nodes
		return nil
	}
}

func handleGossipMessage(gospMsg GossipMessage, c net.Conn) {
	isDup := false
	log.Println("Message id", gospMsg.ID, "ttl", gospMsg.TTL)
	// Handle duplicate message
	if _, ok := (*GMCache)[gospMsg.ID]; ok { // Message is duplicate, gets put in cache at Broadcast
		log.Println("Message is duplicate!")
		isDup = true
	}
	// Handle message
	switch gospMsg.Type {
	case "hello":
		log.Println("New node just waved hello!")
		var gospMsgData GossipMessageTypeHello
		err := json.Unmarshal(gospMsg.Data, &gospMsgData)
		if err != nil {
			logger.Error(err)
			c.Write([]byte("JSON decoding error!\n"))
			c.Close()
			return
		}
		if gospMsgData.GlobalCluster != *GlobalClusterName {
			logger.Error("Node tried to join with another global cluster name!")
			c.Write([]byte("Different Global Cluster Name, Rejecting!\n"))
			c.Close()
			return
		}
		if !isDup { // de-duplication
			err = addNodeToCluster(gospMsgData.LocalCluster, gospMsgData.Port, gospMsgData.Name, gospMsgData.Interface, gospMsgData.APIPort)
			if err != nil {
				logger.Error(err)
				c.Write([]byte("JSON decoding error!\n"))
				c.Close()
				return
			}
		}
		// Respond with all the nodes I have
		sendAllNodes := GossipMessageTypeHelloResponse{}
		for k, v := range AllNodes {
			if k != gospMsgData.Name { // Don't send it itself
				sendAllNodes.ClusterNodes = append(sendAllNodes.ClusterNodes, *v)
			}
		}
		respData, err := json.Marshal(sendAllNodes)
		if err != nil {
			logger.Error(err)
			return
		}
		finalObj := GossipMessage{
			Data: respData,
			Type: "helloResponse",
		}
		finalData, err := json.Marshal(finalObj)
		if err != nil {
			logger.Error(err)
			return
		}
		c.Write([]byte(string(finalData) + "\n"))
		c.Close()
		data := GossipMessageTypeHello{
			LocalCluster: gospMsgData.LocalCluster,
			Name:         gospMsgData.Name,
			Port:         gospMsgData.Port,
			Interface:    gospMsgData.Interface,
		}
		msgData, err := json.Marshal(data)
		if err != nil {
			logger.Error(err)
			return
		}
		if gospMsg.TTL-1 > 0 {
			BroadcastGossipMessage(msgData, []string{gospMsgData.Name}, "addNode", gospMsg.TTL, gospMsg.ID)
		}
		// client will close connection

	case "addNode":
		log.Println("Gossip about new node!")
		var gospMsgData GossipMessageTypeHello
		err := json.Unmarshal(gospMsg.Data, &gospMsgData)
		if err != nil {
			logger.Error(err)
			c.Write([]byte("JSON decoding error!\n"))
			c.Close()
			return
		}
		if !isDup { // de-duplication
			err = addNodeToCluster(gospMsgData.LocalCluster, gospMsgData.Port, gospMsgData.Name, gospMsgData.Interface, gospMsgData.APIPort)
			if err != nil {
				logger.Error(err)
				c.Write([]byte("error adding new node!\n"))
				c.Close()
				return
			}
		}
		c.Write([]byte("Got it.\n"))
		data := GossipMessageTypeHello{
			LocalCluster: gospMsgData.LocalCluster,
			Name:         gospMsgData.Name,
			Port:         gospMsgData.Port,
			Interface:    gospMsgData.Interface,
		}
		msgData, err := json.Marshal(data)
		if err != nil {
			panic(err)
		}
		if gospMsg.TTL-1 > 0 {
			BroadcastGossipMessage(msgData, gospMsg.Visited, "addNode", gospMsg.TTL, gospMsg.ID)
		}
		// client will close connection

	case "addIndex":
		log.Println("\n\nGossip about new document!")
		var gospMsgData GossipMessageTypeAddIndex
		err := json.Unmarshal(gospMsg.Data, &gospMsgData)
		if err != nil {
			logger.Error(err)
			c.Write([]byte("JSON decoding addIndex error!\n"))
			c.Close()
			return
		}
		if !isDup {
			docID := addIndexFromGossip(gospMsgData.DocID, gospMsgData.App, gospMsgData.Fields) // TODO: Make dynamic
			log.Println("added doc id", docID, "from gossip!")
		}
		c.Write([]byte("Got it.\n"))
		if gospMsg.TTL-1 > 0 {
			BroadcastGossipMessage(gospMsg.Data, gospMsg.Visited, "addIndex", gospMsg.TTL, gospMsg.ID)
		}
		// client will close connection

	default:
		log.Println("Unrecognized message")
		c.Write([]byte("Unrecognized message\n"))
		c.Close()
	}

	return
}

func handleConnection(c net.Conn) {
	fmt.Printf("Serving %s\n", c.RemoteAddr().String())
	c.SetDeadline(time.Now().Add(time.Second * 5))
	defer func() {
		fmt.Printf("Closing connection to %s\n", c.RemoteAddr().String())
		c.Close()
	}()

	// request, err := bufio.NewReader(c).ReadString('\n')
	// result := "heyyy\n"
	// c.Write([]byte(string(result)))
	// }

	var gospMsg GossipMessage
	decoder := json.NewDecoder(c)
	err := decoder.Decode(&gospMsg)
	if err != nil {
		if strings.Contains(err.Error(), "i/o timeout") {
			logger.Warn("Timeout for connection", c.RemoteAddr().String())
		} else {
			logger.Error("Error handling tcp connection:")
			logger.Error(err)
		}
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
		if strings.Contains(err.Error(), "i/o timeout") {
			// Already logged somewhere else
		} else {
			log.Printf("error: %v\n", err)
		}
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

func AddNodeGossipMessage(nodeAddr string) {
	con, err := net.Dial("tcp4", nodeAddr)
	if err != nil {
		log.Println("Could not connect to fellow node!")
		log.Fatalln(err)
	}

	con.SetDeadline(time.Now().Add(time.Second * 5))

	serverReader := bufio.NewReader(con)

	switch err {
	case nil:
		data := GossipMessageTypeHello{
			LocalCluster:  MyClusterNode.LocalCluster,
			Name:          MyClusterNode.Name,
			Port:          MyClusterNode.Port,
			Interface:     MyClusterNode.IP,
			APIPort:       MyClusterNode.APIPort,
			GlobalCluster: *GlobalClusterName,
		}
		msgData, err := json.Marshal(data)
		if err != nil {
			panic(err)
		}
		rand.Seed(time.Now().UnixNano())
		newID := rand.Uint64()
		msg := GossipMessage{
			Data: msgData,
			Type: "hello",
			ID:   newID,
			TTL:  6,
		}
		jsonData, err := json.Marshal(msg)
		if err != nil {
			panic(err)
		}
		if _, err = con.Write([]byte(string(jsonData) + "\n")); err != nil {
			log.Printf("failed to send the client request: %v\n", err)
		}
		if err != nil {
			log.Println("erro:", err)
			logger.Error(err)
		}
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
	serverResponse, err := serverReader.ReadString('\n') // TODO: Read as bytes instead maybe?

	switch err {
	case nil:
		// Add all of the existing nodes
		if strings.Contains(serverResponse, "Different Global Cluster Name, Rejecting!") {
			log.Println("Got rejected! Tried to join a different global cluster!")
			con.Close()
			return
		}
		var helloResp GossipMessage
		err = json.Unmarshal([]byte(serverResponse), &helloResp)
		if err != nil {
			logger.Error("Error unmarshalling hello response:")
			logger.Error(err)
			con.Close()
			return
		}
		var helloRespData GossipMessageTypeHelloResponse
		err = json.Unmarshal(helloResp.Data, &helloRespData)
		if err != nil {
			logger.Error("Error unmarshalling hello response data:")
			logger.Error(err)
			con.Close()
			return
		}
		for _, node := range helloRespData.ClusterNodes {
			addNodeToCluster(node.LocalCluster, node.Port, node.Name, node.IP, node.APIPort)
		}
	case io.EOF:
		log.Println("server closed the connectionn")
		return
	default:
		if strings.Contains(err.Error(), "i/o timeout") {
			log.Println("IO Timeout...")
		} else {
			log.Printf("server error: %v\n", err)
		}
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
		for _, node := range nodeList {
			AddNodeGossipMessage(node)
		}
	} else {
		log.Println("No other nodes")
	}
}

func LeaveGossipCluster() {
	// TODO: Broadcast leave
	GossipServer.Close()
}

func SendGossipMessage(msg []byte, addr string) {
	log.Println("Sending to", addr)
	con, err := net.Dial("tcp4", addr)
	if err != nil {
		log.Println("Could not connect to fellow node!")
		log.Fatalln(err)
	}

	serverReader := bufio.NewReader(con)

	switch err {
	case nil:
		if _, err = con.Write([]byte(string(msg) + "\n")); err != nil {
			log.Printf("failed to send the client request: %v\n", err)
		}
		if err != nil {
			logger.Error(err)
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
		log.Println("Got from serverrr:", strings.TrimSpace(serverResponse))
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

func ClusterAddIndex(newDocStruct *GossipMessageTypeAddIndex) error {
	addIndexFromGossip(newDocStruct.DocID, newDocStruct.App, newDocStruct.Fields)
	return nil
}

func TellClusterAddIndex(newDocStructData []byte) {
	fmt.Println("Telling cluster to add index")
	for _, node := range AllNodes {
		if node.Name == MyClusterNode.Name { // Don't do me
			continue
		}
		resp, err := http.Post(fmt.Sprintf("http://%s:%d/cluster/addIndex", node.IP, node.APIPort), "application/json", bytes.NewBuffer(newDocStructData))
		if err != nil {
			log.Println("Error making request to another node:")
			log.Println(err)
			logger.Error(err)
			continue
		}
		if resp.StatusCode > 299 {
			log.Println(">300 status code:")
			log.Println(resp.StatusCode)
			body, _ := ioutil.ReadAll(resp.Body)
			log.Println(string(body))
		}
		resp.Body.Close()
	}
}

// void is just an empty struct
type void struct{}

// BroadcastGossipMessage picks 3 random nodes in the cluster
// msg is the struct of a JSON encoded message to send
// sourceName is the source of the message, which will be used to avoid sending a duplicate message to itself. Functions adds itself to list
// msgType is the type of the message e.g. "addNode"
// ttl is TTL of message. Function decrements ttl
// id is de-duplication id of message
func BroadcastGossipMessage(data []byte, sourceNames []string, msgType string, ttl int, id uint64) {
	(*GMCache)[id] = void{}
	// Pick random nodes on in visited list
	msg := GossipMessage{
		Data:    data,
		Type:    msgType,
		Visited: append(sourceNames, MyClusterNode.Name),
		TTL:     ttl - 1,
		ID:      id,
	}
	jsonData, err := json.Marshal(msg)
	if err != nil {
		logger.Error("error marshalling json data to broadcast gossip message:")
		logger.Error(err)
		return
	}
	// Send to up to 3 random nodes
	if len(AllNodes) > 3 { // TODO: Make a faster way to pick a random node from a map (this is O(3n)), maybe also have a slice?
		for i := 0; i < 3; i++ {
			// Get random index
			randIndex := rand.Intn(len(AllNodes))
			// Send to that node
			indx := 0
			for _, node := range AllNodes {
				if randIndex == indx {
					SendGossipMessage(jsonData, node.Name)
				}
				indx++
			}
		}
	} else {
		for _, node := range AllNodes {
			if node.Name != MyClusterNode.Name {
				SendGossipMessage(jsonData, node.Name)
			}
		}
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
