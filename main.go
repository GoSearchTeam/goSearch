package main

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"io"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var (
	logger *logrus.Logger
)

func main() {
	syscall.Umask(0) // file mode perms

	// Logger
	logFile, _ := os.OpenFile("debug.log", os.O_CREATE|os.O_APPEND|os.O_RDWR, 0666)
	eventFile, err := os.OpenFile("events.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal(err)
	}
	logger = logrus.New()
	logger.SetFormatter(&logrus.JSONFormatter{})
	logger.SetOutput(eventFile)
	mw := io.MultiWriter(os.Stdout, logFile)
	log.SetOutput(mw)
	// Start log file rotation
	go monitorFileSize()

	// Command line arguments and flags
	ParseFlags()

	// Activate Clustering
	if *ClusterMode {
		BeginClusterDiscovery()
	}

	// Apps
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
	// next 2 lines temp for noise testing
	// <-c
	// os.Exit(1)
	initApp("Example App") // For testing
	go func() {
		// cluster := initCluster()
		CheckDocumentsFolder()
		LoadAppsFromDisk()
		for _, app := range Apps {
			start := time.Now()
			app.LoadIndexesFromDisk()
			end := time.Now()
			fmt.Println("Load index time", end.Sub(start))
			log.Println("starting webservers...")
			StartWebserver(app) // FIXME: HOLY SHIT THIS IS DUMB, need to update to pass app name in url
		}
	}()
	<-c
	log.Println("### Serializing apps and indexes before exiting...")
	for _, app := range Apps {
		start := time.Now()
		app.SerializeIndex()
		end := time.Now()
		fmt.Println("Serialize index time", end.Sub(start))
		app.SerializeApp()
	}
	os.Exit(1)
}
