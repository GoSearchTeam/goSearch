package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/gin-gonic/gin"
)

type ClusterNodesResponseNode struct {
	LocalCluster string `json:"localCluster"`
	IP           string `json:"ip"`
	Port         int    `json:"port"`
	Status       string `json:"status"`
	Name         string `json:"name"`
}

type ClusterNodesResponse struct {
	Nodes []ClusterNodesResponseNode `json:"nodes"`
}

func HandleClusterRoutes(r *gin.Engine) {
	clusterGroup := r.Group("/cluster")

	clusterGroup.GET("/nodes", func(c *gin.Context) {
		nodeList := make([]*ClusterNode, 0)
		for _, v := range AllNodes {
			nodeList = append(nodeList, v)
		}
		c.JSON(200, nodeList)
	})

	clusterGroup.POST("/addIndex", func(c *gin.Context) {
		fmt.Println("Being told to add index")
		data, _ := ioutil.ReadAll(c.Request.Body)
		jDat := GossipMessageTypeAddIndex{}
		json.Unmarshal([]byte(data), &jDat)
		ClusterAddIndex(&jDat)
		c.String(200, "Done")
	})
}
