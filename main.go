package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strings"

	"golang.org/x/net/html"
)

func main() {
	if len(os.Args) != 3 {
		fmt.Println("Usage: myScrape url keyword")
		return
	}
	rawURL := os.Args[1]
	keyword := os.Args[2]

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

	body, err := ioutil.ReadAll(resp.Body)
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
	for _, href := range filteredHrefs {
		fmt.Println(href)
	}

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
