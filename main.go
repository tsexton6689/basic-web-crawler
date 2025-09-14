package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	"golang.org/x/net/html"
)

func main() {
	// Flags ex go run . --page=https://golang.org
	webPage := flag.String("page", "", "The Webpage to crawl for links")

	flag.Parse()

	if *webPage == "" {
		log.Fatal("Error: Webpage flag must be set")
	}

}

func crawl(url string, depth int, visited map[string]bool) error {
	if (depth == 3) {
		return nil
	}

	visited[url] = true

	links, err := fetch(url)

	if err != nil {
		return err
	}

	for _, link := range links {
		crawl(link, depth + 1, visited)
	}

	return nil
}

func fetch(url string) ([]string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	links := []string{}

	tokens := html.NewTokenizer(resp.Body)
	for {
		next := tokens.Next()
		if next == html.ErrorToken {
			break
		}

		token := tokens.Token()
		if token.Data == "a" {
			for _, attr := range token.Attr {
				if attr.Key == "href" {
					links = append(links, attr.Val)
				}
			}
		}
	}

	return links, nil
}