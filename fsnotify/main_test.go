package main

import (
	"github.com/fsnotify/fsnotify"
	"log"
	"testing"
)

func TestName(t *testing.T) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal("NewWatcher failed: ", err)
	}
	defer watcher.Close()

	done := make(chan bool)
	go func() {
		defer close(done)
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}
				log.Printf("%s %s\n", event.Name, event.Op)
			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				log.Println("error:", err)
			}
		}
	}()

	err = watcher.Add("E:\\vuln\\log\\log1")
	err = watcher.Add("E:\\vuln\\log\\log2")
	if err != nil {
		log.Fatal("Add failed:", err)
	}
	<-done
}
