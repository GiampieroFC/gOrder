package main

import (
	"embed"
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"
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
	fmt.Print("loading formats")
	// if you don't want embed the json. And just read it when execute the program, you can change this line for:
	// b, err := os.ReadFile("formats.json")
	b, err := embedJson.ReadFile("formats.json")
	errorHandler(err, "can't read .json")
	e := json.Unmarshal(b, &directories)
	errorHandler(e, "can't convert .json")
	fmt.Println(".json read successfully")
}

func exten(file fs.DirEntry) (string, error) {
	if !file.IsDir() {
		var tipo string
		var ext []string = strings.Split(file.Name(), ".")
		tipo = "." + ext[len(ext)-1]
		tipo = strings.ToLower(tipo)
		return tipo, nil
	}
	if file.IsDir() {
		return "", nil
	}
	return file.Name(), errors.New("\n❗It isn't a directory and isn't a file 🤔")
}

func guardar(obj *Directory, file fs.DirEntry) {
	ext, fail := exten(file)
	errorHandler(fail, "can't determinate extension")
	for _, extension := range obj.Formats {
		if extension == ext {
			here, _ := os.Getwd()
			err := os.MkdirAll(obj.Name, os.ModePerm)
			errorHandler(err, "can't create directory")
			oldPath := filepath.Join(here, file.Name())
			var newPath string = filepath.Join(here, obj.Name, file.Name())
			err = os.Rename(oldPath, newPath)
			errorHandler(err, "can't move directory")
		}
	}
}

func readingDir(de []fs.DirEntry) {
	for _, d := range de {
		if d.IsDir() {
			info, _ := d.Info()
			files, _ := os.ReadDir(info.Name())
			if len(files) == 0 {
				os.Remove(info.Name())
			}
		}
	}
}

func main() {
	files, e := os.ReadDir(".")
	errorHandler(e, "can't read current directory")
	loadDir()
	for i := 0; i < len(directories); i++ {
		for _, file := range files {
			guardar(&directories[i], file)
			fmt.Printf("'%s' saved correctly in [./%s]", file.Name(), directories[i].Name)
		}
	}
	readingDir(files)
	fmt.Println("We have finished 🙃")
}
