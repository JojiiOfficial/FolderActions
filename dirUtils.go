package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
)

func createDirIfNotExists(dir string) error {
	stat, err := os.Stat(dir)
	if err == nil {
		if !stat.IsDir() {
			return errors.New("Name already taken by a file")
		}
		return nil
	}

	return os.Mkdir(dir, 0700)
}

func listCurrFolders() []string {
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
