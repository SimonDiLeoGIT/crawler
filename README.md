# Concurrent Web Crawler in Go

A lightweight, concurrent web crawler implemented in Go that recursively discovers and crawls URLs starting from a seed webpage. It limits crawls to the same root domain to prevent runaway external indexing and utilizes worker goroutines to scale execution performance.

---

## 🚀 Key Features

*   **Concurrency Pool**: Employs a pool of configurable worker Goroutines executing crawling tasks in parallel.
*   **Domain Jailing**: Automatically filters URLs so that crawling is restricted strictly to the **same root domain** as the seed page (e.g., if the seed URL is `bbc.com`, it will only crawl subdomains and pages resolving to the `bbc.com` root registry domain).
*   **Depth Restriction**: Prevents infinite recursion by stopping crawls once they reach a predefined depth step limit from the seed URL.
*   **Deduplication & Visited Registry**: Keeps track of crawled directories to avoid revisiting the same page twice.
*   **URL Normalization**: Normalizes referenced URLs by stripping query parameters and hash fragments (e.g., `?ref=page` and `#section`), resolving relative paths matching the seed host.

---

## 📁 Project Architecture

The crawler codebase is organized as follows:

```text
├── cmd/
│   └── main.go       # Program entry-point, parses flags, triggers Crawl
├── internal/
│   ├── crawl.go      # Orchestrator coordinating worker loops, visited state, and final metrics
│   ├── fetcher.go    # HTTP client, HTML parser (goquery), same-domain filtering, and url normalization
│   └── worker.go     # Concurrency types (job, fetchResult) and worker goroutine execution loop
├── go.mod            # Dependency manager declaring Go 1.25.0
├── go.sum            # Checksums for external packages (goquery, golang.org/x/net)
└── README.md         # Document index (this file)
```

---

## 🛠️ Getting Started

### 📋 Prerequisites

*   Go installed (version **1.25.0** or above is recommended).

### ⚙️ Installation

First, pull down any required package dependencies locally (primarily `goquery` and `publicsuffix` libraries):

```bash
go mod download
```

---

## 🏃 How to Run the Project

You can run the web crawler directly using `go run` or compile it into a distribution binary.

### Option 1: Run standard execution directly
```bash
go run cmd/main.go [flags]
```

### Option 2: Build and run compilation executable
```bash
go build -o crawler cmd/main.go
./crawler [flags]
```

---

## ⚙️ Configuration Flags

Customize the crawler's behavior using command-line arguments:

| Flag | Type | Default | Description |
| :--- | :--- | :--- | :--- |
| `-url` | `string` | `"http://bbc.com/"` | Seed URL where the crawler begins crawling. |
| `-depth` | `int` | `3` | Maximum recursion limit/hops from the initial seed page. |
| `-workers` | `int` | `8` | Number of concurrent worker goroutines fetching pages in parallel. |

### Configuration Examples

1.  **Crawl with custom depth and worker counts:**
    ```bash
    go run cmd/main.go -url "https://crawler-test.com" -depth 2 -workers 12
    ```

2.  **Run with default configuration:**
    ```bash
    go run cmd/main.go
    ```

---

## 📊 Sample Output

Upon finishing the execution, the crawler prints out URLs discovered in real-time and logs final metrics on completion:

```text
https://crawler-test.com/
https://subdomain.crawler-test.com
http://crawler-test.com/

Crawl completed

URLs discovered 415
URLs crawled 5
URLs skiped 410

Depth: 1
Workers: 2
Duration: 630.048ms
```
