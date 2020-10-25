package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
	app := initApp("Example App")
	go func() {
		// cluster := initCluster()
		CheckDocumentsFolder()
		// LoadIndexesFromDisk(app)
		StartWebserver(app)
	}()
	<-c
	fmt.Println("### Serializing index before exiting...")
	app.SerializeIndex()
	os.Exit(1)
}
