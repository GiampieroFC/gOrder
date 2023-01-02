package main

import (
	"embed"
	"encoding/json"
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
)

type Directory struct {
	Name    string   `json:"name"`
	Formats []string `json:"formats"`
}

var directories []Directory

/*
A //go:embed directive above a variable declaration specifies which files to embed:
*/
//go:embed formats.json
var embedJson embed.FS

func errorHandler(e error, msg string) {
	if e != nil {
		fmt.Printf("\a \n❌ %s: %v\n", msg, e)
		log.Fatalln(e)
		panic(e)
	}
}

func loadDir() {
	fmt.Print("💫 loading formats")
	// if you don't want embed the json. And just read it when execute the program, you can change this line for:
	// b, err := os.ReadFile("formats.json")
	b, err := embedJson.ReadFile("formats.json")
	errorHandler(err, "can't read .json")
	e := json.Unmarshal(b, &directories)
	errorHandler(e, "can't convert .json")
	fmt.Println(".json read successfully 📖")
}

func guardar(obj *Directory, file fs.DirEntry) {

	for _, extension := range obj.Formats {
		if !file.IsDir() {
			if extension == filepath.Ext(file.Name()) {
				here, _ := os.Getwd()
				err := os.MkdirAll(obj.Name, os.ModePerm)
				errorHandler(err, "can't create directory")
				oldPath := filepath.Join(here, file.Name())
				var newPath string = filepath.Join(here, obj.Name, file.Name())
				err = os.Rename(oldPath, newPath)
				errorHandler(err, "can't move directory")
				fmt.Printf("👍 '%s' saved correctly in [./%s]\n", file.Name(), newPath)
			}
		}
	}
}

func cleanEmtydirs(file fs.DirEntry) {
	if file.IsDir() {
		info, err := file.Info()
		errorHandler(err, "I don't know what is this")
		files, err := os.ReadDir(info.Name())
		errorHandler(err, "I can't read it. Is this a directory?")
		if len(files) == 0 {
			os.Remove(info.Name())
			fmt.Printf("💢 The empty directory %s was deleted\n", info.Name())
		}
	}
}

func main() {
	files, e := os.ReadDir(".")
	for _, v := range files {
		cleanEmtydirs(v)
	}
	errorHandler(e, "can't read current directory")
	loadDir()
	for i := 0; i < len(directories); i++ {
		for _, file := range files {
			guardar(&directories[i], file)
		}
	}
	fmt.Println("We have finished 🙃")
}
