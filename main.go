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
	logFile, _ := os.OpenFile("logs.log", os.O_CREATE|os.O_APPEND|os.O_RDWR, 0777)
	mw := io.MultiWriter(os.Stdout, logFile)
	log.SetOutput(mw)
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
	app := initApp("Example App")
	go func() {
		// cluster := initCluster()
		CheckDocumentsFolder()
		LoadIndexesFromDisk(app)
		StartWebserver(app)
	}()
	<-c
	log.Println("### Serializing index before exiting...")
	app.SerializeIndex()
	os.Exit(1)
}
