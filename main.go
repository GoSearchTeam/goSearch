package main

import (
	"io"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	syscall.Umask(0) // file mode perms
	logFile, _ := os.OpenFile("logs.log", os.O_CREATE|os.O_APPEND|os.O_RDWR, 0755)
	mw := io.MultiWriter(os.Stdout, logFile)
	log.SetOutput(mw)
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
	var apps []*appIndexes
	app := initApp("Example App") // For testing
	apps = append(apps, app)      // For testing
	go func() {
		// cluster := initCluster()
		CheckDocumentsFolder()
		// apps = append(apps, LoadAppsFromDisk()...)
		// apps = LoadAppsFromDisk()
		log.Println("Finished Loading apps")
		log.Println(apps)
		for _, app := range apps {
			// app.LoadIndexesFromDisk()
			log.Println("starting webservers...")
			StartWebserver(app) // FIXME: HOLY SHIT THIS IS DUMB, need to update to pass app name in url
		}
	}()
	<-c
	log.Println("### Serializing apps and indexes before exiting...")
	for _, app := range apps {
		// app.SerializeIndex()
		app.SerializeApp()
	}
	os.Exit(1)
}
