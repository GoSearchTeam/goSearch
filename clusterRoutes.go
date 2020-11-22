package main

import (
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
}
