package main

import ()

func main() {
	app := initApp("Example App")
	// cluster := initCluster()
	CheckDocumentsFolder()
	LoadIndexesFromDisk(app)
	StartWebserver(app)
}
