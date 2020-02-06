package main

import (
	"fmt"
	"log"
	"os/exec"
	"strings"

	"github.com/fsnotify/fsnotify"
)

func startWatcher(dir string) {
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
				handleEvent(dir, event)
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

func eventToScriptPath(dir, event string) string {
	if len(dir) == 0 {
		return ""
	}
	nameDir := strings.ReplaceAll(dir, "/", "_")
	if strings.HasSuffix(nameDir, "_") {
		nameDir = nameDir[:len(nameDir)-1]
	}
	return nameDir + "_" + event + ".sh"
}

func handleEvent(dir string, event fsnotify.Event) {
	name := event.Name
	var scriptFile string
	if event.Op&fsnotify.Create == fsnotify.Create {
		if verbose && !quiet {
			fmt.Println("create", name)
		}
		scriptFile = scriptPath + eventToScriptPath(dir, "create")
	} else if event.Op&fsnotify.Remove == fsnotify.Remove {
		if verbose && !quiet {
			fmt.Println("remove:", dir, name)
		}
		scriptFile = scriptPath + eventToScriptPath(dir, "delete")
	} else {
		return
	}
	if !quiet && verbose {
		fmt.Println("to run scriptfile:", scriptFile)
	}
	err := runScript(scriptFile, name)
	if err != nil && verbose && !quiet {
		fmt.Println(err.Error())
	}
}

func runScript(scriptFile, name string) error {
	return exec.Command(scriptFile, name).Start()
}
