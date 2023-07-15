const http = require("http");
const fs = require("fs");
const path = require("path");
const { spawnSync } = require("child_process");
const os = require("os");

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

  if (path.extname(filePath) === ".xml") {
    // Serve the RSS feed file directly
    fs.readFile(filePath, (err, content) => {
      if (err) {
        if (err.code === "ENOENT") {
          console.log(
            `HTTP ${new Date().toLocaleString()} ${
              req.socket.remoteAddress
            } GET ${req.url}`,
          );
          console.log(
            `HTTP ${new Date().toLocaleString()} ${
              req.socket.remoteAddress
            } Returned 404 in 1 ms`,
          );

          res.writeHead(404, { "Content-Type": "text/html" });
          res.end("<h1>404 Not Found</h1>");
        } else {
          console.log(
            `HTTP ${new Date().toLocaleString()} ${
              req.socket.remoteAddress
            } GET ${req.url}`,
          );
          console.log(
            `HTTP ${new Date().toLocaleString()} ${
              req.socket.remoteAddress
            } Returned 500 in 1 ms`,
          );

          res.writeHead(500, { "Content-Type": "text/html" });
          res.end("<h1>500 Server Error</h1>");
        }
      } else {
        res.writeHead(200, {
          "Content-Type": "application/rss+xml",
        });
        res.end(content);
      }
    });
  } else {
    // Serve other files (HTML, CSS, JS, etc.)
    fs.readFile(filePath, (err, content) => {
      if (err) {
        if (err.code === "ENOENT") {
          console.log(
            `HTTP ${new Date().toLocaleString()} ${
              req.socket.remoteAddress
            } GET ${req.url}`,
          );
          console.log(
            `HTTP ${new Date().toLocaleString()} ${
              req.socket.remoteAddress
            } Returned 404 in 1 ms`,
          );

          res.writeHead(404, { "Content-Type": "text/html" });
          res.end("<h1>404 Not Found</h1>");
        } else {
          console.log(
            `HTTP ${new Date().toLocaleString()} ${
              req.socket.remoteAddress
            } GET ${req.url}`,
          );
          console.log(
            `HTTP ${new Date().toLocaleString()} ${
              req.socket.remoteAddress
            } Returned 500 in 1 ms`,
          );

          res.writeHead(500, { "Content-Type": "text/html" });
          res.end("<h1>500 Server Error</h1>");
        }
      } else {
        const fileExtension = path.extname(filePath);
        const contentType = getContentType(fileExtension);

        console.log(
          `HTTP ${new Date().toLocaleString()} ${
            req.socket.remoteAddress
          } GET ${req.url}`,
        );
        console.log(
          `HTTP ${new Date().toLocaleString()} ${
            req.socket.remoteAddress
          } Returned 200 in 1 ms`,
        );

        res.writeHead(200, {
          "Content-Type": contentType,
        });
        res.end(content);
      }
    });
  }
});

server.listen(PORT, () => {
  console.log(`\n   ┌───────────────────────────────────────────┐`);
  console.log(`   │                                           │`);
  console.log(`   │   Serving!                                │`);
  console.log(`   │                                           │`);
  console.log(`   │   - Local:    http://localhost:${PORT}       │`);
  console.log(`   │   - Network:  http://192.168.1.130:${PORT}   │`);
  console.log(`   │                                           │`);
  console.log(`   └───────────────────────────────────────────┘\n`);

  const ipAddress = getIPAddress();
  copyToClipboard(ipAddress);
  console.log(`IP address copied to clipboard: ${ipAddress}`);
});

function getIPAddress() {
  const networkInterfaces = os.networkInterfaces();
  for (const interfaceName of Object.keys(networkInterfaces)) {
    const interfaces = networkInterfaces[interfaceName];
    for (const { family, address, internal } of interfaces) {
      if (family === "IPv4" && !internal) {
        return address;
      }
    }
  }
  return "Unknown";
}

function copyToClipboard(text) {
  if (process.platform === "linux") {
    const xclipProcess = spawnSync("xclip", ["-selection", "clipboard"], {
      input: text,
    });
    if (xclipProcess.error) {
      console.error("Failed to copy to clipboard:", xclipProcess.error.message);
    }
  } else if (process.platform === "darwin") {
    const pbcopyProcess = spawnSync("pbcopy", [], { input: text });
    if (pbcopyProcess.error) {
      console.error(
        "Failed to copy to clipboard:",
        pbcopyProcess.error.message,
      );
    }
  } else {
    console.warn("Clipboard access not supported on this platform.");
  }
}

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
    case ".xml":
      return "text/xml";
    default:
      return "application/octet-stream";
  }
}

server.on("request", (req, res) => {
  let filePath = path.join(STATIC_DIRECTORY, req.url);
  if (filePath === path.join(STATIC_DIRECTORY, "/")) {
    filePath = path.join(STATIC_DIRECTORY, "index.html");
  } else if (!path.extname(filePath)) {
    filePath += ".html";
  }

  fs.readFile(filePath, (err, content) => {
    if (err) {
      // Error handling code...
    } else {
      const fileExtension = path.extname(filePath);
      const contentType = getContentType(fileExtension);

      res.writeHead(200, {
        "Content-Type": contentType,
        "Content-Disposition":
          fileExtension === ".xml" ? "inline" : "attachment",
      });
      res.end(content);
    }
  });
});
