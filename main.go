package main

import (
	"encoding/json"
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

// Footer represents the footer section in the template.
type Footer struct {
	Year string
}

// Navbar represents the navigation bar section in the template.
type Navbar struct {
	Pages []Page
}

// Settings represents the configuration settings.
type Settings struct {
	GithubUsername string `json:"githubUsername"`
	WebsiteName    string `json:"websiteName"`
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
		goldmark.WithExtensions(extension.GFM, extension.Table, extension.Footnote, extension.Linkify),
		goldmark.WithParserOptions(parser.WithAutoHeadingID()),
		goldmark.WithRendererOptions(html.WithHardWraps()),
	)

	var buf strings.Builder
	if err := md.Convert(content, &buf); err != nil {
		log.Fatal(err)
	}

	return buf.String()
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

func main() {
	// Directory containing the markdown files
	dirPath := "markdown"

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

	// Load settings from settings.json file
	settingsFile, err := ioutil.ReadFile("settings.json")
	if err != nil {
		log.Fatal(err)
	}

	var settings Settings
	err = json.Unmarshal(settingsFile, &settings)
	if err != nil {
		log.Fatal(err)
	}

	// Generate individual pages
	for _, page := range pages {
		// Create the output file
		pageFile, err := os.Create(filepath.Join(outputDir, fmt.Sprintf("%s.html", page.Title)))
		if err != nil {
			log.Fatal(err)
		}
		defer pageFile.Close()

		// Combine the templates to generate the final HTML content for the page
		data := struct {
			Title          string
			Content        string
			GithubUsername string
			WebsiteName    string
			Navbar         Navbar
			Footer         Footer
		}{
			Title:          page.Title,
			Content:        page.Content,
			GithubUsername: settings.GithubUsername,
			WebsiteName:    settings.WebsiteName,
			Navbar:         Navbar{Pages: pages},
			Footer:         Footer{Year: "2023"}, // Update with the appropriate year
		}

		err = templates.ExecuteTemplate(pageFile, "page.html", data)
		if err != nil {
			log.Fatal(err)
		}

		log.Printf("Generated page: %s.html\n", page.Title)
	}

	// Read the README.md file
	readmeContent, err := ioutil.ReadFile(filepath.Join(dirPath, "README.md"))
	if err != nil {
		log.Fatal(err)
	}

	// Convert README.md content to HTML
	readmeHTML := renderMarkdown(readmeContent)

	// Combine the templates to generate the final HTML content for the index page
	indexData := struct {
		WebsiteName    string
		GithubUsername string
		ReadmeContent  string
		Buttons        string
		Navbar         Navbar
		Footer         Footer
	}{
		WebsiteName:    settings.WebsiteName,
		GithubUsername: settings.GithubUsername,
		ReadmeContent:  readmeHTML,
		Buttons:        generateButtons(pages),
		Navbar:         Navbar{Pages: pages},
		Footer:         Footer{Year: "2023"}, // Update with the appropriate year
	}

	// Create the index.html file
	indexFile, err := os.Create(filepath.Join(outputDir, "index.html"))
	if err != nil {
		log.Fatal(err)
	}
	defer indexFile.Close()

	err = templates.ExecuteTemplate(indexFile, "index.html", indexData)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Static pages generated successfully.")
}

// generateButtons generates the HTML buttons for each page (excluding index.html).
func generateButtons(pages []Page) string {
	var buttons strings.Builder
	for _, page := range pages {
		if page.Title != "index" {
			button := fmt.Sprintf("<a href=\"%s.html\" class=\"button\">%s</a>", page.Title, page.Title)
			buttons.WriteString(button)
		}
	}

	return buttons.String()
}
