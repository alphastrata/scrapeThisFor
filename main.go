package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"golang.org/x/net/html"
)

func getResponse(url string) (*http.Response, error) {
	resp, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("receiving status of %d for url: %s", resp.StatusCode, url)
	}
	return resp, nil
}

func DownloadFile(url, filepath *string) error {
	out, err := os.Create(*filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	resp, err := http.Get(*url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return err
	}

	return nil
}

func worker(url, filename *string) {
	fmt.Printf("Worker %s starting...\n", *filename)
	e := DownloadFile(url, filename)
	if e != nil {
		fmt.Println(e)
	}
	fmt.Printf("worker %s finishing!\n", *filename)

}

var wg sync.WaitGroup

func main() {
	if len(os.Args) < 3 {
		fmt.Println("Usage: myScrape url keyword [download]\nExample: go run . https://huggingface.co/bigscience/bloom/tree/main model_000 download")
		return
	}
	rawURL := os.Args[1]
	keyword := os.Args[2]
	downloadAll := false
	if len(os.Args) > 3 && os.Args[3] == "download" {
		downloadAll = true
	}

	u, err := url.Parse(rawURL)
	if err != nil {
		fmt.Println("Error: invalid URL")
		return
	}

	resp, err := http.Get(u.String())
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	doc, err := html.Parse(strings.NewReader(string(body)))
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	hrefs := extractHrefs(doc)

	var filteredHrefs []string
	for _, href := range hrefs {
		absURL, err := u.Parse(href)
		if err != nil {
			continue // Skip invalid URLs.
		}
		if strings.Contains(absURL.String(), keyword) {
			filteredHrefs = append(filteredHrefs, absURL.String())
		}
	}
	fmt.Printf("Number of hrefs looked at: %d, Number of hits: %d\n", len(hrefs), len(filteredHrefs))

	for _, url := range filteredHrefs {
		if downloadAll {

			resp, err := getResponse(url)
			if err != nil {
				fmt.Println("could not get response", err)
				os.Exit(1)
			}
			defer resp.Body.Close()

			if resp.ContentLength <= 0 {
				fmt.Printf("can't parse content length for %s, aborting...\n", url)
				continue
			}

			filename := filepath.Base(url)
			if _, err := os.Stat(filename); !os.IsNotExist(err) {
				fmt.Printf("%s, already exists on disk", filename)
				continue
			}
			file, err := os.Create(filename)
			if err != nil {
				fmt.Println("could not create file:", err)
				os.Exit(1)
			}
			defer file.Close()

			// Start the download
			wg.Add(1)
			go func(url *string, filename *string) {
				worker(url, filename)
				wg.Done()
			}(&url, &filename)

		} else {
			fmt.Println(url)
		}

	}
	wg.Wait()

}

func extractHrefs(n *html.Node) []string {
	if n == nil {
		return nil
	}

	var hrefs []string

	if n.Type == html.ElementNode && n.Data == "a" {
		for _, attr := range n.Attr {
			if attr.Key == "href" {
				hrefs = append(hrefs, attr.Val)
			}
		}
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		hrefs = append(hrefs, extractHrefs(c)...)
	}

	return hrefs
}
