package main

import (
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/russross/blackfriday/v2"
)

// Page represents a markdown page
type Page struct {
	Title string
	Body  template.HTML
}

// Navbar represents the navigation bar
type Navbar struct {
	Home       string
	GitHub     string
	Email      string
	DarkMode   bool
	DarkModeJS template.JS
}

// Footer represents the footer
type Footer struct {
	Text string
}

// GeneratePageHTML generates the HTML for a markdown page
func GeneratePageHTML(markdownPath string) (template.HTML, error) {
	data, err := ioutil.ReadFile(markdownPath)
	if err != nil {
		return "", err
	}

	renderer := blackfriday.NewHTMLRenderer(blackfriday.HTMLRendererParameters{
		Flags: blackfriday.CommonHTMLFlags | blackfriday.HrefTargetBlank,
	})

	html := blackfriday.Run(data, blackfriday.WithRenderer(renderer))
	return template.HTML(html), nil
}

// GenerateIndexHTML generates the HTML for the index page
func GenerateIndexHTML(pages []string) (string, error) {
	indexTmpl, err := template.ParseFiles("templates/index.html")
	if err != nil {
		return "", err
	}

	builder := strings.Builder{}
	err = indexTmpl.Execute(&builder, struct {
		Navbar Navbar
		Pages  []string
		Footer Footer
	}{
		Navbar: Navbar{
			Home:       "index.html",
			GitHub:     "https://github.com/",
			Email:      "mailto:example@example.com",
			DarkMode:   true,
			DarkModeJS: template.JS(`toggleDarkMode()`),
		},
		Pages:  pages,
		Footer: Footer{Text: "© 2023 My Website. All rights reserved."},
	})

	if err != nil {
		return "", err
	}

	return builder.String(), nil
}

func main() {
	markdownDir := "markdown"
	outputDir := "public"

	err := os.MkdirAll(outputDir, os.ModePerm)
	if err != nil {
		log.Fatal(err)
	}

	files, err := ioutil.ReadDir(markdownDir)
	if err != nil {
		log.Fatal(err)
	}

	pages := []string{}
	for _, file := range files {
		if file.IsDir() {
			continue
		}

		markdownPath := filepath.Join(markdownDir, file.Name())
		pagePath := filepath.Join(outputDir, strings.TrimSuffix(file.Name(), ".md")+".html")

		html, err := GeneratePageHTML(markdownPath)
		if err != nil {
			log.Fatal(err)
		}

		pageTmpl, err := template.ParseFiles("templates/page.html")
		if err != nil {
			log.Fatal(err)
		}

		fileContent := strings.Builder{}
		err = pageTmpl.Execute(&fileContent, Page{
			Title: strings.TrimSuffix(file.Name(), ".md"),
			Body:  html,
		})

		if err != nil {
			log.Fatal(err)
		}

		err = ioutil.WriteFile(pagePath, []byte(fileContent.String()), 0644)
		if err != nil {
			log.Fatal(err)
		}

		pages = append(pages, strings.TrimPrefix(pagePath, outputDir+"/"))
	}

	indexHTML, err := GenerateIndexHTML(pages)
	if err != nil {
		log.Fatal(err)
	}

	indexPath := filepath.Join(outputDir, "index.html")
	err = ioutil.WriteFile(indexPath, []byte(indexHTML), 0644)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Static webpage built successfully in the '%s' directory.\n", outputDir)
}
