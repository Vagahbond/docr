package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer/html"
)

// Page represents a single page with its title and HTML content.
type Page struct {
	Title   string
	Content string
}

// FooterData represents the data for the footer.
type FooterData struct {
	Year int
}

// NavbarData represents the data for the navbar.
type NavbarData struct {
	Pages []Page
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

			htmlContent := renderMarkdown(content)

			var page Page
			page.Title = strings.TrimSuffix(info.Name(), ".md")
			page.Content = htmlContent

			pages = append(pages, page)
		}

		return nil
	})

	return pages, err
}

// renderMarkdown converts the given Markdown content to HTML using goldmark.
func renderMarkdown(content []byte) string {
	md := goldmark.New(
		goldmark.WithExtensions(extension.GFM),
		goldmark.WithParserOptions(parser.WithAutoHeadingID()),
		goldmark.WithRendererOptions(html.WithHardWraps()),
	)

	var buf strings.Builder
	if err := md.Convert(content, &buf); err != nil {
		log.Fatal(err)
	}

	return buf.String()
}

// getReadmeContent reads the contents of the README.md file, converts it to HTML,
// and returns it as a string.
func getReadmeContent(readmePath string) (string, error) {
	content, err := ioutil.ReadFile(readmePath)
	if err != nil {
		return "", err
	}

	htmlContent := renderMarkdown(content)
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

	// Copy static files to output directory
	err = copyStaticFiles(outputDir)
	if err != nil {
		log.Fatal(err)
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
		err = templates.ExecuteTemplate(file, "page.html", struct {
			Page
			Footer FooterData
			Navbar NavbarData
		}{
			Page:   page,
			Footer: FooterData{Year: 2023}, // Replace with the actual year if needed
			Navbar: NavbarData{Pages: pages},
		})
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
		Footer        FooterData
		Navbar        NavbarData
	}{
		Buttons:       buttonHTML.String(),
		ReadmeContent: readmeContent,
		Footer:        FooterData{Year: 2023}, // Replace with the actual year if needed
		Navbar:        NavbarData{Pages: pages},
	}

	err = templates.ExecuteTemplate(indexFile, "index.html", data)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Static pages generated successfully.")
}

// copyStaticFiles copies the static files (CSS, JS, etc.) to the output directory.
func copyStaticFiles(outputDir string) error {
	// Source directories containing static files
	cssDir := "templates/css"
	jsDir := "templates/js"

	// Create the CSS directory in the output directory
	err := os.MkdirAll(filepath.Join(outputDir, "css"), os.ModePerm)
	if err != nil {
		return err
	}

	// Copy CSS files
	err = filepath.Walk(cssDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() {
			// Read the CSS file
			content, err := ioutil.ReadFile(path)
			if err != nil {
				return err
			}

			// Create the corresponding file in the output directory
			outputPath := filepath.Join(outputDir, "css", filepath.Base(path))
			err = ioutil.WriteFile(outputPath, content, os.ModePerm)
			if err != nil {
				return err
			}
		}

		return nil
	})
	if err != nil {
		return err
	}

	// Create the JS directory in the output directory
	err = os.MkdirAll(filepath.Join(outputDir, "js"), os.ModePerm)
	if err != nil {
		return err
	}

	// Copy JS files
	err = filepath.Walk(jsDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() {
			// Read the JS file
			content, err := ioutil.ReadFile(path)
			if err != nil {
				return err
			}

			// Create the corresponding file in the output directory
			outputPath := filepath.Join(outputDir, "js", filepath.Base(path))
			err = ioutil.WriteFile(outputPath, content, os.ModePerm)
			if err != nil {
				return err
			}
		}

		return nil
	})
	if err != nil {
		return err
	}

	return nil
}
