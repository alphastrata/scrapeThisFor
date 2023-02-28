package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	"golang.org/x/net/html"
)

func main() {
	if len(os.Args) != 3 {
		fmt.Println("Usage: myScrape url keyword")
		return
	}

	url := os.Args[1]
	keyword := os.Args[2]

	resp, err := http.Get(url)
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
	filteredHrefs := filterHrefs(hrefs, keyword)

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

func filterHrefs(hrefs []string, keyword string) []string {
	var filteredHrefs []string

	for _, href := range hrefs {
		if strings.Contains(href, keyword) {
			filteredHrefs = append(filteredHrefs, href)
		}
	}

	return filteredHrefs
}
