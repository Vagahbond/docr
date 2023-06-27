const http = require("http");
const fs = require("fs");
const path = require("path");

const PORT = 3000;

if (process.argv.length < 3) {
  console.error("Please provide the directory to be served.");
  process.exit(1);
}

const STATIC_DIRECTORY = process.argv[2];

const server = http.createServer((req, res) => {
  let filePath = path.join(STATIC_DIRECTORY, req.url);
  if (filePath === path.join(STATIC_DIRECTORY, "/")) {
    filePath = path.join(STATIC_DIRECTORY, "index.html");
  } else if (!path.extname(filePath)) {
    filePath += ".html";
  }

  fs.readFile(filePath, (err, content) => {
    if (err) {
      if (err.code === "ENOENT") {
        console.log(
          `HTTP ${new Date().toLocaleString()} ${
            req.connection.remoteAddress
          } GET ${req.url}`
        );
        console.log(
          `HTTP ${new Date().toLocaleString()} ${
            req.connection.remoteAddress
          } Returned 404 in 1 ms`
        );

        res.writeHead(404, { "Content-Type": "text/html" });
        res.end("<h1>404 Not Found</h1>");
      } else {
        console.log(
          `HTTP ${new Date().toLocaleString()} ${
            req.connection.remoteAddress
          } GET ${req.url}`
        );
        console.log(
          `HTTP ${new Date().toLocaleString()} ${
            req.connection.remoteAddress
          } Returned 500 in 1 ms`
        );

        res.writeHead(500, { "Content-Type": "text/html" });
        res.end("<h1>500 Server Error</h1>");
      }
    } else {
      const fileExtension = path.extname(filePath);
      const contentType = getContentType(fileExtension);

      console.log(
        `HTTP ${new Date().toLocaleString()} ${
          req.connection.remoteAddress
        } GET ${req.url}`
      );
      console.log(
        `HTTP ${new Date().toLocaleString()} ${
          req.connection.remoteAddress
        } Returned 200 in 1 ms`
      );

      res.writeHead(200, { "Content-Type": contentType });
      res.end(content);
    }
  });
});

server.listen(PORT, () => {
  console.log(`\n   ┌───────────────────────────────────────────┐`);
  console.log(`   │                                           │`);
  console.log(`   │   Serving!                                │`);
  console.log(`   │                                           │`);
  console.log(`   │   - Local:    http://localhost:${PORT}       │`);
  console.log(`   │   - Network:  http://192.168.1.130:${PORT}   │`);
  console.log(`   │                                           │`);
  console.log(`   │   Copied local address to clipboard!      │`);
  console.log(`   │                                           │`);
  console.log(`   └───────────────────────────────────────────┘\n`);
});

function getContentType(fileExtension) {
  switch (fileExtension) {
    case ".html":
      return "text/html";
    case ".css":
      return "text/css";
    case ".js":
      return "text/javascript";
    case ".png":
      return "image/png";
    case ".jpg":
      return "image/jpeg";
    default:
      return "application/octet-stream";
  }
}
