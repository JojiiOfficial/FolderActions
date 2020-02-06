package main

import (
	"fmt"
	"log"

	"github.com/fsnotify/fsnotify"
)

func startWatcher(dir, currPath string) {
	if verbose && !quiet {
		fmt.Println("Started watcher for", dir)
	}

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer watcher.Close()

	done := make(chan bool)
	go func() {
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}

				if event.Op&fsnotify.Create == fsnotify.Create {
					if verbose && !quiet {
						fmt.Println("create", event.Name)
					}
				} else if event.Op&fsnotify.Remove == fsnotify.Remove {
					if verbose && !quiet {
						fmt.Println("remove:", event.Name)
					}
				}
			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				log.Println("error:", err)
			}
		}
	}()

	err = watcher.Add(dir)
	if err != nil {
		log.Fatal(err)
	}
	<-done
}
