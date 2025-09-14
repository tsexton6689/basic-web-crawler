package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"strings"
	"sync"

	"golang.org/x/net/html"
)

func main() {
	// Flags ex go run . --page=https://golang.org
	webPage := flag.String("page", "", "The Webpage to crawl for links")

	flag.Parse()

	if *webPage == "" {
		log.Fatal("Error: Webpage flag must be set")
	}

	maxDepth := 1
	visited := make(map[string]bool)
	var mu sync.Mutex
	var wg sync.WaitGroup

	var crawl func(string, int)

	crawl = func(url string, depth int) {
		defer wg.Done()

		if depth > maxDepth {
			return
		}

		mu.Lock()
		if visited[url] {
			mu.Unlock()
			return
		}

		visited[url] = true
		mu.Unlock()

		links, err := fetch(url)
		if err != nil {
			fmt.Println("error:", err)
			return
		}

		fmt.Println("Found:", url)
		for _, link := range links {
			wg.Add(1)
			go crawl(link, depth+1)
		}
	}

	wg.Add(1)
	go crawl(*webPage, 0)
	wg.Wait()
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
