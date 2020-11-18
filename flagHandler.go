package main

import (
	"flag"
)

var (
	// FellowNodes String of other nodes in global cluster
	FellowNodes *string
	// ClusterMode Whether to activate cluster mode features
	ClusterMode *bool
	// LocalClusterName Name of the local cluster
	LocalClusterName *string
	// NodePort Port of the node
	NodePort *int
	// NodeInterface Interface of the node
	NodeInterface *string
	APIPort       *int
)

func ParseFlags() {
	ClusterMode = flag.Bool("cluster-mode", false, "Bool whether to activate clustering for this node")
	FellowNodes = flag.String("fellow-nodes", "", "A CSV of fellow nodes to begin discovery if not first node in cluster. Format: \"x.x.x.x:port,y.y.y.y:port\"")
	LocalClusterName = flag.String("local-cluster", "", "Name of the local cluster. Typically this is the name of the Availability Zone")
	NodePort = flag.Int("noise-port", -1, "Port of the node to join the cluster with")
	APIPort = flag.Int("port", -1, "Port for the API to listen to")
	NodeInterface = flag.String("iface", "", "Interface IP to use when joining the cluster")
	flag.Parse()
}
