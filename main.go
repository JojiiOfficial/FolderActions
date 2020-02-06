package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
	"time"

	"github.com/fsnotify/fsnotify"
	"gopkg.in/alecthomas/kingpin.v2"
)

var folders []string
var verbose, quiet bool

func init() {
	a := kingpin.Flag("dir", "The list of directories to watch").HintAction(listCurrItems).Short('d').Strings()
	v := kingpin.Flag("verbose", "Increase debug message").Short('v').Bool()
	q := kingpin.Flag("quiet", "Don't output messages").Short('q').Bool()
	kingpin.Parse()

	for _, e := range *a {
		folders = append(folders, e)
	}
	verbose = *v
	quiet = *q
}

func listCurrItems() []string {
	files, err := ioutil.ReadDir("./")
	if err != nil {
		fmt.Println("Error listing curr dir:", err.Error())
		return []string{}
	}
	var s []string

	for _, e := range files {
		if e.IsDir() {
			s = append(s, e.Name())
		}
	}
	return s
}

func main() {
	exec, err := os.Executable()
	if err != nil {
		fmt.Println("Error getting current exec path:", err.Error())
		os.Exit(1)
		return
	}
	currPath, _ := path.Split(exec)

	if len(folders) > 0 {
		for _, dir := range folders {
			go startWatcher(dir, currPath)
		}

		for {
			time.Sleep(1 * time.Hour)
		}
	} else if !quiet {
		fmt.Println("No dir found! use --dir to specify directories")
	}
}

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
