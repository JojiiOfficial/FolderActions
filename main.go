package main

import (
	"fmt"
	"os"
	"path"
	"time"

	"gopkg.in/alecthomas/kingpin.v2"
)

var (
	//OS args
	folders        []string
	verbose, quiet bool

	//path data
	scriptPath string
)

const (
	envVarScriptPath = "FA_SCRIPT_PATH"
)

func init() {
	a := kingpin.Flag("dir", "The list of directories to watch").HintAction(listCurrFolders).Short('d').Strings()
	v := kingpin.Flag("verbose", "Increase debug message").Short('v').Bool()
	q := kingpin.Flag("quiet", "Don't output messages").Short('q').Bool()
	kingpin.Parse()

	for _, e := range *a {
		folders = append(folders, e)
	}
	verbose = *v
	quiet = *q
}

func getCurrPath() string {
	exec, err := os.Executable()
	if err != nil {
		fmt.Println("Error getting current exec path:", err.Error())
		os.Exit(1)
		return ""
	}
	currPath, _ := path.Split(exec)
	return currPath
}

func main() {
	envvar := os.Getenv(envVarScriptPath)
	if len(envvar) > 0 {
		scriptPath = envvar
	} else {
		scriptPath = getCurrPath() + "scripts/"
		createDirIfNotExists(scriptPath)
	}
	if len(folders) > 0 {
		for _, dir := range folders {
			go startWatcher(dir, scriptPath)
		}

		for {
			time.Sleep(1 * time.Hour)
		}
	} else if !quiet {
		fmt.Println("No dir found! use --dir to specify directories")
	}
}
