package main

import (
	"fmt"
	"os"
	"path"
	"strings"
	"time"

	"gopkg.in/alecthomas/kingpin.v2"
)

var (
	//OS args
	folders     []string
	verbose     bool
	quiet       bool
	allowUnsafe bool

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
	aus := kingpin.Flag("allow-unsafe", "Allow special characters in path and filename").Short('a').Bool()
	sp := kingpin.Flag("scripdir", "The folder containing the scripts").Short('s').HintAction(listCurrFolders).String()
	kingpin.Parse()

	for _, e := range *a {
		folders = append(folders, e)
	}
	verbose = *v
	quiet = *q
	allowUnsafe = *aus
	if len(*sp) > 0 {
		scriptPath = *sp
	}
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
	if len(scriptPath) == 0 {
		envvar := os.Getenv(envVarScriptPath)
		if len(envvar) > 0 {
			scriptPath = envvar
		} else {
			scriptPath = getCurrPath() + "scripts/"
		}
	} else {
		s, err := os.Stat(scriptPath)
		if err == nil && !s.IsDir() {
			fmt.Println("Error! Given ScriptDir is no dir!")
			os.Exit(1)
			return
		}
	}
	if err := createDirIfNotExists(scriptPath); err != nil {
		fmt.Println("Error:", err.Error())
		return
	}

	if !strings.HasSuffix(scriptPath, "/") {
		scriptPath += "/"
	}

	if !strings.HasPrefix(scriptPath, "/") {
		sp := scriptPath
		if strings.HasPrefix(sp, "./") {
			sp = sp[2:]
		}
		scriptPath = getCurrPath() + sp
	}

	if verbose && !quiet {
		fmt.Println("ScriptPath:", scriptPath)
	}

	if len(folders) > 0 {
		for _, dir := range folders {
			go startWatcher(dir)
		}

		for {
			time.Sleep(1 * time.Hour)
		}
	} else if !quiet {
		fmt.Println("No dir found! use --dir to specify directories")
	}
}
