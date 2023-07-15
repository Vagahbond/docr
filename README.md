# 📚 Docr

> Docr is a barebones and dead simple static site generator in Go. It converts Markdown content to HTML using the goldmark library, applies a template, and generates static HTML files according to your configuration and templates.

## Features

- Converts Markdown files to HTML pages, with Github Flavored Markdown support
- Generates an RSS feed for the website
- Supports customization through configuration file or environment variables
- Provides a simple template system for consistent page layout (supports CSS, JS, HTML and XLS templates)

## Usage

### Installing Docr

Built binaries are automatically published to Github Releases, so you do not need to install the Go toolkit and build Docr yourself. Simply obtain
the latest release from the [releases tab]() and run the binary.

You will also need the contents of the [templates/] directory, which you can obtain from the same source as they are uploaded there.
