package main

import (
	"fmt"
	"log"
	"os"
	"github.com/yuin/goldmark"
	"path/filepath"
	"strings"
	"io/fs"
)

func main() {

	dirPath := "internal/content"

	err := filepath.WalkDir(dirPath, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
				return nil
		} 
		if filepath.Ext(d.Name()) != ".md" {
			return nil
		}
		
		content, err := os.ReadFile(path)
		if err != nil {
			fmt.Printf("Error reading file %s: %v\n", d.Name(), err)
			return nil
		}
		relativePath, err := filepath.Rel(dirPath, path)
		if err != nil {
			return err
		}

		relativePath = strings.TrimSuffix(relativePath, ".md") + ".html"

		outputPath := filepath.Join("public", relativePath)

		file, err := os.Create(outputPath)
		if err != nil {
			fmt.Printf("Error creating file %s: %v\n", outputPath, err)
			return nil
		}

		err = goldmark.Convert(content, file)
		if err != nil {
			fmt.Printf("Error converting %s: %v\n", path, err)
		}

		file.Close()

		fmt.Printf("Generated: %s\n", outputPath)

		return nil
	})
	if err != nil {
			log.Fatal(err)
		}
}