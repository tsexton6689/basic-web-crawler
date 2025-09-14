package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"golang.org/x/net/html"
)

func main() {
	// Flags ex go run . --page=https://golang.org
	webPage := flag.String("page", "", "The Webpage to crawl for links")

	flag.Parse()

	if *webPage == "" {
		log.Fatal("Error: Webpage flag must be set")
	}

	visited := make(map[string]bool)

	crawl(*webPage, visited)

	for key, _ := range visited {
		fmt.Println(key)
	}

	os.Exit(0)

}

func crawl(url string, visited map[string]bool) error {

	visited[url] = true

	links, err := fetch(url)

	if err != nil {
		return err
	}

	for _, link := range links {
		visited[link] = true
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
				if attr.Key == "href" && strings.Contains(attr.Val, "https://") {
					links = append(links, attr.Val)
				}
			}
		}
	}

	return links, nil
}
