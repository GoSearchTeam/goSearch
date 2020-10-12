package main

import ()

func main() {
	app := initApp("Example App")
	CheckDocumentsFolder()
	LoadIndexesFromDisk(app)
	StartWebserver(app)
}
