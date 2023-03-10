# scrapeThisFor

A command line tool for scraping URLs that contain a given keyword and downloading their contents.

## Roadmap:
- Nil, this is feature complete for what I needed.
---

## Requirements:

1. Clone this repo: `git clone https://github.com/alphastrata/scrapeThisFor.git`

2. Golang >= 1.18

3. A working internet connection (for `go mod` to work) and for the app to download content.

4. Some command line/terminal proficiency.

5. A browser if you want to read the documentation generated by godoc.

**Note**: You may need to restart your machine/re-login to your shell (depending on your OS).

## Building:

1. Run `go mod tidy` to download the required packages.
2. Run `go build main.go` to build the application.

## Example Usage:
* This command: 
`go run main.go https://huggingface.co/bigscience/bloom/tree/main model_000`
* will produce:
```bash
https://huggingface.co/bigscience/bloom/blob/main/model_00001-of-00072.safetensors
https://huggingface.co/bigscience/bloom/resolve/main/model_00001-of-00072.safetensors
... Snipped for brevity ...
https://huggingface.co/bigscience/bloom/blob/main/model_00045-of-00072.safetensors
https://huggingface.co/bigscience/bloom/resolve/main/model_00045-of-00072.safetensors
```
## Installation:
1. `go build -o scrapeThisFor main.go`
2. `sudo mv QuickScraper /usr/local/bin/` or it may be `/usr/bin` if that's your thing.
> Or:
1. `go get github.com/alphastrata/scrapeForThis`
