package main

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"time"
)

func monitorFileSize() {
	for true {
		time.Sleep(time.Second * 60) // Every minute
		checkFileSize()
	}
}

func checkFileSize() {
	stat, err := os.Stat("events.log")
	if os.IsExist(err) {
		// Ignore
	} else if err != nil {
		panic(err)
	}
	if stat.Size() > 10000000 { // TODO: Make customizable in config, default 10MB
		// rotate the file
		newFileName := fmt.Sprintf("events-%s.log", time.Now().Format("2006-01-02_15-04-05"))
		os.Rename("events.log", newFileName)
		os.Create("events.log")
		newFile, err := os.Create(fmt.Sprintf("events-%s.zip", time.Now().Format("2006-01-02_15-04-05")))
		if err != nil {
			panic(err)
		}
		w := zip.NewWriter(newFile)
		defer w.Close()
		// f, err := w.Create(newFileName)
		if err != nil {
			panic(err)
		}
		openFile, err := os.Open(newFileName)
		defer openFile.Close()
		if err != nil {
			panic(err)
		}
		info, err := openFile.Stat()
		if err != nil {
			panic(err)
		}
		header, err := zip.FileInfoHeader(info)
		header.Method = zip.Deflate
		writer, err := w.CreateHeader(header)
		if err != nil {
			panic(err)
		}
		_, err = io.Copy(writer, openFile)
		if err != nil {
			panic(err)
		}
	}
}
