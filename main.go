package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/shurcooL/github_flavored_markdown"
)

// Page represents a single page with its title and HTML content.
type Page struct {
	Title   string
	Content string
}

// generatePages traverses the specified directory, reads markdown files,
// converts them to HTML, and generates Page objects for each file.
func generatePages(dirPath string) ([]Page, error) {
	var pages []Page

	err := filepath.Walk(dirPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() && filepath.Ext(path) == ".md" && filepath.Base(path) != "README.md" {
			content, err := ioutil.ReadFile(path)
			if err != nil {
				return err
			}

			htmlContent := string(github_flavored_markdown.Markdown(content))

			var page Page
			page.Title = strings.TrimSuffix(info.Name(), ".md")
			page.Content = htmlContent

			pages = append(pages, page)
		}

		return nil
	})

	return pages, err
}

// getReadmeContent reads the contents of the README.md file, converts it to HTML,
// and returns it as a string.
func getReadmeContent(readmePath string) (string, error) {
	content, err := ioutil.ReadFile(readmePath)
	if err != nil {
		return "", err
	}

	htmlContent := string(github_flavored_markdown.Markdown(content))
	return htmlContent, nil
}

func main() {
	// Directory containing the markdown files
	dirPath := "markdown"

	// Path to the README.md file
	readmePath := filepath.Join(dirPath, "README.md")

	// Output directory for generated HTML pages
	outputDir := "output"

	// Load templates
	templates := template.Must(template.ParseGlob("templates/*.html"))

	// Generate pages
	pages, err := generatePages(dirPath)
	if err != nil {
		log.Fatal(err)
	}

	// Create output directory if it doesn't exist
	if _, err := os.Stat(outputDir); os.IsNotExist(err) {
		os.Mkdir(outputDir, os.ModePerm)
	}

	// Generate static HTML pages
	for _, page := range pages {
		// Create the output file
		outputFile := filepath.Join(outputDir, page.Title+".html")
		file, err := os.Create(outputFile)
		if err != nil {
			log.Fatal(err)
		}

		// Combine the templates to generate the final HTML content
		err = templates.ExecuteTemplate(file, "page.html", page)
		if err != nil {
			log.Fatal(err)
		}

		file.Close()
	}

	// Generate the index page
	indexFile, err := os.Create(filepath.Join(outputDir, "index.html"))
	if err != nil {
		log.Fatal(err)
	}
	defer indexFile.Close()

	// Generate the buttons for each page
	var buttonHTML strings.Builder
	counter := 0
	for _, page := range pages {
		// Add a newline and padding if more than 3 buttons have been added
		if counter > 0 && counter%3 == 0 {
			buttonHTML.WriteString("<br>")
		}

		// Add the button
		buttonHTML.WriteString(fmt.Sprintf(`<a href="%s.html" class="button">%s</a>`, page.Title, page.Title))

		counter++
	}

	// Get the contents of the README.md file
	readmeContent, err := getReadmeContent(readmePath)
	if err != nil {
		log.Fatal(err)
	}

	// Combine the templates to generate the final HTML content for the index page
	data := struct {
		Buttons       string
		ReadmeContent string
	}{
		Buttons:       buttonHTML.String(),
		ReadmeContent: readmeContent,
	}

	err = templates.ExecuteTemplate(indexFile, "index.html", data)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Static pages generated successfully.")
}
