package main

import (
	"fmt"
	"log"
	"os"
	"github.com/yuin/goldmark"
	"path/filepath"
	"strings"
)

func main() {

	dirPath := "internal/content"
	files, err := os.ReadDir(dirPath)
	if err != nil {
		log.Fatal(err)
	}
	for _, entry := range files {
		if entry.IsDir() {
			continue
		}
		if filepath.Ext(entry.Name()) != ".md" {
				continue
		}
		fullPath := filepath.Join(dirPath, entry.Name())
		content, err := os.ReadFile(fullPath)
		if err != nil {
			fmt.Printf("Error reading file %s: %v\n", entry.Name(), err)
			continue
		}
		name := strings.TrimSuffix(entry.Name(), ".md")
		outputPath := filepath.Join("public", name+".html")

		file, err := os.Create(outputPath)
		if err != nil {
			log.Printf("Error creating %s: %v\n", outputPath, err)
			continue
		}
		err = goldmark.Convert(content, file)
		if err != nil {
			log.Printf("Error converting %s: %v\n", entry.Name(), err)
		}

		file.Close()
	}	
}



