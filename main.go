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
	logFile, _ := os.OpenFile("process.log", os.O_CREATE|os.O_APPEND|os.O_RDWR, 0666)
	eventFile, err := os.OpenFile("events.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal(err)
	}
	logger = logrus.New()
	logger.SetFormatter(&logrus.JSONFormatter{})
	logger.SetOutput(eventFile)
	mw := io.MultiWriter(os.Stdout, logFile)
	log.SetOutput(mw)
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
	initApp("Example App") // For testing
	go func() {
		// cluster := initCluster()
		CheckDocumentsFolder()
		LoadAppsFromDisk()
		log.Println("Finished Loading apps")
		log.Println(Apps)
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
