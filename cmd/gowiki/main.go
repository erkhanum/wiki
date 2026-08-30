package main

import (
	"fmt"
	"log"
	"os"
	"github.com/yuin/goldmark"
	"path/filepath"
	"strings"
	"io/fs"
	"gopkg.in/yaml.v3"
	"html/template"
	"bytes"
)

type FrontMatter struct {
	Title string  `yaml:"title"`
	Author string `yaml:"author"`
	Date string   `yaml:"date"`
}


type Page struct {
	Title   string
	Author  string
	Date    string
	Content template.HTML
}


func parseFrontMatter(data []byte) (FrontMatter, []byte, error) {
	var fm FrontMatter
	text := string(data)
	if !strings.HasPrefix(strings.TrimSpace(text), "---") {
		return fm, data, nil
	}
	firstEnd := strings.Index(text, "---")
	if firstEnd == -1 {
		return fm, data, nil
	}
	remaining := text[firstEnd+3:]

	secondEnd := strings.Index(remaining, "---")
	if secondEnd == -1 {
		return fm, data, fmt.Errorf("closing front matter delimiter not found")
	}

	yamlData := remaining[:secondEnd]

	err := yaml.Unmarshal([]byte(yamlData), &fm)
	if err != nil {
		return fm, data, fmt.Errorf("invalid front matter: %w", err)
	}
	markdown := remaining[secondEnd+3:]

	return fm, []byte(markdown), nil
}

func main() {
	
	dirPath := "internal/content"
	pageTemplate, err := template.ParseFiles("templates/page.html")
	if err != nil {
		log.Fatal(err)
	}

	err = filepath.WalkDir(dirPath, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
				return nil
		} 
		if filepath.Ext(d.Name()) != ".md" {
			return nil
		}
		
		data, err := os.ReadFile(path)
		if err != nil {
			fmt.Printf("Error reading %s: %v\n", path, err)
			return nil
		}

		frontMatter, markdown, err := parseFrontMatter(data)
		if err != nil {
			fmt.Printf("Error parsing %s: %v\n", path, err)
			return nil
		}
		var buffer bytes.Buffer

		err = goldmark.Convert(markdown, &buffer)
		if err != nil {
			fmt.Printf("Error converting %s: %v\n", path, err)
			return nil
		}

		// Create Page data
		page := Page{
			Title:   frontMatter.Title,
			Author:  frontMatter.Author,
			Date:    frontMatter.Date,
			Content: template.HTML(buffer.String()),
		}

		relativePath, err := filepath.Rel(dirPath, path)
		if err != nil {
			return err
		}

		relativePath = strings.TrimSuffix(relativePath, ".md") + ".html"

		outputPath := filepath.Join("public", relativePath)

		err = os.MkdirAll(filepath.Dir(outputPath), 0755)
		if err != nil {
			fmt.Printf("Error creating directory for %s: %v\n", outputPath, err)
			return nil
		}
		file, err := os.Create(outputPath)
		if err != nil {
			fmt.Printf("Error creating file %s: %v\n", outputPath, err)
			return nil
		}
		err = pageTemplate.Execute(file, page)
		if err != nil {
			fmt.Printf("Error executing template for %s: %v\n", path, err)
		}

		file.Close()

		fmt.Printf("Generated: %s\n", outputPath)

		return nil
	})
	if err != nil {
			log.Fatal(err)
		}
	
}