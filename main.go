package main

import (
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"os"
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
	tmpl, err := template.New("index").Parse(`
		<!DOCTYPE html>
		<html>
		<head>
			<meta charset="UTF-8">
			<title>My Website</title>
			<style>
				html, body {
					height: 100%;
					margin: 0;
					padding: 0;
				}
				body {
					display: flex;
					flex-direction: column;
				}
				.container {
					flex: 1;
				}
				.navbar {
					background-color: lightgray;
					padding: 10px;
				}
				.navbar a {
					margin-right: 10px;
					text-decoration: none;
				}
				.main {
					padding: 20px;
				}
				.page-list {
					list-style-type: none;
				}
				.page-list li {
					margin-bottom: 10px;
				}
				.footer {
					padding: 20px;
					text-align: center;
					background-color: lightgray;
				}
				.dark-mode-button {
					position: absolute;
					top: 10px;
					right: 10px;
				}
				.dark-mode {
					background-color: #000;
					color: #fff;
				}
			</style>
			<script>
				function toggleDarkMode() {
					var body = document.querySelector('body');
					body.classList.toggle('dark-mode');
				}
			</script>
		</head>
		<body>
			<div class="navbar">
				<a href="{{.Navbar.Home}}">Home</a>
				<a href="{{.Navbar.GitHub}}">GitHub</a>
				<a href="{{.Navbar.Email}}">Email</a>
				<button class="dark-mode-button" onclick="toggleDarkMode()">Dark Mode</button>
			</div>
			<div class="container">
				<div class="main">
					<h1>Welcome to My Website</h1>
					<ul class="page-list">
						{{range .Pages}}
							<li><a href="{{.}}">{{.}}</a></li>
						{{end}}
					</ul>
				</div>
			</div>
			<div class="footer">
				<p>{{.Footer.Text}}</p>
			</div>
		</body>
		</html>
	`)

	if err != nil {
		return "", err
	}

	builder := strings.Builder{}
	err = tmpl.Execute(&builder, struct {
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

		markdownPath := fmt.Sprintf("%s/%s", markdownDir, file.Name())
		pagePath := fmt.Sprintf("%s/%s.html", outputDir, strings.TrimSuffix(file.Name(), ".md"))

		html, err := GeneratePageHTML(markdownPath)
		if err != nil {
			log.Fatal(err)
		}

		tmpl, err := template.New("page").Parse(`
			<!DOCTYPE html>
			<html>
			<head>
				<meta charset="UTF-8">
				<title>{{.Title}}</title>
					<style>
					html, body {
						height: 100%;
						margin: 0;
						padding: 0;
					}
					body {
						display: flex;
						flex-direction: column;
					}
					.container {
						flex: 1;
					}
					.navbar {
						background-color: lightgray;
						padding: 10px;
					}
					.navbar a {
						margin-right: 10px;
						text-decoration: none;
					}
					.main {
						padding: 20px;
					}
					.page-content {
						margin-top: 20px;
					}
					.footer {
						padding: 20px;
						text-align: center;
						background-color: lightgray;
					}
					.dark-mode-button {
						position: absolute;
						top: 10px;
						right: 10px;
					}
					.dark-mode {
						background-color: #000;
						color: #fff;
					}
				</style>
				<script>
					function toggleDarkMode() {
						var body = document.querySelector('body');
						body.classList.toggle('dark-mode');
					}
				</script>
			</head>
			<body>
				<div class="navbar">
					<a href="../index.html">Home</a>
					<a href="#{{.Title}}">{{.Title}}</a>
					<a href="../index.html#footer">Footer</a>
					<button class="dark-mode-button" onclick="toggleDarkMode()">Dark Mode</button>
				</div>
				<div class="container">
					<div class="main">
						<div class="page-content">
							{{.Body}}
						</div>
					</div>
				</div>
				<div class="footer" id="footer">
					<p>© 2023 My Website. All rights reserved.</p>
				</div>
			</body>
			</html>
		`)

		if err != nil {
			log.Fatal(err)
		}

		fileContent := strings.Builder{}
		err = tmpl.Execute(&fileContent, Page{
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

	indexPath := fmt.Sprintf("%s/index.html", outputDir)
	err = ioutil.WriteFile(indexPath, []byte(indexHTML), 0644)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Static webpage built successfully in the '%s' directory.\n", outputDir)
}
